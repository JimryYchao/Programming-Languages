## CSharp 并行编程

---
### 1. Parallel

任务并行库 TPL 基于 `Parallel` 支持数据并行。
- `Parallel.For` 并行执行 *For* 循环。
- `Parallel.ForEach` 并行遍历集合。
- `Parallel.Invoke` 并行执行多个任务。

```csharp
// 累加 1 ~ 1_0000_0000
// For
{
    long total  = 0;
    var nums = Enumerable.Range(0, 1_0000_0000).ToArray();
    Parallel.For<long>(0, nums.Length, () => 0, (i, loop, subtotal) => {
        subtotal += nums[i];
        return subtotal;
    }, (final) => Interlocked.Add(ref total, final));
    Console.WriteLine(total);
}
// Foreach
{
    long total = 0;
    var nums = Enumerable.Range(0, 1_0000_0000).ToArray();
    Parallel.ForEach<int, long>(nums, () => 0, (e, loop, subtotal) => {
        subtotal += e;
        return subtotal;
    }, (final) => Interlocked.Add(ref total, final));
    Console.WriteLine(total);
}
```

`Partitioner` 为 `Parallel.For` 创建数据分区，以便每个分区仅调用一次 *Loop* 委托，而不是每个迭代调用一次委托。

```csharp
{ 
    long total = 0;
    var nums = Enumerable.Range(0, 1_0000_0000).ToArray();
    var rangePartitioner = Partitioner.Create(0, nums.Length, 100000);
    Parallel.ForEach(rangePartitioner, () => 0, (range, loopState, subtotal) => {
        for (int i = range.Item1; i < range.Item2; i++)
            subtotal += nums[i];
        return subtotal;
    }, (final) => Interlocked.Add(ref total, final));
}
```

---
### 2. Task

`Task` 代表异步任务，`Task.ContinueWith` 创建任务延续。对于 `Task<TResult>` 任务延续，`Unwrap` 对延续进行解包。`Parallel.Invoke` 在后台隐式创建 `Task` 执行并发任务。`TaskFactory.StartNew` 创建并启用一个任务，`Task` 包含一个默认的静态实例 `Factory`。`Task.Factory` 可以创建附加（AttachedToParent）子任务。

```csharp
class Sample {
    class CustomData {
        public long CreationTime;
        public int Name;
        public int ThreadId;
    }
    static void Main(string[] args) {
        Task.Factory.StartNew(() => {
            for (int i = 0; i < 10; i++) {
                Task.Factory.StartNew((obj) => {
                    if (obj is CustomData data) 
                        data.ThreadId = Thread.CurrentThread.ManagedThreadId;
                }, new CustomData() { Name = i, CreationTime = DateTime.Now.Ticks },
                 TaskCreationOptions.AttachedToParent)
                .ContinueWith(t => {
                    if (t.AsyncState is CustomData data)
                        Console.WriteLine($"Task #{data.Name} created at {data.CreationTime} on thread #{data.ThreadId}.");
                }, TaskContinuationOptions.OnlyOnRanToCompletion | TaskContinuationOptions.AttachedToParent);
            }
        }).Wait();
    }
}
```

>---
#### 2.1. TaskCompletionSource

`TaskCompletionSource` 直接用于手动创建和控制 `Task` 状态，未绑定到代理。

> 包装 EAP, APM 为 TAP

```csharp
interface APM<T> {
    IAsyncResult BeginInvoke(AsyncCallback? callback, object? state, CancellationToken token);
    T EndInvoke(IAsyncResult asyncResult);
}
static class APMTask {
    public static async Task<T> InvokeAsync<T>(this APM<T> apm, dynamic[]? args, AsyncCallback? callback, object? state, CancellationToken token) {
        var tcs = new TaskCompletionSource<T>();
        IAsyncResult result = apm.BeginInvoke(callback, state, token);
        try {
            var rt = apm.EndInvoke(result);  // block
            tcs.SetResult(rt);
        } catch (OperationCanceledException) {
            tcs.SetCanceled(token);
        } catch (Exception ex) {
            tcs.SetException(ex);
        }
        return await tcs.Task;
    }
}
```

---
### 3. 数据流

TPL 包含三种数据流块：`ISourceBlock<T>` 源块，`ITargetBlock<T>` 目标块，`IPropagatorBlock<T>` 传播器块。预定义数据流块类型包含缓冲块、执行块、分组块。

>---
#### 3.1. 缓冲块

`BufferBlock` 先进先出异步消息传送结构，可多源写入，多目标读取，目标接收消息后从缓冲区移除。

```csharp
BufferBlock<int> bufferBlock = new BufferBlock<int>();
for (int i = 0; i < 10; i++) 
    bufferBlock.Post(i);
bufferBlock.Complete();
await foreach (var item in bufferBlock.ReceiveAllAsync()) 
    Console.WriteLine(item);  // 0,1,2,3,4,5,6,7,8,9
```

`BroadcastBlock<T>` 广播最新消息到所有目标，消息不会移除。

```csharp
var broadcastBlock = new BroadcastBlock<double>(null);
broadcastBlock.Post(Math.PI);
for (int i = 0; i < 3; i++)
   Console.WriteLine(broadcastBlock.Receive());
/* Output:
   3.14159265358979
   3.14159265358979
   3.14159265358979
 */
```

`WaitOnceBlock<T>` 仅写入一次，该消息不会删除，且后续消息会被忽略。

```csharp
var writeOnceBlock = new WriteOnceBlock<string>(null);
Parallel.Invoke(
   () => writeOnceBlock.Post("Message 1"),
   () => writeOnceBlock.Post("Message 2"),
   () => writeOnceBlock.Post("Message 3"));
Console.WriteLine(writeOnceBlock.Receive());  // Message 1
```

>---
#### 3.2. 执行块

`ActionBlock<TInput>` 在接收数据时调用委托。

```csharp
ActionBlock<int> actionBlock = new ActionBlock<int>(n => Console.WriteLine(n));
for (int i = 0; i < 10; i++) 
    actionBlock.Post(i);
actionBlock.Complete();
actionBlock.Completion.Wait(); // 0,1,2,...
```

`TransformBlock<TInput, TOutput>` 类似于 `ActionBlock<TInput>`，既是源块也是目标块。

```csharp
TransformBlock<int, int> transformBlock = new TransformBlock<int, int>(n => n * 2);
for (int i = 1; i < 10; i++)
    transformBlock.Post(i);
transformBlock.Complete();
await foreach (var i in transformBlock.ReceiveAllAsync())
    Console.Write($"{i},");  // 2,4,6,8,10,12,14,16,18,
```

`TransformManyBlock<TInput, TOutput>` 类似于 `TransformBlock<TInput, TOutput>`，为每一个输入值生成零到多个输出值。

```csharp
var transformManyBlock = new TransformManyBlock<string, char>(s => s.ToCharArray());  // Func<TInput, IEnumerable<TOutput>>
transformManyBlock.Post("Hello");
transformManyBlock.Post("World");
transformManyBlock.Complete();
await foreach (var c in transformManyBlock.ReceiveAllAsync()) 
    Console.Write($"{c},");  // H,e,l,l,o,W,o,r,l,d,
```

>---
#### 3.3. 分组块

`BatchBlock<T>` 批处理输入的数据，合并到输出中。贪婪模式下，接受每一个输出，输入计数达到阈值时传播数组；非贪婪模式下，推迟所有输入，知道足够的源数据后形成批处理。

```csharp
var batchBlock = new BatchBlock<int>(5);
for (int i = 0; i < 13; i++) 
    batchBlock.Post(i);
batchBlock.Complete();
await foreach (var batch in batchBlock.ReceiveAllAsync())
    Console.WriteLine(batch.Sum());   // 10，35，33
```

`JoinBlock<T1, T2>` 和 `JoinBlock<T1, T2, T3>` 收集输入合并为 `Tuple<T1, T2>` 或 `Tuple<T1, T2, T3>`。

```csharp
JoinBlock<int, int> joinBlock = new JoinBlock<int, int>();
for (int i = 0; i < 5; i++) {
    joinBlock.Target1.Post(i);
    joinBlock.Target2.Post(i + 5);
}
joinBlock.Complete();
await foreach(var tuple in joinBlock.ReceiveAllAsync())
    Console.WriteLine(tuple);  // (0,5), (1,6), (2,7), (3,8), (4,9)
```

`BatchedJoinBlock<T1, T2>` 类似于 `JoinBlock<T1, T2>`，输入批处理，达到计数时异步传播 `Tuple<IList<T1>, IList<T2>>`。

```csharp
Func<int, int> DoWork = n => n < 0 ? throw new ArgumentOutOfRangeException() : n;
var batchedJoinBlock = new BatchedJoinBlock<int, Exception>(4);
foreach (int i in new int[] { 5, 6, -7, -22, 13, 55, 0 }) {
    try {
        batchedJoinBlock.Target1.Post(DoWork(i));
    } catch (ArgumentOutOfRangeException e) {
        batchedJoinBlock.Target2.Post(e);
    }
}
batchedJoinBlock.Complete();
await foreach (var results in batchedJoinBlock.ReceiveAllAsync()) {
    Console.WriteLine("Receice : ");
    foreach (int n in results.Item1)
        Console.WriteLine(n);
    foreach (Exception e in results.Item2)
        Console.WriteLine(e.Message);
}
/* Output:
Receice :
    5
    6
    Specified argument was out of the range of valid values.
    Specified argument was out of the range of valid values.
Receice :
    13
    55
    0
 */
```

>---
#### 3.4. 数据流管道

`block.LinkTo(targetBlock, linkOpts)`  连接源块到目标块。`linkOpts.MaxMessages` 为 1 表示目标从源接收一次消息后断开链接，默认为 -1。

```csharp
var downloadStringBlock = new TransformBlock<string, string>(async uri => {
    Console.WriteLine($"Downloading '{uri}'...");
    return await new HttpClient(new HttpClientHandler { AutomaticDecompression = System.Net.DecompressionMethods.GZip }).GetStringAsync(uri);
});
var createWordListBlock = new TransformBlock<string, string[]>(text => {
    Console.WriteLine("Creating word list...");
    char[] tokens = text.Select(c => char.IsLetter(c) ? c : ' ').ToArray();
    text = new string(tokens);
    return text.Split(new[] { ' ' }, StringSplitOptions.RemoveEmptyEntries);
});
var filterWordListBlock = new TransformBlock<string[], string[]>(words => {
    Console.WriteLine("Filtering word list...");
    return words.Where(word => word.Length > 3).Distinct().ToArray();
});
var findReversedWordsBlock = new TransformManyBlock<string[], string>(words => {
    Console.WriteLine("Finding reversed words...");
    var wordsSet = new HashSet<string>(words);
    return from word in words.AsParallel()
           let reverse = new string(word.Reverse().ToArray())
           where word != reverse && wordsSet.Contains(reverse)
           select word;
});
var printReversedWordsBlock = new ActionBlock<string>(reversedWord => {
    Console.WriteLine($"Found reversed words {reversedWord}/{new string(reversedWord.Reverse().ToArray())}");
});

// linkTo : download > createWord > filterWord > findReversedWord > printReversedWord
var linkOpts = new DataflowLinkOptions { PropagateCompletion = true };
downloadStringBlock.LinkTo(createWordListBlock, linkOpts);
createWordListBlock.LinkTo(filterWordListBlock, linkOpts);
filterWordListBlock.LinkTo(findReversedWordsBlock, linkOpts);
findReversedWordsBlock.LinkTo(printReversedWordsBlock, linkOpts);

// Post
downloadStringBlock.Post("http://www.gutenberg.org/cache/epub/16452/pg16452.txt");
downloadStringBlock.Complete();
printReversedWordsBlock.Completion.Wait();
```

---