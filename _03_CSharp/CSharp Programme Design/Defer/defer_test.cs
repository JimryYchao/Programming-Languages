
public class DeferTest
{
    public static async void Test()
    {
        using var defer = Defer.New(() => Console.WriteLine(1));
        defer.Run(() => Console.WriteLine(2));
        Console.WriteLine("Hello");
        defer.Run(() => Console.WriteLine(3));
        using (defer.Run(() => Console.WriteLine(4)))
        {
            //throw new Exception();   // Hello, Exception, 4,3,2,1
            // ...
        } // >> 提前 defer
        defer.Run(() => Console.WriteLine(5));
        Console.WriteLine("World");
        // Hello, 4,3,2,1, World ,5
    }
}