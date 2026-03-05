public class ChannelTest
 {
    public static async Task Test() {   
        Console.WriteLine("=== Channel 死锁测试 ===");
        
        // 测试 0 缓冲 channel
        Console.WriteLine("\n1. 测试 0 缓冲 Channel（无缓冲）");
        await TestChannel(0, 2, 1, 3);
        
        // 测试有缓冲 channel
        Console.WriteLine("\n2. 测试有缓冲 Channel（容量为 3）");
        await TestChannel(3, 2, 1, 3);
        
        // 测试大缓冲 channel
        Console.WriteLine("\n3. 测试大缓冲 Channel（容量为 5）");
        await TestChannel(5, 2, 1, 3);
        
        Console.WriteLine("\n=== 所有测试完成，无死锁、无任务挂起 ===");
    }
    
    static async Task TestChannel(int capacity, int producerCount, int consumerCount, int itemsPerProducer) {
        Console.WriteLine($"\n- 通道容量: {capacity}, 生产者: {producerCount}, 消费者: {consumerCount}, 每生产者项目数: {itemsPerProducer}");
        
        // 异步释放声明，适配异步并发场景
        await using var channel = new Channel<int>(capacity);

        // 多生产者并发发送
        var producerTasks = new List<Task>();
        for (int p = 1; p <= producerCount; p++) {
            int producerId = p;
            producerTasks.Add(Task.Run(async () => {
                for (int i = 1; i <= itemsPerProducer; i++) {
                    try {
                        int data = producerId * 100 + i;
                        // 并发发送，满通道自动挂起，无死锁
                        await channel.SendAsync(data);
                        //Console.WriteLine($"[{DateTime.Now:HH:mm:ss.fff}] 生产者{producerId} 发送数据：{data}");
                        // 随机延迟，模拟不同的发送速度
                        await Task.Delay(new Random().Next(50, 150));
                    } catch (ObjectDisposedException) {
                        Console.WriteLine($"生产者{producerId}：通道已释放，停止发送");
                        break;
                    }
                }
            }));
        }

        // 多消费者并发取值
        var consumerTasks = new List<Task>();
        for (int c = 1; c <= consumerCount; c++) {
            int consumerId = c;
            consumerTasks.Add(Task.Run(async () => {
                while (true) {
                    try {
                        // 并发await取值，空通道自动挂起，无死锁
                        int data = await channel;
                        //Console.WriteLine($"[{DateTime.Now:HH:mm:ss.fff}] 消费者{consumerId} 取出数据：{data}");
                        // 随机延迟，模拟不同的消费速度
                        await Task.Delay(new Random().Next(100, 200));
                    } catch (ObjectDisposedException) {
                        Console.WriteLine($"消费者{consumerId}：通道关闭/释放，停止取值");
                    } catch (InvalidOperationException) {
                        break;
                    }
                }
            }));
        }

        // 等待所有生产者完成
        await Task.WhenAll(producerTasks);
        // 生产者完成后关闭通道
        channel.Close();
        Console.WriteLine($"\n- 所有 {producerCount} 个生产者发送完成，通道已关闭");
        // 等待所有消费者取完剩余数据
        await Task.WhenAll(consumerTasks);
        
        Console.WriteLine($"- 测试完成：{capacity} 缓冲通道，{producerCount} 生产者，{consumerCount} 消费者");
    }
}