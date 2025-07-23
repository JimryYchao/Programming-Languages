#pragma comment(lib, "LuaDll.lib") 

#include "lua.h"
#include "lauxlib.h"
#include "lualib.h"

// 用户定义布尔数组
#include <limits.h>
#define BITS_PER_WORD (CHAR_BIT * sizeof(unsigned))
#define I_WORD(i)     ((unsigned)(i)/BITS_PER_WORD)
#define I_BIT(i)	  (1 << ((unsigned)(i) % BITS_PER_WORD))


typedef struct BitArray {
	int size;
	unsigned int values[];		// VLA
} BitArray;

#define checkArray(L) (BitArray*)luaL_checkudata((L), 1, "LuaBook.array")
// BitArray.new(sz)
static int newArray(lua_State* L) {
	int n = (int)luaL_checkinteger(L, 1);
	luaL_argcheck(L, n >= 1, 1, "invalid size");
	size_t nbytes = sizeof(BitArray) + I_WORD(n) * sizeof(unsigned);
	BitArray* a = (BitArray*)lua_newuserdata(L, nbytes);     // 分配用户数据内存
	a->size = n;
	for (size_t i = 0; i < I_WORD(n); i++)
		a->values[i] = 0;

	luaL_getmetatable(L, "LuaBook.array");
	lua_setmetatable(L, -2);
	return 1;
}

static unsigned* getParams(lua_State* L, unsigned* mask) {
	BitArray* a = checkArray(L);
	int index = (int)luaL_checkinteger(L, 2) - 1;
	luaL_argcheck(L, index >= 0 && index < a->size, 2, "index out of range");
	*mask = I_BIT(index);
	return &a->values[I_WORD(index)];
}

// arr[idx] = v
static int setArray(lua_State* L) {
	unsigned mask;
	unsigned* entry = getParams(L, &mask);
	luaL_checkany(L, 3);
	if (lua_toboolean(L, 3))
		*entry |= mask;
	else *entry &= ~mask;
	return 0;
}
// arr[idx]
static int getArray(lua_State* L) {
	unsigned mask;
	unsigned* entry = getParams(L, &mask);
	lua_pushboolean(L, *entry & mask);
	return 1;
}
// #arr
static int getSize(lua_State* L) {
	BitArray* a = checkArray(L);
	
	lua_pushinteger(L, a->size);
	return 1;
}

// arr:__tostring
static int arrToString(lua_State* L) {
	BitArray* a = checkArray(L);
	lua_pushfstring(L, "array(%d)", a->size);
	return 1;
}


static const luaL_Reg arraylib_f[] = {
	{"new", newArray},
	{NULL, NULL}
};
static const luaL_Reg arraylib_m[] = {
	{"__tostring", arrToString},
	{"__newindex", setArray},
	{"__index", getArray},
	{"__len", getSize},
	{NULL, NULL}
};

int luaopen_BitArray(lua_State* L) {
	// 为用户数据创建元表，防止后续库函数调用传递的用户数据类型参数非 BitArray
	luaL_newmetatable(L, "LuaBook.array");
	luaL_setfuncs(L, arraylib_m, 0);    // 为 LuaBook.array 注册元方法

	luaL_newlib(L, arraylib_f);		       // 创建库
	return 1;
}

int main() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	luaopen_BitArray(L);
	lua_setglobal(L, "BitArray");

	if (luaL_dofile(L, "bitArray.lua"))
		printf("%s", lua_tostring(L, -1));
	lua_close(L);
}