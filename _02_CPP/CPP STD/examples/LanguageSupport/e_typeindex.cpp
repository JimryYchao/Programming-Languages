#include <typeindex>
#include <unordered_map>
#include <iostream>
#include <vector>
#include <memory>
//#include <typeinfo>

module LanguageSupport;
using namespace std;

// type_index 和 typeid( T | expr )
void example_typeid_expr() {
	std::cout << "=== type_index and typeid-expr ===\n";
	// type_index 包装 typeid(T)的返回值类型 const typeinfo&
	std::type_index int_type = typeid(int);
	//const type_info& int_type = typeid(int);
	std::type_index double_type = typeid(double);
	std::type_index vec_type = typeid(std::vector<int>);
	std::cout << "int type name: " << int_type.name() << "\n";
	std::cout << "double type name: " << double_type.name() << "\n";
	std::cout << "vector<int> type name: " << vec_type.name() << "\n";

	// 多态类型和 typeid
	struct Base { virtual ~Base() = default; };
	struct Derived1 : Base {};
	struct Derived2 : Base {};
	std::unique_ptr<Base> b1 = std::make_unique<Derived1>();
	std::unique_ptr<Base> b2 = std::make_unique<Derived2>();
	std::cout << "b1 actual type: " << typeid(*b1).name() << "\n";  
	std::cout << "b2 actual type: " << typeid(*b2).name() << "\n";
}

// 类型工厂模式
void example_type_factory() {
	std::cout << "\n=== Type factory pattern ===\n";

	class Object {
	public:
		virtual ~Object() = default;
		virtual void print() const = 0;
	};

	class A : public Object {
	public:
		void print() const override { std::cout << "Object A\n"; }
	};

	class B : public Object {
	public:
		void print() const override { std::cout << "Object B\n"; }
	};

	std::unordered_map<std::type_index, std::unique_ptr<Object>> factory;
	factory[typeid(A)] = std::make_unique<A>();
	factory[typeid(B)] = std::make_unique<B>();

	auto create = [&](const std::type_index& type) -> Object* {
		auto it = factory.find(type);
		return it != factory.end() ? it->second->print(), it->second.get() : nullptr;
		};

	create(typeid(A));
	create(typeid(B));
}

// bad_cast
struct Foo { virtual ~Foo() {} };
struct Bar { virtual ~Bar() { std::cout << "~Bar\n"; } };
struct Pub : Bar { ~Pub() override { std::cout << "~Pub\n"; } };
void example_bad_cast() {
	std::cout << "\n=== bad type cast ===\n";
	Pub pub;
	try
	{
		(void)dynamic_cast<Bar&>(pub); // OK, upcast
		(void)dynamic_cast<Foo&>(pub); // throws bad_cast
	}
	catch (const std::bad_cast& e)
	{
		std::cout << "e.what(): " << e.what() << '\n';
	}
}

void test_typeindex() {
	example_typeid_expr();
	example_type_factory();
	example_bad_cast();
}