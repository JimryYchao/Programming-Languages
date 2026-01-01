#include <iostream>
#include <memory>
#include <vector>
#include <cassert>
#include <atomic>
#include <type_traits>
#include <stdexcept>
#include <algorithm>
#include <functional>
#include <memory_resource>

module MemoryManagement;
using namespace std;

// Smart pointer 
void example_smart_ptr() {
	cout << "\n=== Smart pointer  ===\n";
	// 1 unique_ptr 
	{
		unique_ptr<int> up1 = make_unique<int>(42);
		unique_ptr<int> up2 = move(up1);
		assert(up1 == nullptr);
		unique_ptr<int, void(*)(int*)> up3(new int(100), [](int* p) {
			cout << "Custom deleter: Memory released\n\n";
			delete p;
			});
		unique_ptr<int> up4 = make_unique_for_overwrite<int>();
		*up4 = 200;
		cout << "unique_ptr values: " << *up2 << ", " << *up3 << ", " << *up4 << "\n";
		cout << "unique_ptr values: " << *up2 << ", " << *up3 << ", " << *up4 << "\n\n";
	};
	// 2 shared_ptr
	{
		shared_ptr<std::string> sp1(new std::string("Hello C++23"));
		shared_ptr<std::string> sp2 = sp1;
		cout << "Shared value: " << *sp1 << "\n";
		cout << "Reference count: " << sp1.use_count() << "\n";
		sp1.reset();
		cout << "Reference count after reset sp1: " << sp2.use_count() << "\n";
		shared_ptr<int[]> sp_arr(new int[5] {1, 2, 3, 4, 5});
		cout << "Array elements: ";
		for (int i = 0; i < 5; ++i) cout << sp_arr[i] << " ";
		cout << "\n\n";
	}
	// 3 weak_ptr 
	{
		auto sp = make_shared<int>(100);
		weak_ptr<int> wp = sp;
		cout << "Weak pointer expired status: " << boolalpha << wp.expired() << "\n";
		cout << "Weak pointer lock result: " << wp.lock().use_count() << "\n";
		cout << "Reference count: " << sp.use_count() << "\n";
		sp.reset();
		cout << "Expired status after shared_ptr reset: " << wp.expired() << "\n";
	}
}

// addressof and to_address
void example_address_functions() {
	cout << "\n=== addressof and to_address ===\n";
	int x = 10;
	int* px = &x;
	// addressof 获取对象的原始指针，忽略 & 重载
	int* addr = addressof(x);
	// to_address 从指针，智能指针获取原始指针
	int* to_addr = to_address(px);
	cout << "x value: " << x << ", addressof: " << addr
		<< ", to_address*: " << to_addr << endl;
}

// basic allocator
void example_allocator() {
	cout << "\n=== allocator ===\n";
	allocator<int> alloc;
	// 分配
	int* arr = alloc.allocate(5);
	// 构造
	for (int i = 0; i < 5; ++i) {
		std::construct_at(arr + i, i * 10);
	}
	cout << "Allocated array elements: ";
	for (int i = 0; i < 5; ++i) cout << arr[i] << " ";
	cout << "\n";
	// 销毁
	for (int i = 0; i < 5; ++i) {
		destroy_at(arr + i);
	}
	// 释放
	alloc.deallocate(arr, 5);
	cout << "\n";

	// allocation_result
	auto p = alloc.allocate_at_least(10);
	for (size_t i = 0; i < p.count; i++)
		p.ptr[i] = 10 * i;
	alloc.deallocate(p.ptr, p.count);
}

// uses_allocator
void example_uses_allocator() {
	cout << "=== uses_allocator ===\n";
	struct Value {
		int i;
		Value(int i) : i{ i } { cout << "Value(), i = " << i << endl; }
		~Value() { cout << "~Value(), i = " << i << endl; }
	};
	vector<shared_ptr<Value>> ptr;
	if constexpr (uses_allocator<std::vector<shared_ptr<Value>>, allocator<Value>>::value) {
		byte buffer[sizeof(Value) * 16];
		pmr::monotonic_buffer_resource res(buffer, sizeof(buffer));
		pmr::polymorphic_allocator<Value> alloc(&res);
		for (int i = 0; i < 4; ++i)
			ptr.emplace_back(allocate_shared<Value>(alloc, i));
		for (auto& v : ptr)
			cout << v->i << ",";
		cout << endl;
	}
}

// owner_less 执行 shared weak 之间的指针比较
void example_owner_less() {
	cout << "\n=== owner_less ===\n";
	auto sp1 = make_shared<int>(1);
	auto sp2 = sp1;  // 共享所有权
	auto sp3 = make_shared<int>(2);

	owner_less<shared_ptr<int>> cmp;
	cout << "sp1 < sp3: " << cmp(sp1, sp3) << "\n";  // true
	cout << "sp1 < sp2: " << cmp(sp1, sp2) << "\n"; // false
}

// Atomic smart pointers
void example_atomic_smart_pointers() {
	cout << "\n=== atomic smart ptr ===\n";
	atomic<shared_ptr<int>> atomic_sp;
	atomic<weak_ptr<int>> atomic_wp;

	auto sp = make_shared<int>(100);
	atomic_sp.store(sp);
	atomic_wp.store(sp);

	cout << "Atomic shared_ptr value: " << *atomic_sp.load() << "\n";
	if (!atomic_wp.load().expired())  // false
		cout << "Atomic weak_ptr value: " << *atomic_sp.load() << "\n";
}


// 智能指针适配器与转换
void example_smart_ptr_adapters() {
	cout << "\n=== Smart pointer adapters and conversions ===\n";
	struct Base { virtual int get_value() const { return 0; } virtual ~Base() = default; };
	struct Derived : Base { int get_value() const override { return 42; } };
	shared_ptr<Base> base_ptr = make_shared<Derived>();

	// Static cast
	auto static_cast_ptr = static_pointer_cast<Derived>(base_ptr);
	// Dynamic cast
	auto dynamic_cast_ptr = dynamic_pointer_cast<Derived>(base_ptr);

	cout << "dynamic_cast result: " << dynamic_cast_ptr->get_value() << "\n";
	cout << "static_cast result: " << static_cast_ptr->get_value() << "\n";
}

// Memory management utility functions
void example_memory_management_functions() {
	cout << "\n=== memory management functions ===\n";
	int* ptr = new int[5];

	// Construction and destruction
	construct_at(ptr, 10);
	construct_at(ptr + 1, 20);

	cout << "Values after construction: " << ptr[0] << ", " << ptr[1] << "\n";

	destroy_at(ptr);
	destroy_at(ptr + 1);

	// Destroy range
	uninitialized_fill_n(ptr, 5, 0);
	destroy(ptr, ptr + 5);

	delete[] ptr;
	cout << "\n";
}

// out_ptr & inout_ptr 
void c_style_create_buffer(int** buffer, size_t size) {
	// 模拟 C 风格 API 函数，用于演示 out_ptr
	*buffer = new int[size];
	for (size_t i = 0; i < size; ++i)
		(*buffer)[i] = static_cast<int>(i * 10);
}
void c_style_resize_buffer(int** buffer, size_t old_size, size_t new_size) {
	// 模拟 C 风格 API 函数，用于演示 inout_ptr
	int* new_buffer = new int[new_size];
	// 复制旧数据到新缓冲区
	for (size_t i = 0; i < old_size && i < new_size; ++i)
		new_buffer[i] = (*buffer)[i];
	// 初始化新分配的空间
	for (size_t i = old_size; i < new_size; ++i)
		new_buffer[i] = static_cast<int>(i * 10 + 5);
	// 释放旧缓冲区
	delete[] * buffer;
	// 返回新缓冲区
	*buffer = new_buffer;
}
// out_ptr 会自动管理 C API 分配的内存，避免内存泄漏
void example_out_ptr() {  // out_ptr：用于从 C 风格 API 获取分配的内存
	cout << "\n=== out_ptr ===\n";
	// 使用 out_ptr 将 C 风格 API 返回的原始指针包装为智能指针
	unique_ptr<int[]> ptr;
	c_style_create_buffer(std::out_ptr(ptr), 5);
	cout << "Buffer contents: ";
	for (size_t i = 0; i < 5; ++i)
		cout << ptr[i] << " ";
	cout << "\n";
}
// inout_ptr：用于需要传入并可能被重新分配的内存
void example_inout_ptr() {
	cout << "\n=== inout_ptr ===\n";
	unique_ptr<int[]> ptr(new int[3] {10, 20, 30});
	// 使用 inout_ptr 传递智能指针给可能会重新分配内存的 C API
	// inout_ptr 会处理旧内存的释放和新内存的接管
	c_style_resize_buffer(std::inout_ptr(ptr), 3, 5);
	cout << "After resize: ";
	for (size_t i = 0; i < 5; ++i)
		cout << ptr[i] << " ";
	cout << "\n";
}

// 未初始化存储操作示例
void example_uninitialized_storage() {
	cout << "\n=== uninitialized storage ===\n";

	struct TestObj {
		int value{ 0 };
		TestObj() { cout << "TestObj default constructor called\n"; }
		~TestObj() { cout << "TestObj destructor called, value=" << value << "\n"; }
	};
	const size_t size = 3;
	TestObj* dest = static_cast<TestObj*>(::operator new(sizeof(TestObj) * size));
	// 默认构造对象
	uninitialized_default_construct(dest, dest + size);  // call 3 ctor
	for (size_t i = 0; i < size; ++i)
		dest[i].value = i * 10;
	cout << "uninitialized_default_construct result: ";
	for (size_t i = 0; i < size; ++i)
		cout << dest[i].value << " ";
	cout << "\n";
	// 清理
	destroy(dest, dest + size);  // call 3 ~()
	::operator delete(dest);
}

// 使用分配器进行对象构造的工具函数示例
struct AllocAware {
	int value;
	std::string name;
	// 普通构造函数
	AllocAware(int v, std::string n) : value(v), name(std::move(n)) {
		cout << "AllocAware ctor: value=" << value << ", name=" << name << "\n";
	}
	// 带分配器的构造函数
	template<typename Alloc>
	AllocAware(int v, std::string n, const Alloc&) : value(v), name(std::move(n)) {
		cout << "AllocAware ctor with Alloc: value=" << value << ", name=" << name << "\n";
	}
	~AllocAware() {
		cout << "~AllocAware() \n\n";
	}
};
template<typename Alloc>  // 特化 uses_allocator 以支持 AllocAware 类
struct uses_allocator<AllocAware, Alloc> : true_type {};
void example_allocator_construction() {
	cout << "\n=== allocator construction ===\n";

	// 创建分配器
	std::allocator<AllocAware> alloc;

	// default ctor 创建对象
	{
		auto obj = AllocAware(0, "empty");
	}

	// make_obj_using_allocator 创建对象
	{
		auto obj1 = std::make_obj_using_allocator<AllocAware>(alloc, 42, "object1");
	}

	// uses_allocator_construction_args 准备构造参数
	{
		// 获取构造参数包
		auto args = std::uses_allocator_construction_args<AllocAware>(alloc, 100, "object2");
		// 手动使用参数包构造对象
		AllocAware obj2(std::get<0>(args), get<1>(args));
	}

	// 使用 uninitialized_construct_using_allocator 在未初始化内存上构造对象
	{
		// 分配未初始化内存
		AllocAware* memory = alloc.allocate(1);
		std::uninitialized_construct_using_allocator(memory, alloc, 200, "object3");
		std::destroy_at(memory);
		alloc.deallocate(memory, 1);
	}
}

// 辅助打印函数
void print(const auto& value, const string& label) {
	cout << label << ": " << value << endl;
}

void printArray(const auto& arr, size_t size, const string& label) {
	cout << label << ": ";
	for (size_t i = 0; i < size; ++i) {
		cout << arr[i] << " ";
	}
	cout << endl;
}
// start_lifetime
void example_start_lifetime() {
	cout << "\n=== start_lifetime TODO ===\n";
}

// Test entry function
void test_memory() {
	example_smart_ptr();
	example_address_functions();
	example_allocator();
	example_uses_allocator();
	example_owner_less();
	example_atomic_smart_pointers();
	example_smart_ptr_adapters();
	example_memory_management_functions();
	example_out_ptr();
	example_inout_ptr();
	example_uninitialized_storage();
	example_allocator_construction();

	example_start_lifetime();
}