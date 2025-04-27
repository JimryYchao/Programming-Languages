#pragma once
typedef int _defined_;

typedef _defined_ once_flag;                         // 标识一次性初始化的对象类型
#define ONCE_FLAG_INIT (once_flag){0}                // 一次性初始化对象的初始值
void call_once(once_flag *flag, void (*func)(void)); // 执行一次性初始化函数 func，flag 为一次性初始化对象的指针

// 线程
typedef _defined_ thrd_t;                                                    // 标识线程的完整对象类型
thrd_t thrd_current(void);                                                   // 获取当前线程的ID
typedef int (*thrd_start_t)(void *);                                         // 线程函数类型
int thrd_create(thrd_t *thr, thrd_start_t func, void *arg);                  // 创建线程
int thrd_equal(thrd_t thr0, thrd_t thr1);                                    // 比较两个线程ID是否相等
int thrd_sleep(const struct timespec *duration, struct timespec *remaining); // 休眠指定时间，提前恢复将剩余时间保存至 remaining
int thrd_detach(thrd_t thr);                                                 // 分离线程
int thrd_join(thrd_t thr, int *res);                                         // 等待线程结束，获取返回值
void thrd_yield(void);                                                       // 让出当前线程的执行权
void thrd_exit(int res) [[noreturn]];                                        // 终止线程
enum                                                                         // 线程返回状态
{
    thrd_success,  // 线程成功
    thrd_nomen,    // 内存耗尽
    thrd_timedout, // 超时
    thrd_busy,     // 资源不可用
    thrd_error,    // 其他错误
}
// 互斥
typedef _defined_ mtx_t; // 标识互斥体
enum                     // 互斥体类别
{
    mtx_plain,     // 常规互斥体
    mtx_recursive, // 递归互斥体
    mtx_timed,     // 定时互斥体
    // mtx_plain | mtx_recursive 常规的递归互斥锁；
    // mtx_timed | mtx_recursive 支持超时的递归互斥锁。
};
int mtx_init(mtx_t *mtx, int type);    // 创建，type 为互斥体类型
int mtx_lock(mtx_t *mtx);              // 阻塞当前线程到锁定互斥体
int mtx_trylock(mtx_t *mtx);           // 锁定，已锁定则返回 thrd_busy
int mtx_unlock(mtx_t *mtx);            // 解锁互斥体
void mtx_destroy(mtx_t *mtx);          // 销毁互斥锁
int mtx_timedlock(mtx_t *restrict mtx, // 阻塞当前线程到锁定互斥体，超时则返回 thrd_timedout
                  const struct timespec *restrict ts);
// 条件变量
typedef _defined_ cnd_t;                // 标识条件变量
int cnd_init(cnd_t *cond);              // 创建条件变量
int cnd_signal(cnd_t *cond);            // 取消阻塞 cond 上的一个线程
int cnd_broadcast(cnd_t *cond);         // 取消阻塞 cond 上的所有线程
void cnd_destroy(cnd_t *cond);          // 销毁条件变量
int cnd_wait(cnd_t *cond, mtx_t *mtx);  // 在 cond 上阻塞，mtx 必须在调用前加锁
int cnd_timedwait(cnd_t *restrict cond, // 在 cond 阻塞一段时长
                  mtx_t *restrict mtx, const struct timespec *restrict ts);
// 线程局部存储
typedef void (*tss_dotr_t)(void *);          // 用作 TSS 析构器的函数类型
typedef _defined_ tss_t;                     // 线程特定存储的指针
#define TSS_DTOR_ITERATIONS                  // 析构器调用最大次数
int tss_create(tss_t *key, tss_dtor_t dtor); // 创建线程特定存储，dtor 为析构器函数
void *tss_get(tss_t key);                    // 获取线程特定存储的值
int tss_set(tss_t key, void *val);           // 设置线程特定存储的值
void tss_delete(tss_t key);                  // 删除线程特定存储