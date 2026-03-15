// C 结构体（纯数据对象 POD）二进制序列化 / 反序列化 API
// 适用前提：结构体内只含固定大小的数据成员（无指针、无嵌套指针）

// 安全函数说明：C11 Annex K 的 fwrite_s/fread_s 为可选扩展。
// 本机 gcc(mingw+msvcrt)/clang/MSVC 运行库均未提供 fwrite_s（MSVC 亦无），
// fread_s 在 mingw 仅声明不可链接。故采用条件编译：
//   运行库完整实现 Annex K（__STDC_LIB_EXT1__）时使用 fread_s/fwrite_s，
//   否则回退到对返回值严格校验的 fread/fwrite（安全语义等价）。
#define __STDC_WANT_LIB_EXT1__ 1
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#if defined(__STDC_LIB_EXT1__) && (__STDC_LIB_EXT1__ + 0)
#define SER_USE_EXT1 1
#endif

// ============================================================================
// 序列化结果
// ============================================================================

typedef enum {
    SER_SUCCESS = 0,
    SER_ERR_NULL,
    SER_ERR_SIZE,
    SER_ERR_OPEN,
    SER_ERR_WRITE,
    SER_ERR_READ
} SerResult;

static const char *ser_msg(SerResult r) {
    static const char *t[] = {
        "success", "null pointer", "invalid size", "open failed", "write failed", "read failed"
    };
    return (r >= 0 && r <= SER_ERR_READ) ? t[r] : "unknown";
}

// ============================================================================
// 核心写入：将结构体写入已打开的文件流
//   obj  : 指向结构体对象的指针
//   size : 结构体大小（sizeof(T)）
//   fp   : 已以二进制写模式打开的 FILE*
// 数据布局：size | data(size)
// ============================================================================

SerResult ser_write(const void *obj, size_t size, FILE *fp) {
    if (!obj || !fp) return SER_ERR_NULL;
    if (size == 0)   return SER_ERR_SIZE;

    if (fwrite(&size, sizeof size, 1, fp) != 1) return SER_ERR_WRITE;
#ifdef SER_USE_EXT1
    // Annex K: errno_t fwrite_s(buffer, size, count, stream)，返回 0 表示成功
    if (fwrite_s(obj, size, 1, fp) != 0) return SER_ERR_WRITE;
#else
    if (fwrite(obj,  size,        1, fp) != 1) return SER_ERR_WRITE;
#endif
    return SER_SUCCESS;
}

// ============================================================================
// 核心读取：从已打开的文件流恢复结构体
//   obj  : 指向已分配的结构体对象的指针
//   size : 预期结构体大小（sizeof(T)）
//   fp   : 已以二进制读模式打开的 FILE*
// ============================================================================

SerResult ser_read(void *obj, size_t size, FILE *fp) {
    if (!obj || !fp) return SER_ERR_NULL;
    if (size == 0)   return SER_ERR_SIZE;

    size_t saved;
#ifdef SER_USE_EXT1
    // Annex K: size_t fread_s(buffer, bufsz, size, count, stream)，返回完整读入的元素数
    if (fread_s(&saved, sizeof saved, sizeof saved, 1, fp) != 1) return SER_ERR_READ;
#else
    if (fread(&saved, sizeof saved, 1, fp) != 1) return SER_ERR_READ;
#endif
    if (saved != size) return SER_ERR_SIZE;
#ifdef SER_USE_EXT1
    if (fread_s(obj, size, size, 1, fp) != 1) return SER_ERR_READ;
#else
    if (fread(obj, size, 1, fp) != 1) return SER_ERR_READ;
#endif
    return SER_SUCCESS;
}

// ============================================================================
// ser_save：打开用于写入的存储文件
//   path : 文件路径
//   fp   : 输出参数，成功时导出已定位到文件头的 FILE*
//   行为："wb" 模式：新建文件或清空原内容，位置定于 0
//   失败：*fp 置为 NULL，返回 SER_ERR_OPEN
// ============================================================================

SerResult ser_save(const char *path, FILE **fp) {
    if (!path || !fp) return SER_ERR_NULL;

    FILE *f = fopen(path, "wb");
    if (!f) { *fp = NULL; return SER_ERR_OPEN; }

    rewind(f);              // 确保定位到文件头
    *fp = f;
    return SER_SUCCESS;
}

// ============================================================================
// ser_load：打开用于读取的存储文件
//   path : 文件路径
//   fp   : 输出参数，成功时导出已定位到文件头的 FILE*
//   失败：*fp 置为 NULL，返回 SER_ERR_OPEN
// ============================================================================

SerResult ser_load(const char *path, FILE **fp) {
    if (!path || !fp) return SER_ERR_NULL;

    FILE *f = fopen(path, "rb");
    if (!f) { *fp = NULL; return SER_ERR_OPEN; }

    rewind(f);
    *fp = f;
    return SER_SUCCESS;
}

// ============================================================================
// ser_close：关闭存储文件
//   fp : 指向 FILE* 的指针；关闭后 *fp 置为 NULL（已为 NULL 时安全）
// ============================================================================

void ser_close(FILE **fp) {
    if (!fp || !*fp) return;
    fclose(*fp);
    *fp = NULL;
}

// ============================================================================
// 示例结构体（纯数据对象）
// ============================================================================

typedef struct {
    int    id;
    char   name[32];
    float  score;
} Student;

// ============================================================================
// 演示：save → write(可多次) → close / load → read(可多次) → close
// ============================================================================

int main(void) {
    FILE *fp = NULL;

    // 1. 保存：打开写入流，连续写入两个对象
    Student s1 = { .id = 1, .name = "Alice", .score = 95.5 };
    Student s2 = { .id = 2, .name = "Bob",   .score = 88.0 };
    if (ser_save("student.dat", &fp) == SER_SUCCESS) {
        ser_write(&s1, sizeof(Student), fp);
        ser_write(&s2, sizeof(Student), fp);
        ser_close(&fp);
        printf("Saved: { %d, %s, %.1f } { %d, %s, %.1f }\n",
               s1.id, s1.name, s1.score, s2.id, s2.name, s2.score);
    }

    // 2. 加载：打开读取流，连续读回两个对象
    Student r1 = { 0 }, r2 = { 0 };
    if (ser_load("student.dat", &fp) == SER_SUCCESS) {
        ser_read(&r1, sizeof(Student), fp);
        ser_read(&r2, sizeof(Student), fp);
        ser_close(&fp);
        printf("Loaded: { %d, %s, %.1f } { %d, %s, %.1f }\n",
               r1.id, r1.name, r1.score, r2.id, r2.name, r2.score);
    }
}
