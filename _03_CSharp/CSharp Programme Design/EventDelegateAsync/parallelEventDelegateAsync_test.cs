class ParallelEventDelegateAsyncTest {
    public static void Test() {
        // 测试并发操作
        Console.WriteLine("Testing concurrent operations...");

        var asyncOp = EventDelegateAsync.Create(() => {
            Console.WriteLine($"Operation started on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(1000);
            Console.WriteLine($"Operation completed on thread {Thread.CurrentThread.ManagedThreadId}");
        });

        asyncOp.InvokeCompleted += (sender, e) => {
            Console.WriteLine($"Operation completed. User state: {e.UserState}");
        };

        // 启动多个并发操作
        for (int i = 0; i < 5; i++) {
            object userState = $"Operation{i}";
            asyncOp.InvokeAsync(null, userState);
            Console.WriteLine($"Started operation {i}");
        }

        // 等待所有操作完成
        Console.WriteLine("Waiting for all operations to complete...");
        Thread.Sleep(2000);

        // 测试带返回值的并发操作
        Console.WriteLine("\nTesting concurrent functions...");

        var asyncFunc = EventDelegateAsync.Create<Func<int>, int>(() => {
            Console.WriteLine($"Function started on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(500);
            return new Random().Next(1, 100);
        });

        asyncFunc.InvokeCompleted += (sender, e) => {
            Console.WriteLine($"Function completed. Result: {e.Result}, User state: {e.UserState}");
        };

        // 启动多个并发函数
        for (int i = 0; i < 3; i++) {
            object userState = $"Function{i}";
            asyncFunc.InvokeAsync(null, userState);
            Console.WriteLine($"Started function {i}");
        }
        Thread.Sleep(2000);
    }
}