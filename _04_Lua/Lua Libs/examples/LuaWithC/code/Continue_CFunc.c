#include <math.h>
#include "test.h"

static int finishpcall(lua_State* L, int status, intptr_t ctx) {
	printf("Call finishpcall\n");
	(void)ctx;
	status = (status != LUA_OK && status != LUA_YIELD);
	lua_pushboolean(L, (status == 0));
	lua_insert(L, 1);
	return lua_gettop(L);
}

static int luaB_pcall(lua_State* L) {
	printf("Call luaB_pcall\n");
	int status;
	luaL_checkany(L, 1);
	status = lua_pcallk(L, lua_gettop(L) - 1, LUA_MULTRET, 0, 0, finishpcall);
	return finishpcall(L, status, 0);
}

void test_Continue_pcall() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	lua_pushcfunction(L, luaB_pcall);
	lua_setglobal(L, "pcall");

	luaL_dostring(L,
		"co = coroutine.create(function()			\
			print('yield co')						\
			print(pcall(coroutine.yield))			\
			print('resume co')						\
			end)									\
		coroutine.resume(co)");

	luaL_dostring(L, "coroutine.resume(co)");

	lua_close(L);
}


// 548 KFunction
static int counter;
static void _C_DeCounter(lua_State* L, int n);

static int C_DeCounter_yield(lua_State* L, int status, lua_KContext ctx) {
	_C_DeCounter(L, counter);
}

static void _C_DeCounter(lua_State* L, int n) {
	if (!(n < 0)) {
		lua_pushinteger(L, n);
		counter = n - 1;
		lua_yieldk(L, 1, 0, C_DeCounter_yield);
	}
	else
		printf("counter end\n");
}

static int C_DeCounter(lua_State* L) {
	counter = lua_tointeger(L, -1);
	lua_pop(L, 1);
	if (counter <= 0)
		return 0;
	_C_DeCounter(L, counter);
	return 0;
}

void test_Continue_CFunc() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	lua_pushcfunction(L, C_DeCounter);
	lua_setglobal(L, "counter");

	luaL_dostring(L,
		"co = coroutine.create(function()			\
				counter(5)							\
			end)									\
		for i = 5, -1, -1 do						\
			print(coroutine.resume(co))				\
		end");
}