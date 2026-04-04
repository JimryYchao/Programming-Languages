

public class DelegateAsyncTest {
    public static void Test() {
        List<WaitHandle> handles = new();
        for (int i = 0; i < 10; i++) {
            var ac = DelegateAsync.Create<Func<int, int>, int>((i) => {
                Thread.Sleep(new Random().Next(50, 200));
                return i;
            });
            handles.Add(ac.BeginInvoke([i], (asRT) => {
                var rt = ac.EndInvoke(asRT);
                Console.WriteLine(rt);
            }).AsyncWaitHandle);
        }
        // 等待全部异步完成
        WaitHandle.WaitAll([.. handles]);
    }
}