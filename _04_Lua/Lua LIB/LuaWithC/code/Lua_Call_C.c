#include "test.h"

#include <math.h>
#include <string.h>

static int l_sin(lua_State* L)
{
	double d = luaL_checknumber(L, 1);
	lua_pushnumber(L, sin(d));
	return 1;
}

// Split(str, sep)
static int l_split(lua_State* L) {
	const char* s = luaL_checkstring(L, 1);
	const char* sep = luaL_checkstring(L, 2);
	const char* e;
	int i = 1;

	lua_newtable(L);  // result

	while ((e = strstr(s, sep)) != NULL) {
		lua_pushlstring(L, s, e - s);
		lua_rawseti(L, -2, i++);
		s = e + 1;
	}

	lua_pushstring(L, s);
	lua_rawseti(L, -2, i);
	return 1;
}

// printTable(table)
static int l_printTable(lua_State* L) {
	luaL_checktype(L, 1, LUA_TTABLE);
	int len = lua_rawlen(L, -1);
	for (int i = 1; i <= len; i++)
		if (lua_geti(L, 1, i)) {
			printf("%s,", lua_tostring(L, -1));
		}
	printf("\n");
	return 0;
}


static void test_call_c()
{
	// 创建状态机
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);

	// 压入 C 函数并执行
	lua_pushcfunction(L, l_sin);
	lua_setglobal(L, "l_sin");

	luaL_dostring(L, "return l_sin(math.pi  / 6)");
	luaL_dostring(L, "return l_sin(math.pi  / 2)");

	double rt = lua_tonumber(L, -2);
	printf("call l_sin: rt = %g\n", rt);

	rt = lua_tonumber(L, -1);
	printf("call l_sin: rt = %g\n", rt);

	lua_close(L);
}

// 等效于 test_call_c
void test_Lua_Call_C() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);

	lua_register(L, "Sin", l_sin);  // 注册 C 函数 Sin

	// do file
	if (luaL_loadfile(L, "LuaCallC.lua") || lua_pcall(L, 0, 0, 0))
		luaL_error(L, "cannot load lua file: %s\n", lua_tostring(L, -1));

	// do luastring
	lua_register(L, "Split", l_split);
	lua_register(L, "printTable", l_printTable);
	luaL_dostring(L, "local rt = Split('test split a long string', ' '); printTable(rt)");

	lua_close(L);
}