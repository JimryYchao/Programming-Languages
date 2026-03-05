
/// <summary>
/// defer.Do(action) like a go defer.
/// <code>
/// {
///     using var defer = Defer.New(ac);
///     ...
///     defer.Do(ac2)
///     ...
/// } // >>> defer do ac2, ac1 ... ac
/// </code>
/// </summary>
public sealed class Defer : IDisposable
{
    static Stack<Defer> deferPool = new Stack<Defer>();
    private Stack<Action> deferStack = new Stack<Action>();
    public static Defer New(Action action)
    {
        Defer defer;
        lock (deferPool)
        {
            if (deferPool.TryPop(out Defer d))
                defer = d;
            else defer = new Defer();
        }
        defer.Run(action);
        return defer;
    }
    public Defer Run(Action action)
    {
        if (action is not null)
            this.deferStack.Push(action);
        return this;
    }
    void IDisposable.Dispose()
    {
        // 异常时跳过并继续
        if (this.deferStack.Count > 0)
        {
            using (this)
            {
                try
                {
                    for (int i = 0; i < deferStack.Count;)
                        this.deferStack.Pop().Invoke();
                }
                catch { }
            }
        }
        lock (deferPool)
            Defer.deferPool.Push(this);
    }
}