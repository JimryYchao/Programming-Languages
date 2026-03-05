#include "Task.hpp"

Task<int> async_add(int a, int b)
{
    std::cout << "开始计算 " << a << " + " << b << std::endl;
    co_return a + b;
}
Task<int> async_computation()
{
    int x = co_await async_add(1, 2);
    int y = co_await async_add(x, 3);
    co_return y * 2;
}

int main()
{
    auto task = async_computation();
    int result = task.get();
    std::cout << "结果: " << result << std::endl; // 输出: 12
}
