#include <iostream>
#include <vector>
#include <array>
#include <mdspan>
#include <algorithm>
#include <numeric>

module Containers;
using namespace std;

// Basic usage of std::mdspan
void example_mdspan_basic() {
    std::cout << "\n=== Basic std::mdspan Usage ===\n";

    // Create a data buffer for the mdspan
    std::vector<int> data(12);
    std::iota(data.begin(), data.end(), 1); // Fill with 1, 2, 3, ..., 12
    
    // Create a 2D mdspan (3 rows x 4 columns)
    std::mdspan<int, 3, 4> ms2d(data.data());
    
    // Access elements using operator()
    std::cout << "2D mdspan contents: " << std::endl;
    for (size_t i = 0; i < ms2d.extent(0); ++i) {
        for (size_t j = 0; j < ms2d.extent(1); ++j) {
            std::cout << ms2d(i, j) << "\t";
        }
        std::cout << std::endl;
    }
    
    // Modify elements
    ms2d(1, 2) = 100; // Change value at row 1, column 2
    std::cout << "After modification, ms2d(1, 2) = " << ms2d(1, 2) << std::endl;
    
    // Create a 3D mdspan (2 x 3 x 2)
    std::array<double, 12> double_data;
    for (size_t i = 0; i < double_data.size(); ++i) {
        double_data[i] = i * 1.5;
    }
    
    std::mdspan<double, 2, 3, 2> ms3d(double_data.data());
    
    // Access 3D elements
    std::cout << "\n3D mdspan element at (1, 1, 1): " << ms3d(1, 1, 1) << std::endl;
    
    // Use dynamic extents
    std::mdspan<int, std::dynamic_extent, std::dynamic_extent> dyn_ms(
        data.data(), 3, 4
    );
    
    std::cout << "\nDynamic mdspan dimensions: (" << dyn_ms.extent(0) 
              << ", " << dyn_ms.extent(1) << ")" << std::endl;
    
    // Check mdspan properties
    std::cout << "ms2d rank (number of dimensions): " << ms2d.rank() << std::endl;
    std::cout << "ms3d rank: " << ms3d.rank() << std::endl;
    std::cout << "ms2d is exhaustive (covers entire data): " 
              << std::boolalpha << ms2d.is_exhaustive() << std::endl;
    std::cout << "ms2d is unique (no overlapping elements): " 
              << std::boolalpha << ms2d.is_unique() << std::endl;
}

// Advanced usage of std::mdspan
void example_mdspan_advanced() {
    std::cout << "\n=== Advanced std::mdspan Usage ===\n";
    
    // Create a larger data buffer
    std::vector<int> large_data(20);
    std::iota(large_data.begin(), large_data.end(), 100);
    
    // Create a submdspan (view of a portion of the data)
    std::mdspan<int, std::dynamic_extent, std::dynamic_extent> full_ms(
        large_data.data(), 5, 4
    );
    
    // Create a subspan - rows 1-3, columns 0-2
    auto sub_ms = std::submdspan(full_ms, std::tuple(1, 4), std::tuple(0, 3));
    
    std::cout << "Submdspan contents: " << std::endl;
    for (size_t i = 0; i < sub_ms.extent(0); ++i) {
        for (size_t j = 0; j < sub_ms.extent(1); ++j) {
            std::cout << sub_ms(i, j) << "\t";
        }
        std::cout << std::endl;
    }
    
    // Custom layout mapping (row-major vs column-major)
    using row_major_layout = std::layout_right;
    using col_major_layout = std::layout_left;
    
    std::mdspan<int, 3, 4, row_major_layout> row_major_ms(large_data.data());
    std::mdspan<int, 3, 4, col_major_layout> col_major_ms(large_data.data());
    
    std::cout << "\nRow-major vs Column-major layout: " << std::endl;
    std::cout << "row_major_ms(1, 2) = " << row_major_ms(1, 2) << std::endl;
    std::cout << "col_major_ms(1, 2) = " << col_major_ms(1, 2) << std::endl;
    
    // Using mdspan with mathematical operations
    std::vector<double> matrix_a(6, 1.0);  // 2x3 matrix of 1.0s
    std::vector<double> matrix_b(6, 2.0);  // 2x3 matrix of 2.0s
    std::vector<double> result(6, 0.0);    // Result matrix
    
    auto ms_a = std::mdspan<double, 2, 3>(matrix_a.data());
    auto ms_b = std::mdspan<double, 2, 3>(matrix_b.data());
    auto ms_result = std::mdspan<double, 2, 3>(result.data());
    
    // Element-wise addition
    for (size_t i = 0; i < ms_result.extent(0); ++i) {
        for (size_t j = 0; j < ms_result.extent(1); ++j) {
            ms_result(i, j) = ms_a(i, j) + ms_b(i, j);
        }
    }
    
    std::cout << "\nResult of matrix addition: " << std::endl;
    for (size_t i = 0; i < ms_result.extent(0); ++i) {
        for (size_t j = 0; j < ms_result.extent(1); ++j) {
            std::cout << ms_result(i, j) << "\t";
        }
        std::cout << std::endl;
    }
    
    // Multi-dimensional array traversal with structured bindings
    std::cout << "\nTraversing with structured bindings: " << std::endl;
    auto print_element = [](auto&&... indices) {
        std::cout << "Element at (" << (... << (indices << ", ")) << "\b\b): " 
                  << ms_a(indices...) << std::endl;
    };
    
    // Visit specific elements
    print_element(0, 0);
    print_element(1, 2);
    
    // Example of using mdspan with non-contiguous data (strided access)
    // Note: This is a more advanced use case and requires C++23 or later
    std::cout << "\nNote: Strided access patterns can be achieved with custom layout mappings" << std::endl;
    std::cout << "in C++23 and later, allowing for more complex data access patterns" << std::endl;
}

void test_mdspan() {
    example_mdspan_basic();
    example_mdspan_advanced();
}