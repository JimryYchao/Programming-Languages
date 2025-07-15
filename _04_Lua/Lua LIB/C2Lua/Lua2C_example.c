#include <stdio.h>
#include "../src/lua.h"
#include "../src/lauxlib.h"

int getglobint(lua_State *L, const char *var) {
	int isnum, result;
	lua_getglobal(L, var);
	result = (int)lua_tointegerx(L, -1, &isnum);
	if (!isnum) 
		luaL_error(L, "'%s' should be a number\n", var);
	lua_pop(L, 1);
	return result;
}

void load(lua_State* L, const char* fname, int* w, int* h) {
	if (luaL_loadfile(L, fname) || lua_pcall(L, 0, 0, 0))
		luaL_error(L, "cannot run config. file: %s\n", lua_tostring(L, -1));
	*w = getglobint(L, "width");
	*h = getglobint(L, "height");
}


int Test_Config() {
	const char* fname = "screen.lua";
	int w, h;

	lua_State* L = luaL_newstate();
	load(L, fname, &w, &h);
	printf("w = %d; h = %d;", w, h);
	lua_close(L);
	return 0;
}