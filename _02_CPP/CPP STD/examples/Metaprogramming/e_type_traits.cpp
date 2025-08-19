#include <iostream>
#include <type_traits>
#include <string>
#include <vector>
#include <array>

module Metaprogramming;
using namespace std;

void example_type_traits() {
    std::cout << "\n=== Type Traits ===\n";
    
    // 编译时类型信息
    std::cout << "int is integral: " << std::boolalpha << std::is_integral<int>::value << std::endl;
    std::cout << "double is floating point: " << std::is_floating_point<double>::value << std::endl;
    std::cout << "std::string is a class: " << std::is_class<std::string>::value << std::endl;
    std::cout << "int& is a reference: " << std::is_reference<int&>::value << std::endl;
    std::cout << "int* is a pointer: " << std::is_pointer<int*>::value << std::endl;
    std::cout << "const int is const: " << std::is_const<const int>::value << std::endl;
    
    // CV 限定符 remove / add
    std::cout << "Remove const: const int -> " << typeid(typename std::remove_const<const int>::type).name() << std::endl;
    std::cout << "Add const: int -> " << typeid(typename std::add_const<int>::type).name() << std::endl;
    std::cout << "Remove reference: int& -> " << typeid(typename std::remove_reference<int&>::type).name() << std::endl;
    std::cout << "Add pointer: int -> " << typeid(typename std::add_pointer<int>::type).name() << std::endl;
    
    // 类型关系
    std::cout << "int is same as int: " << std::is_same<int, int>::value << std::endl;
    std::cout << "int is same as double: " << std::is_same<int, double>::value << std::endl;
    struct Base {};
    struct Derived : Base {};
    std::cout << "Derived is derived from Base: " << std::is_base_of<Base, Derived>::value << std::endl;
    std::cout << "Base is derived from Derived: " << std::is_base_of<Derived, Base>::value << std::endl;
    std::cout << "int is convertible to double: " << std::is_convertible<int, double>::value << std::endl;   // 可转换
    std::cout << "const char* is convertible to std::string: " << std::is_convertible<const char*, std::string>::value << std::endl;
}

// Conditional type selection
template<typename T>
void print_type_info() {
    using EffectiveType = typename std::conditional<std::is_arithmetic<T>::value, T, std::string>::type;
    std::cout << "For type " << typeid(T).name() << ", effective type is: " << typeid(EffectiveType).name() << std::endl;
}
void example_conditional_types() {
    std::cout << "\n=== Conditional Type Selection ===\n";
    
    // Select type based on condition
    using Type1 = typename std::conditional<true, int, double>::type;
    using Type2 = typename std::conditional<false, int, double>::type;
    std::cout << "Select int when condition is true: " << typeid(Type1).name() << std::endl;
    std::cout << "Select double when condition is false: " << typeid(Type2).name() << std::endl;
    
    // Apply in function template
    print_type_info<int>();
    print_type_info<string>();
    print_type_info<vector<int>>();
}

void test_type_traits() {
    example_type_traits();
    example_conditional_types();
}