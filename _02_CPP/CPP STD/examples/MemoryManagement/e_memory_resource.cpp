#include <iostream>
#include <memory_resource>
#include <vector>
#include <string>
#include <cassert>
#include <algorithm>

module MemoryManagement;
using namespace std;
using namespace std::pmr;

// 多态分配器示例
void example_basic_polymorphic_allocator() {
	std::cout << "\n=== Basic polymorphic_allocator ===\n";
	// 创建一个单调缓冲区资源作为内存池
	monotonic_buffer_resource pool(1024);
	// 创建使用该内存池的多态分配器
	polymorphic_allocator<int> alloc(&pool);
	int* ptr = alloc.allocate(5);
	for (int i = 0; i < 5; ++i)
		std::construct_at(ptr + i, i * 2);
	std::cout << "Elements: ";
	for (int i = 0; i < 5; ++i)
		std::cout << ptr[i] << " ";
	std::cout << "\n";
	// 销毁对象，释放内存
	for (int i = 0; i < 5; ++i)
		std::destroy_at(ptr + i);
	alloc.deallocate(ptr, 5);
}

// 容器中使用多态分配器
void example_polymorphic_allocator_in_containers() {
	std::cout << "\n=== Polymorphic_allocator in Containers ===\n";
	synchronized_pool_resource pool;
	polymorphic_allocator<char> char_alloc(&pool);	// 为 std::string 类型创建多态分配器
	std::vector<std::string, polymorphic_allocator<std::string>> vec(char_alloc);
	// 向 vector 中添加元素
	vec.push_back("Hello");
	vec.push_back("C++");
	vec.push_back("Memory");
	vec.push_back("Resource");
	for (const auto& str : vec)
		std::cout << str << " ";
	std::cout << "\n";

	// 演示多态分配器的传播性
	vec.reserve(10);  // 会使用相同的分配器分配更多内存
	std::cout << "Vector cap: " << vec.capacity() << "\n";
}

// 不同类型的内存资源
void example_different_memory_resources() {
	std::cout << "\n=== Different memory_resource Types ===\n";

	// 单调缓冲区资源
	unsigned char buffer[1024];
	monotonic_buffer_resource mono_buf(buffer, sizeof(buffer));
	polymorphic_allocator<int> mono_alloc(&mono_buf);
	int* mono_data = mono_alloc.allocate(10);
	mono_alloc.deallocate(mono_data, 10);

	// 同步池资源
	synchronized_pool_resource sync_pool;
	polymorphic_allocator<int> sync_alloc(&sync_pool);
	int* sync_data = sync_alloc.allocate(10);
	sync_alloc.deallocate(sync_data, 10);

	// 非同步池资源
	unsynchronized_pool_resource unsync_pool;
	polymorphic_allocator<int> unsync_alloc(&unsync_pool);
	int* unsync_data = unsync_alloc.allocate(10);
	std::cout << "Elements: ";
	for (int i = 0; i < 10; ++i)
		std::cout << unsync_data[i] << " ";
	std::cout << "\n";
	for (int i = 0; i < 10; ++i)
		std::destroy_at(unsync_data + i);
	unsync_alloc.deallocate(unsync_data, 10);
}

// 自定义内存资源
class CustomMemoryResource : public memory_resource {
private:
	// 简单内存池实现
	static constexpr size_t POOL_SIZE = 1024 * 1024; // 1MB
	unsigned char* pool ;
	size_t next_free = 0;

protected:
	void* do_allocate(size_t bytes, size_t alignment) override {
		if (next_free + bytes > POOL_SIZE) throw std::bad_alloc();
		// 实际应用中应考虑对齐要求
		void* result = pool + next_free;
		next_free += bytes;
		std::cout << "CustomMemoryResource: Allocated " << bytes << " bytes\n";
		return result;
	}
	void do_deallocate(void* p, size_t bytes, size_t alignment) override {
		// 简单实现 - 不实际释放内存（这是一个内存池）
		std::cout << "CustomMemoryResource: Deallocated " << bytes << " bytes\n";
	}
	bool do_is_equal(const memory_resource& other) const noexcept override {
		return this == &other;
	}
public:
	CustomMemoryResource() : next_free(0) {
		std::cout << "CustomMemoryResource: Initialized pool\n";
		pool = new unsigned char[POOL_SIZE];
	}
	~CustomMemoryResource() override {
		delete[POOL_SIZE] pool;
		std::cout << "CustomMemoryResource: Destroyed pool\n";
	}
	void reset() {
		next_free = 0;
		std::cout << "CustomMemoryResource: Reset pool\n";
	}
	size_t used() const {
		return next_free;
	}
};
void example_custom_memory_resource() {
	std::cout << "\n=== Custom memory_resource ===\n";
	// 创建自定义内存资源
	CustomMemoryResource custom_mem;
	polymorphic_allocator<int> custom_alloc(&custom_mem);

	// 分配内存并构造对象
	int* data = custom_alloc.allocate(5);
	for (int i = 0; i < 5; ++i)
		std::construct_at(data + i, i * 10);
	// 使用数据
	std::cout << "Elements: ";
	for (int i = 0; i < 5; ++i)
		std::cout << data[i] << " ";
	std::cout << "\n";
	std::cout << "Memory used: " << custom_mem.used() << " bytes\n";
	// 清理
	for (int i = 0; i < 5; ++i)
		std::destroy_at(data + i);
	custom_alloc.deallocate(data, 5);
	// 重置内存池
	custom_mem.reset();
	std::cout << "Memory reset: " << custom_mem.used() << " bytes\n";
}

// 内存资源的继承与组合
void example_memory_resource_composition() {
	std::cout << "\n=== memory_resource Composition ===\n";
	// 创建一个上游内存资源（作为后备）
	monotonic_buffer_resource upstream;
	// 创建一个使用上游资源的同步池资源
	synchronized_pool_resource pool(&upstream);
	polymorphic_allocator<int> alloc(&pool);
	// 创建使用该分配器的 vector
	std::vector<int, decltype(alloc)> vec(alloc);
	for (int i = 0; i < 20; ++i)
		vec.push_back(i);
	std::cout << "Vector size: " << vec.size() << "\n";
	std::cout << "Vector elements: ";
	for (int i = 0; i < vec.size(); ++i)
		std::cout << vec[i] << " ";
	std::cout << "\n";
}

void test_memory_resource() {
	example_basic_polymorphic_allocator();
	example_polymorphic_allocator_in_containers();
	example_different_memory_resources();
	example_custom_memory_resource();
	example_memory_resource_composition();
}
