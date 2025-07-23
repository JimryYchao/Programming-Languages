#pragma comment(lib, "LuaDll.lib") 

#include "lua.hpp"
#include <string.h>
#include <threads.h>

typedef struct Proc {
	lua_State* L;
	thrd_t thread;
	cnd_t cond;
	const char* channel;
	struct Proc* previous, * next;
} Proc;

static Proc* waitsend = NULL;
static Proc* waitreceive = NULL;

static mtx_t kernel_access;

static Proc* getself(lua_State* L);
static void movevalues(lua_State* send, lua_State* rec);
static Proc* searchmatch(const char* channel, Proc** list);
static void waitonlist(lua_State* L, const char* channel, Proc** list);
static int ll_send(lua_State* L);
static int ll_receive(lua_State* L);
static int ll_thread(void* arg);
static int ll_start(lua_State* L);
static int ll_exit(lua_State* L);
static void registerlib(lua_State* L, const char* name, lua_CFunction f);
static void openlibs(lua_State* L);
int luaopen_lproc(lua_State* L);


static Proc* getself(lua_State* L) {
	Proc* p;
	lua_getfield(L, LUA_REGISTRYINDEX, "__self");
	p = (Proc*)lua_touserdata(L, -1);
	lua_pop(L, 1);
	return p;
}
// 将发送进程的栈中所有的值（除了第一个，它是通道）移动到接收进程的栈中
static void movevalues(lua_State* send, lua_State* rec) {
	int n = lua_gettop(send);
	luaL_checkstack(rec, n, "too many resutls");
	for (int i = 2; i <= n; i++)
		lua_pushstring(rec, lua_tostring(send, i));
}
// 寻找等待通道的进程的函数
static Proc* searchmatch(const char* channel, Proc** list) {
	Proc* node;
	for (node = *list; node != NULL; node = node->next)
	{
		if (strcmp(channel, node->channel) == 0) {
			if (*list == node)
				*list = (node->next == node) ? NULL : node->next;
			node->previous->next = node->next;
			node->next->previous = node->previous;
			return node;
		}
	}
	return NULL;  // 未匹配
}

// 用于在等待列表中新增一个进程的函数
static void waitonlist(lua_State* L, const char* channel, Proc** list) {
	Proc* p = getself(L);
	if (*list == NULL) {
		*list = p;
		p->previous = p->next = p;
	}
	else {
		p->previous = (*list)->previous;
		p->next = *list;
		p->previous->next = p->next->previous = p;
	}

	p->channel = channel;
	do {
		cnd_wait(&p->cond, &kernel_access);
	} while (p->channel);
}

// 用于发送消息的函数
static int ll_send(lua_State* L) {
	Proc* p;
	const char* channel = luaL_checkstring(L, 1);
	mtx_lock(&kernel_access);
	p = searchmatch(channel, &waitreceive);
	if (p) {
		movevalues(L, p->L);
		p->channel = NULL;
		cnd_signal(&p->cond);
	}
	else
		waitonlist(L, channel, &waitsend);
	mtx_unlock(&kernel_access);
	return 0;
}

// 用于接收消息的函数
static int ll_receive(lua_State* L) {
	Proc* p;
	const char* channel = luaL_checkstring(L, 1);
	lua_settop(L, 1);
	mtx_lock(&kernel_access);
	p = searchmatch(channel, &waitsend);
	if (p) {
		movevalues(p->L, L);
		p->channel = NULL;
		cnd_signal(&p->cond);
	}
	else
		waitonlist(L, channel, &waitreceive);

	mtx_unlock(&kernel_access);
	return lua_gettop(L) - 1;
}

// 用于创建进程的函数
static int ll_start(lua_State* L) {
	thrd_t thread;
	const char* chunk = luaL_checkstring(L, 1);
	lua_State* L1 = luaL_newstate();
	if (L1 == NULL)
		luaL_error(L, "unable to create new state");
	if (luaL_loadstring(L1, chunk) != 0)
		luaL_error(L, "error in thread body: %s", lua_tostring(L1, -1));
	if (thrd_create(&thread, ll_thread, L1) != 0)
		luaL_error(L, "unable to create new thread");
	thrd_detach(thread);
	return 0;
}

static void registerlib(lua_State* L, const char* name, lua_CFunction f) {
	lua_getglobal(L, "package");
	lua_getfield(L, -1, "preload");
	lua_pushcfunction(L, f);
	lua_setfield(L, -2, name);
	lua_pop(L, 2);
}

static void openlibs(lua_State* L) {
	luaL_requiref(L, "_G", luaopen_base, 1);
	luaL_requiref(L, "package", luaopen_package, 1);
	lua_pop(L, 2);
	registerlib(L, "coroutine", luaopen_coroutine);
	registerlib(L, "table", luaopen_table);
	registerlib(L, "io", luaopen_io);
	registerlib(L, "os", luaopen_os);
	registerlib(L, "string ", luaopen_string);
	registerlib(L, "math", luaopen_math);
	registerlib(L, "utf8", luaopen_utf8);
	registerlib(L, "debug", luaopen_debug);
}

static int ll_thread(void* arg) {
	lua_State* L = (lua_State*)arg;
	Proc* self;
	openlibs(L);

	luaL_requiref(L, "lproc", luaopen_lproc, 1);
	lua_pop(L, 1);
	self = (Proc*)lua_newuserdata(L, sizeof(Proc));
	lua_setfield(L, LUA_REGISTRYINDEX, "__self");
	self->L = L;
	self->thread = thrd_current();
	self->channel = NULL;
	cnd_init(&self->cond);

	if (lua_pcall(L, 0, 0, 0) != 0)
		fprintf(stderr, "thread error: %s", lua_tostring(L, -1));

	cnd_destroy(&getself(L)->cond);
	lua_close(L);
	return 0;
}

static int ll_exit(lua_State* L) {
	thrd_exit(0);
	return 0;
}

static int ll_threadID(lua_State* L) {
	lua_pushnumber(L, thrd_current()._Tid);
	return 1;
}

static const struct luaL_Reg ll_funcs[] = {
	{"start", ll_start},
	{"send", ll_send},
	{"receive", ll_receive},
	{"exit", ll_exit},
	{"threadID", ll_threadID},
	{NULL,NULL},
};

int luaopen_lproc(lua_State* L) {
	luaL_newlib(L, ll_funcs);
	return 1;
}

int main() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	luaL_requiref(L, "lproc", luaopen_lproc, 1);
	lua_pop(L, 1);
	luaL_dofile(L, "lproc.lua");
	struct timespec tm = { .tv_sec = 1 };
	lua_close(L);
	return thrd_sleep(&tm, NULL);
}