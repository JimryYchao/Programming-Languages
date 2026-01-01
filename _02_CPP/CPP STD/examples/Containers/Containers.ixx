#include <string>

export module Containers;


export{
	// 序列容器
	void test_array();
	void test_vector();
	void test_deque();
	void test_list();
	void test_forward_list();
    // 排序容器
	void test_map();
	void test_set();
	// 
	void test_unordered_map();
	void test_unordered_set();
	// adaptors
	void test_stack();
	void test_queue();
	void test_flat_map();
	void test_flat_set();
	// view
	void test_span();
	void test_mdspan();

}