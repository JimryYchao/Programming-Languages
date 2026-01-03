#include "test.h"
#include <string.h>
#include <stdio.h>	

typedef void (*test_func)(void);
void (_T)(test_func f) { f(); }
char buf[32];

#define T(func) \
	 (memcpy(buf, #func, strlen(#func)), \
 printf("================= Test %s%s\n", (strtok(buf , "_"), strtok(NULL, "_")), ".h ================="), \
	 (_T)(func))

int main()
{
	T(test_errno);

	//T(test_fenv);

	//T(test_inttype);

	//T(test_locale);

	//T(test_setjmp);

	//T(test_signal);

	//T(test_stdalign);

	//T(test_stdatomic);

	//T(test_stdio);

	//T(test_stdlib);

	//T(test_string);

	//T(test_threads);

	//T(test_time);

	//T(test_uchar);
}

