#pragma once
#pragma comment(lib, "lua550.lib")   

#include "lua.h"
#include "lauxlib.h"
#include "lualib.h"

void test_BitArray();
void test_C_Call_Lua();
void test_Lua_Call_C();
void test_lproc();
void test_Continue_pcall();
void test_Continue_CFunc();
void test_memlimit();