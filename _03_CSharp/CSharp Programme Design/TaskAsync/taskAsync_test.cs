using System.Threading.Tasks;

public class TaskAsyncTest {
    static IEnumerable<TaskAsync<Task>> Generator() {
        int v = 0;
        for (int i = 0; i < 10; i++)
            yield return TaskAsyncFactory.Run(async () => {
                var rt = v++;
                await Task.Delay(new Random().Next(50, 200));
                Console.WriteLine(rt);
            });
    }
    public static async Task Test() {
        await Task.WhenAll(from task in Generator() select task);
    }
}
