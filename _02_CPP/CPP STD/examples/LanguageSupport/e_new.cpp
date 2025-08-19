#include <new>
#include <iostream>
#include <cstdlib>
#include <vector>
#include <mutex>
#include <cstdint>
#include <atomic>
#include <chrono>
#include <cstddef>
#include <iomanip>
#include <thread>
#include <cassert>

module LanguageSupport;
using namespace std;

// 定位 new 操作符
void example_placement_new()
{
	std::cout << "\n=== Placement new ===\n";

	// 预分配内存
	alignas(std::string) char buffer[sizeof(std::string)];
	// 在预分配内存上构造对象
	std::string* pstr = new (buffer) std::string("Hello, placement new!");
	std::cout << *pstr << "\n";
	// 显式调用析构函数
	pstr->~basic_string();
}

// 自定义 new_handler
void example_new_handler()
{
	std::cout << "\n=== Custom new_handler ===\n";

	// 保存原来的 new_handler (get_new_handler)
	std::new_handler old_handler = std::set_new_handler([]() {
		std::cout << "Custom new_handler called - out of memory!\n";
		throw std::bad_alloc(); });

	try
	{
		// 尝试分配过多内存
		std::size_t huge_size = std::numeric_limits<std::size_t>::max();
		auto* p = new double[huge_size];
		delete[] p;
	}
	catch (const std::bad_alloc& e)
	{
		std::cout << "Caught bad_alloc: " << e.what() << "\n";
	}
	catch (...) {}
	// 恢复原来的 new_handler
	std::set_new_handler(old_handler);
}

// nothrow new
void example_nothrow_new() noexcept
{
	std::cout << "\n=== Nothrow new ===\n";

	// 使用 nothrow 版本的 new
	size_t size = std::numeric_limits<std::size_t>::max();
	auto* p = new (std::nothrow) double[size];
	if (p == nullptr)
		std::cout << "Allocation failed (returned nullptr)\n";
	else
	{
		std::cout << "Allocation succeeded\n";
		delete[] p;
	}
}

// 对齐内存分配
void example_aligned_new()
{
	std::cout << "\n=== Aligned new ===\n";

	// 分配 64 字节对齐的内存
	constexpr std::size_t alignment = 64;
	constexpr std::size_t size = 1024;
	void* ptr = ::operator new(size, std::align_val_t{ alignment });
	std::cout << "Allocated memory at: " << ptr << "\n";
	std::cout << "Alignment: " << alignment << "\n";

	// 检查对齐
	if (reinterpret_cast<std::uintptr_t>(ptr) % alignment == 0)
		std::cout << "Memory is properly aligned\n";
	else
		std::cout << "Memory is NOT properly aligned\n";
	::operator delete(ptr, size, std::align_val_t{ alignment });
}

// 自定义分配器使用 std::nothrow
template<class T>
struct NothrowAllocator
{
	using value_type = T;
	NothrowAllocator() = default;
	template<class U>
	constexpr NothrowAllocator(const NothrowAllocator <U>&) noexcept {}
	[[nodiscard]] T* allocate(std::size_t n) noexcept {
		return report(new (std::nothrow) T[n], n);
	}
	void deallocate(T* p, std::size_t n) noexcept {
		delete[] report(p, n, false);
	}
private:
	T* report(T* p, std::size_t n, bool alloc = true) const noexcept {
		if (p)
			std::cout << (alloc ? "Alloc: " : "Dealloc: ") << sizeof(T) * n
			<< " bytes at " << std::hex << std::showbase
			<< reinterpret_cast<void*>(p) << std::dec << " in NothrowAllocator\n";
		else cout << (alloc ? "Alloc " : "Dealloc ") << "nullptr in NothrowAllocator\n";
		return p;
	}
};
template<class T>
struct Mallocator
{
	using value_type = T;
	Mallocator() = default;
	template<class U>
	constexpr Mallocator(const Mallocator <U>&) noexcept {}

	[[nodiscard]] T* allocate(std::size_t n)
	{
		if (n > std::numeric_limits<std::size_t>::max() / sizeof(T))
			throw std::bad_array_new_length();
		if (auto p = static_cast<T*>(std::malloc(n * sizeof(T))))
		{
			report(p, n);
			return p;
		}
		throw std::bad_alloc();
	}
	void deallocate(T* p, std::size_t n) noexcept
	{
		report(p, n, 0);
		std::free(p);
	}
private:
	void report(T* p, std::size_t n, bool alloc = true) const
	{
		std::cout << (alloc ? "Alloc: " : "Dealloc: ") << sizeof(T) * n
			<< " bytes at " << std::hex << std::showbase
			<< reinterpret_cast<void*>(p) << std::dec << " in Mallocator\n";
	}
};
void example_custom_allocator()
{
	std::cout << "\n=== Custom allocator with nothrow ===\n";
	// using Mallocator
	std::vector<int, Mallocator<int>> v(8);
	v.push_back(42);
	// using NothrowAllocator
	try
	{
		std::vector<int, NothrowAllocator<int>> vec;
		// 尝试分配大量内存
		vec.reserve(vec.max_size()); // max_size 4611686018427387903
	}
	catch (const std::bad_alloc& e)
	{
		std::cout << "Caught bad_alloc: " << e.what() << "\n";
	}
}

// 演示优化 CPU 缓冲行方式，提升多线程和缓冲性能 [href=https://zh.cppreference.com/w/cpp/thread/hardware_destructive_interference_size.html]
static std::mutex cout_mutex;
constexpr int max_write_iterations{ 10'000'000 }; // the benchmark time tuning
inline auto now() noexcept { return std::chrono::high_resolution_clock::now(); }
struct alignas(std::hardware_constructive_interference_size) OneCacheLiner
{	// 高频访问数据结构优化，数据成员存储在同一缓冲行以减少缓冲加载次数
	std::atomic_uint64_t x{};
	std::atomic_uint64_t y{};
} oneCacheLiner;
struct TwoCacheLiner
{	// 避免伪共享，确保两个变量不会因存储在同一缓冲行而导致不必要的缓冲失效
	alignas(std::hardware_destructive_interference_size) std::atomic_uint64_t x{};
	alignas(std::hardware_destructive_interference_size) std::atomic_uint64_t y{};
} twoCacheLiner;
template<bool xy, class Liner> void Thread_CacheLiner(Liner& CacheLiner)
{
	const auto start{ now() };
	for (uint64_t count{}; count != max_write_iterations; ++count)
		if constexpr (xy)
			CacheLiner.x.fetch_add(1, std::memory_order_relaxed);
		else
			CacheLiner.y.fetch_add(1, std::memory_order_relaxed);
	const std::chrono::duration<double, std::milli> elapsed{ now() - start };
	std::lock_guard lk{ cout_mutex };
	std::cout << (typeid(CacheLiner) == typeid(oneCacheLiner) ? "one" : "two") << "CacheLinerThread() spent " << elapsed.count() << " ms\n";
	if constexpr (xy) CacheLiner.x = elapsed.count();
	else CacheLiner.y = elapsed.count();
}

void example_optimize_CPU_Cache()
{
	cout << "\n=== Optimize CPU cache lines ===\n";
	std::cout << "hardware_destructive_interference_size == "
		<< std::hardware_destructive_interference_size << '\n'
		<< "hardware_constructive_interference_size == "
		<< std::hardware_constructive_interference_size << "\n\n"
		<< std::fixed << std::setprecision(2)
		<< "sizeof( OneCacheLiner ) == " << sizeof(OneCacheLiner) << '\n'
		<< "sizeof( TwoCacheLiner ) == " << sizeof(TwoCacheLiner) << "\n\n";
	constexpr int max_runs = 4;

	int oneCacheLiner_average{ 0 };
	// 伪共享问题，不同线程修改同一缓冲行不同变量时，可能会导致缓冲一致性协议（MESI），
	// 强制刷新缓冲而降低性能。但是提高了缓冲命中率，同时访问这些变量可以减少缓冲加载次数，提高性能。
	for (auto i{ 0 }; i != max_runs; ++i)
	{
		std::thread th1{ [] { Thread_CacheLiner<0>(oneCacheLiner); } };
		std::thread th2{ [] { Thread_CacheLiner<1>(oneCacheLiner); } };
		th1.join();
		th2.join();
		oneCacheLiner_average += oneCacheLiner.x + oneCacheLiner.y;
	}
	std::cout << "Average T1 time: "
		<< (oneCacheLiner_average / max_runs / 2) << " ms\n\n";

	// 避免伪共享，提高多线程性能，利用 hardware_destructive_interference_size 真共享机制
	// 将不同的变量强制放在不同的缓冲行，避免伪共享的多线程修改导致的强制刷新（MESI）
	int twoCacheLiner_average{ 0 };    // 空间换时间
	for (auto i{ 0 }; i != max_runs; ++i)
	{
		std::thread th1{ [] { Thread_CacheLiner<0>(twoCacheLiner); } };
		std::thread th2{ [] { Thread_CacheLiner<1>(twoCacheLiner); } };
		th1.join();
		th2.join();
		twoCacheLiner_average += twoCacheLiner.x + twoCacheLiner.y;
	}

	std::cout << "Average T2 time: "
		<< (twoCacheLiner_average / max_runs / 2) << " ms\n\n"
		<< "Ratio T1/T2:~ "
		<< 1.0 * oneCacheLiner_average / twoCacheLiner_average << '\n';
}

// std::launder 对象存储重用
void example_storage_reuse() {
	std::cout << "\n=== Object storage reuse ===\n";
	struct X {
		const int n;
	};
	X* px = new X{ 10086 };  // 10086
	// px->n = 10010; // const 无法修改
	std::cout << "Original value: " << px->n << "\n";
	px->~X();

	struct Y {
		int n;
	};
	auto py = std::launder(new (px) Y{ 10010 });   // 重用对象存储
	std::cout << "New value after reuse: " << py->n << "\n";
	py->~Y();
}

// std::launder 与类型转换
void example_launder_with_reinterpret_cast() {
	std::cout << "\n=== std::launder with reinterpret_cast ===\n";

	struct A { int x; };
	struct B { int y; };

	A a{ 42 };	// 创建 A 对象
	B* b = std::launder(reinterpret_cast<B*>(&a));  // &a -> B*
	std::cout << "B y (actually A x): " << b->y << "\n";
}

// 使用 std::launder 处理指向基类的指针
struct Base
{
	virtual int transmogrify();
};

struct Derived : Base
{
	int transmogrify() override
	{
		new(this) Base;
		return 2;
	}
};
int Base::transmogrify()
{
	new(this) Derived;
	return 1;
}
static_assert(sizeof(Derived) == sizeof(Base));
void example_launder_with_inheritance()
{
	std::cout << "\n=== std::launder with inheritance ===\n";

	// case 1：`new(this) Derived;` 新对象未能实现透明替换（Base-> Derived），它是一个基类子对象
	// 旧对象 base 仍是一个完整对象。-------------
	Base base;								// |
	int n = base.transmogrify();    // 1	// ⬇
	// int m = base.transmogrify(); // undefined behavior，编译器仍已依据旧的内存模型进行操作, m = 1
	int m = std::launder(&base)->transmogrify(); // launder 绕过编译器对旧对象原有的内存假设, 去虚化栅栏
	assert(m + n == 3);

	// case 2：处理指向基类的指针
	struct Base { int x = 10; virtual ~Base() = default; };
	struct Derived : Base { int y = 20; };
	Derived* d = new Derived;
	Base* b = d;
	b->x = 10086;
	d->~Derived();
	new (d) Derived;  // 重用 d 存储
	// 必须使用 std::launder 访问新的对象
	b = std::launder(b);
	std::cout << "Base x: " << b->x << "\n";
	delete d;
}

void test_new(void)
{
	example_placement_new();
	example_new_handler();
	example_nothrow_new();
	example_aligned_new();
	example_custom_allocator();
	example_optimize_CPU_Cache();
	// launder
	example_storage_reuse();
	example_launder_with_reinterpret_cast();
	example_launder_with_inheritance();
}