#pragma once
typedef int _defined_;

typedef _defined_ mbstate_t; // 保存在多字节字符序列和宽字符序列之间转换所需的转换状态信息
typedef _defined_ char8_t;   // 8位字符类型
typedef _defined_ char16_t;  // 16位字符类型
typedef _defined_ char32_t;  // 32位字符类型

// 多字节转 UTF-8，UTF-16，UTF-32
size_t mbrtoc8(char8_t *restrict pc8, const char *restrict s, size_t n, mbstate_t *restrict ps);
size_t mbrtoc16(char16_t *restrict pc16, const char *restrict s, size_t n, mbstate_t *restrict ps);
size_t mbrtoc32(char32_t *restrict pc32, const char *restrict s, size_t n, mbstate_t *restrict ps);
// UTF-8，UTF-16，UTF-32 转多字节
size_t c8rtomb(char *restrict s, char8_t c8, mbstate_t *restrict ps);
size_t c16rtomb(char *restrict s, char16_t c16, mbstate_t *restrict ps);
size_t c32rtomb(char *restrict s, char32_t c32, mbstate_t *restrict ps);

// Example: c32rtomb
void _C32rtombTest()
{
    mbstate_t state;
    char32_t in[] = U"zß水🍌"; // 或 "z\u00df\u6c34\U0001F34C"
    size_t in_sz = sizeof in / sizeof *in;
    for (size_t n = 0; n < in_sz; ++n)
        printf("%#x ", in[n]);   // 逐 UTF-32 字符

    char *out = calloc(MB_CUR_MAX * in_sz, sizeof(char));
    char *p = out;
    for (size_t n = 0; n < in_sz; ++n)
    {
        size_t rc = c32rtomb(p, in[n], &state);
        if (rc == -1) // EOF
            break;
        p += rc;
    }
    size_t out_sz = p - out;
    for (size_t x = 0; x < out_sz; ++x)
        printf("%#x ", +(unsigned char)out[x]);   // 逐多字节字符
    free(out);
    // UTF32：0x7a 0xdf 0x6c34 0x1f34c 0
    // MB：   0x7a 0xc3 0x9f 0xe6 0xb0 0xb4 0xf0 0x9f 0x8d 0x8c 0 
}