// Lua 5.4.8
#pragma once
#include "lua_ref.h"

int luaopen_base(lua_State *L);
int luaopen_coroutine(lua_State *L);
int luaopen_table(lua_State *L);
int luaopen_io(lua_State *L);
int luaopen_os(lua_State *L);
int luaopen_string(lua_State *L);
int luaopen_utf8(lua_State *L);
int luaopen_math(lua_State *L);
int luaopen_debug(lua_State *L);
int luaopen_package(lua_State *L);