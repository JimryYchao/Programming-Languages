#include <cstring>
#include <cstdio>
using namespace std;

// === test modules ====
import LanguageSupport;
import Diagnostics;
import MemoryManagement;
import Metaprogramming;
import GeneralUtilities;
import Strings;
import Containers;

// ==== end test modules ====

typedef void (*test_func)(void);
void (_T)(test_func f) { f(); }
char buf[32];

#define T(func) \
	 (memcpy(buf, #func, strlen(#func)), \
 printf("============ Test Header <%s>%s\n", buf+5, " ============"), \
	 (_T)(func))

int main(void) {
	/*     Language support      */
	//T(test_compare);
	//T(test_coroutine);
	//T(test_concepts);
	//T(test_cstddef);
	//T(test_exception);
	//T(test_initializer_list);
	//T(test_limits);
	//T(test_new);
	//T(test_source_location);
	//T(test_typeindex);
	//T(test_version);

	/*     Diagnostics      */
	//T(test_stacktrace);
	//T(test_system_error);

	/*     Memory management      */
	//T(test_memory);
	//T(test_memory_resource);
	//T(test_scoped_allocator);

	/*     Metaprogramming      */
	//T(test_ratio);
	//T(test_type_traits);

	/*     General utilities      */
	//T(test_any);
	//T(test_bit);
	//T(test_bitset);
	//T(test_execution);
	//T(test_expected);
	//T(test_functional);
	//T(test_optional);
	//T(test_tuple);
	//T(test_utility);
	//T(test_variant);

	/*     Strings      */
	//T(test_string);
	//T(test_string_view);

	/*     Containers      */
	//T(test_array);
	//T(test_vector);
	//T(test_deque);
	//T(test_list);
	//T(test_forward_list);
	T(test_stack);
	//T(test_queue);
	//T(test_set);
	//T(test_map);
	//T(test_flat_set);
	//T(test_flat_map);
	//T(test_unordered_set);
	//T(test_unordered_map);

	//T(test_span);
	

	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */

	return 0;
}

