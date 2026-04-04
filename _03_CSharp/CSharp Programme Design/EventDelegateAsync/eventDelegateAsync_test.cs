
public class EventDelegateAsyncTest {
    public static void Test()
    {
        // 测试无返回值
        Console.WriteLine("Testing Action async...");
        var actionAsync = EventDelegateAsync.Create(() => {
            Console.WriteLine($"Action executing on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(2000); // 模拟耗时操作
            Console.WriteLine("Action completed");
        });
        actionAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Action error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Action cancelled");
            } else {
                Console.WriteLine("Action completed successfully");
            }
        };
        actionAsync.InvokeAsync(null);
        while (actionAsync.IsBusy)
            continue;   // 阻塞

        // 测试带返回值
        Console.WriteLine("\nTesting Func async...");
        var funcAsync = EventDelegateAsync.Create<Func<int>, int>(() => {
            Console.WriteLine($"Func executing on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(1500); // 模拟耗时操作
            return 42;
        });
        funcAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Func error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Func cancelled");
            } else {
                Console.WriteLine($"Func completed successfully. Result: {e.Result}");
            }
        };
        funcAsync.InvokeAsync(null);
        while (funcAsync.IsBusy)
            continue;   // 阻塞

        // 测试带用户状态
        Console.WriteLine("\nTesting with user state...");
        var stateAsync = EventDelegateAsync.Create(() => {
            Console.WriteLine($"State action executing on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(1000);
        });
        stateAsync.InvokeCompleted += (sender, e) => {
            Console.WriteLine($"State action completed. User state: {e.UserState}");
        };
        stateAsync.InvokeAsync(null, "TestState123");
        while (stateAsync.IsBusy)
            continue;   // 阻塞

        // 测试取消
        Console.WriteLine("\nTesting with cenceling ...");
        var cancelAsync = EventDelegateAsync.Create(() => {
            Console.WriteLine($"Cancel action executing on thread {Thread.CurrentThread.ManagedThreadId}");
        });
        cancelAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Cancel action error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Cancel action cancelled");
            } else {
                Console.WriteLine($"Cancel action completed. User state: {e.UserState}");
            }
        };
        cancelAsync.InvokeAsync(null, "TestCancel");
        cancelAsync.CancelAsync("TestCancel");
        while (cancelAsync.IsBusy)
            continue;   // 阻塞

        // 测试异常
        Console.WriteLine("\nTesting with exception ...");
        var exceptionAsync = EventDelegateAsync.Create(() => {
            Console.WriteLine($"Exception action executing on thread {Thread.CurrentThread.ManagedThreadId}");
            throw new Exception("Throw an exception");
        });
        exceptionAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Exception action error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Exception action cancelled");
            } else {
                Console.WriteLine($"Exception action completed. User state: {e.UserState}");
            }
        };
        exceptionAsync.InvokeAsync(null, "TestException");
        while (exceptionAsync.IsBusy)
            continue;   // 阻塞
    }
}