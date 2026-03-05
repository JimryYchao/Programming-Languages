#include "AsyncIO.hpp"
#include "Task.hpp"

// 模拟异步读取
AsyncIO<std::string> async_read_file(const std::string& filename) {
    AsyncIO<std::string> op;
    // 模拟异步操作
    EventLoop::instance().post([&op, filename]() {
        // 模拟文件读取延迟
        std::cout << "读取文件: " << filename << std::endl;
        op.complete("文件内容: " + filename);
        });

    return op;
}

// 模拟异步写入
AsyncIO<bool> async_write_file(const std::string& filename, const std::string& content) {
    AsyncIO<bool> op;
    EventLoop::instance().post([&op, filename, content]() {
        std::cout << "写入文件: " << filename << std::endl;
        std::cout << "内容: " << content << std::endl;
        op.complete(true);
        });

    return op;


}

// 使用协程进行异步文件操作
Task<std::string> process_files() {
    auto content1 = co_await async_read_file("input1.txt");
    auto content2 = co_await async_read_file("input2.txt");
    std::string combined = content1 + " + " + content2;
    co_await async_write_file("output.txt", combined);
    co_return combined;
}

int main() {
    auto task = process_files();
    if (!task.done()) 
        task.handle.resume();
    // 运行事件循环
    EventLoop::instance().run();
    std::cout << "最终结果: " << task.get() << std::endl;
}