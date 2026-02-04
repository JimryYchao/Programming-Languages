#include <stdio.h>
#include "test.h"  

typedef struct
{
	size_t current_bytes;     // 当前已分配内存总量
	size_t limit_bytes;       // 内存限制阈值
	lua_Alloc original_alloc; // 原始分配函数
	void* original_ud;        // 原始用户数据
} MemLimitState;

// 自定义内存分配函数
static void* limited_alloc(void* ud, void* ptr, size_t osize, size_t nsize)
{
	MemLimitState* state = (MemLimitState*)ud;

	// 计算内存变化量（考虑realloc场景）
	size_t delta = (nsize > osize) ? (nsize - osize) : 0;

	// 检查内存限制（仅对新分配或扩容请求）
	if (nsize > 0 && (state->current_bytes + delta > state->limit_bytes))
		return NULL; // 超过限制时返回NULL

	// 调用原始分配器（使用原始用户数据）
	void* new_ptr = state->original_alloc(state->original_ud, ptr, osize, nsize);

	if (new_ptr)
		// 更新内存统计（考虑释放、重新分配等情况）
		if (nsize == 0)
			state->current_bytes -= osize; // 释放内存
		else
			state->current_bytes += (nsize - (ptr ? osize : 0));
	return new_ptr;
}

// 设置内存限制的Lua函数
static int setlimit(lua_State* L)
{
	size_t limit = (size_t)luaL_checkinteger(L, 1);
	// 获取当前分配器状态
	MemLimitState* state;
	lua_getallocf(L, (void**)&state);

	if (!state)
	{
		// 首次调用时初始化状态
		state = (MemLimitState*)lua_newuserdata(L, sizeof(MemLimitState));
		state->current_bytes = 0;
		state->original_alloc = lua_getallocf(L, &state->original_ud);
	}

	// 更新内存限制
	state->limit_bytes = limit;

	// 设置新的分配器（保留原始分配器信息）
	lua_setallocf(L, limited_alloc, state);

	return 0;
}

static const struct luaL_Reg libs[] = {
		{"set", setlimit},
		{NULL, NULL}
};

int luaopen_memlimit(lua_State* L)
{
	luaL_newlib(L, libs);
	return 1;
}

void test_memlimit() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	luaL_requiref(L, "memlimit", luaopen_memlimit, 1);
	lua_pop(L, 1);

	if (luaL_dofile(L, "memlimit.lua"))
		luaL_error(L, "%s", lua_tostring(L, -1));

	lua_close(L);
}

