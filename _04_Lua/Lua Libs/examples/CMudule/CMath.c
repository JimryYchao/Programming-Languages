// #pragma comment(lib, "LuaDll.lib") 
#include "include\lua.h"
#include "include\lauxlib.h"

static int l_add(lua_State* L) {
	lua_pushnumber(L, luaL_checknumber(L, 1) + luaL_checknumber(L, 2));
	return 1;
}

static int l_mul(lua_State* L) {
	lua_pushnumber(L, luaL_checknumber(L, 1) * luaL_checknumber(L, 2));
	return 1;
}

static int l_sub(lua_State* L) {
	lua_pushnumber(L, luaL_checknumber(L, 1) - luaL_checknumber(L, 2));
	return 1;
}

static int l_div(lua_State* L) {
	lua_pushnumber(L, luaL_checknumber(L, 1) / luaL_checknumber(L, 2));
	return 1;
}

// 注册 C 库函数
static const struct luaL_Reg CMathlib[] = {
	{"add", l_add},
	{"mul", l_mul},
	{"sub", l_sub},
	{"div", l_div},
	{NULL, NULL}, // 哨兵, 标识数组结尾
};

int luaopen_CMath(lua_State* L) {
	luaL_newlib(L, CMathlib);
	return 1;
}
