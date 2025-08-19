#include <version>
#include <iostream>
#include <vector>
#include <optional>
#include <filesystem>

module LanguageSupport;

// 检查C++标准版本
void example_check_standard_version() {
    std::cout << "=== C++ Standard Version ===\n";

#ifdef __cplusplus
    std::cout << "__cplusplus: " << __cplusplus << "\n";
#endif

    // 检查特定C++版本特性
#if __cpp_lib_optional >= 201606L
    std::cout << "std::optional is available (C++17)\n";
#endif

#if __cpp_lib_filesystem >= 201703L
    std::cout << "std::filesystem is available (C++17)\n";
#endif

#if __cpp_lib_three_way_comparison >= 201907L
    std::cout << "Three-way comparison is available (C++20)\n";
#endif
}

// 检查编译器特定宏
void example_check_compiler_macros() {
    std::cout << "\n=== Compiler Identification ===\n";

#ifdef __GNUC__
    std::cout << "GCC/G++ compiler detected, version: " << __GNUC__ << "." << __GNUC_MINOR__ << "\n";
#endif

#ifdef _MSC_VER
    std::cout << "MSVC compiler detected, version: " << _MSC_VER << "\n";
#endif

#ifdef __clang__
    std::cout << "Clang compiler detected, version: " << __clang_major__ << "." << __clang_minor__ << "\n";
#endif
}

// 特性测试宏
void example_feature_test_macros() {
    std::cout << "\n=== Feature Test Macros ===\n";

#if __has_include(<optional>)
    std::cout << "<optional> header is available\n";
#else
    std::cout << "<optional> header is not available\n";
#endif

#if __has_cpp_attribute(deprecated)
    std::cout << "[[deprecated]] attribute is supported\n";
#endif

#if __cpp_lib_ranges >= 201911L
    std::cout << "Ranges library is available (C++20)\n";
#endif
}

// 条件编译基于可用特性
void example_conditional_compilation() {
    std::cout << "\n=== Conditional Compilation ===\n";

#if __cpp_lib_string_view >= 201606L
    std::cout << "Using std::string_view (modern C++)\n";
#else
    std::cout << "Using const std::string& (fallback)\n";
#endif

#if __cpp_concepts >= 201907L
    std::cout << "Compiling with concepts support\n";
#else
    std::cout << "Compiling without concepts support\n";
#endif
}

// 检查平台特性
void example_platform_checks() {
    std::cout << "\n=== Platform Checks ===\n";

#ifdef _WIN32
    std::cout << "Running on Windows platform\n";
#elif __linux__
    std::cout << "Running on Linux platform\n";
#elif __APPLE__
    std::cout << "Running on macOS platform\n";
#else
    std::cout << "Running on unknown platform\n";
#endif

#ifdef __LP64__
    std::cout << "64-bit platform detected\n";
#else
    std::cout << "32-bit platform detected\n";
#endif
}

void test_version() {
    example_check_standard_version();
    example_check_compiler_macros();
    example_feature_test_macros();
    example_conditional_compilation();
    example_platform_checks();
}