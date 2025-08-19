#include <cstring>
#include <cstdio>
using namespace std;

// === test modules ====
import LanguageSupport;
import Diagnostics;
import MemoryManagement;
import Metaprogramming;

// ==== end test modules ====

typedef void (*test_func)(void);
void (_T)(test_func f) { f(); }
char buf[32];

#define T(func) \
	 (memcpy(buf, #func, strlen(#func)), \
 printf("============ Test Header <%s>%s\n", buf+5, " ============"), \
	 (_T)(func))

int main(void) {
	/*     language support      */
	//T(test_compare);
	//T(test_coroutine);
	T(test_concepts);
	//T(test_cstddef);
	//T(test_exception);
	//T(test_initializer_list);
	//T(test_limits);
	//T(test_new);
	//T(test_source_location);
	//T(test_typeindex);
	//T(test_version);

	/*     diagnostics      */
	//T(test_stacktrace);
	//T(test_system_error);

	/*     memory management      */
	//T(test_memory);
	//T(test_memory_resource);
	//T(test_scoped_allocator);

	/*     metaprogramming      */
	//T(test_ratio);
	//T(test_type_traits);

	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */
	/*     language support      */


	return 0;
}

