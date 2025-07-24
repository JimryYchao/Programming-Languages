#include "test.h"

const char* Op_Add = "Add";
const char* Op_Mul = "Mul";
const char* Op_Sub = "Sub";
const char* Op_Div = "Div";

static double lua_Op(lua_State* L, double x, double y, const char* op) {
	lua_getglobal(L, op);
	lua_pushnumber(L, x);
	lua_pushnumber(L, y);

	if (lua_pcall(L, 2, 1, 0) != LUA_OK)
		luaL_error(L, "error running function '%s':%s", op, lua_tostring(L, -1));
	int isnum;
	double rt = lua_tonumberx(L, -1, &isnum);
	if (!isnum)
		luaL_error(L, "function '%s' should return a number", op);

	lua_pop(L, 1);
	return rt;
}

static double lua_Opi(lua_State* L, int x, int y, const char* op) {
	lua_getglobal(L, op);
	lua_pushinteger(L, x);
	lua_pushinteger(L, y);

	if (lua_pcall(L, 2, 1, 0) != LUA_OK)
		luaL_error(L, "error running function '%s':%s", op, lua_tostring(L, -1));
	int isnum;
	int rt = (int)lua_tonumberx(L, -1, &isnum);
	if (!isnum)
		luaL_error(L, "function '%s' should return a integer", op);

	lua_pop(L, 1);
	return rt;
}

static void lua2C_Op(lua_State* L, double x, double y, const char* op) {
	double rt = lua_Op(L, x, y, op);
	printf("in Lua, %g %s %g = %g\n", x, op, y, rt);
}

static void lua2C_Opi(lua_State* L, int x, int y, const char* op) {
	int rt = lua_Opi(L, x, y, op);
	printf("in Lua, %i %s %i = %i\n", x, op, y, rt);
}

void test_C_Call_Lua() {
	lua_State* L = luaL_newstate();
	luaL_dofile(L, "CCallLua.lua");

	// test call lua function
	lua2C_Op(L, 3.14, 9.81, Op_Mul);
	lua2C_Op(L, 3.14, 9.81, Op_Sub);

	lua2C_Opi(L, 99, 5, Op_Add);
	lua2C_Opi(L, 99, 5, Op_Div);

	lua_close(L);
}