#pragma once

// 字符分类函数
int isalnum(int c);        // 字母或数字
int isalpha(int c);        // 字母
int isblank(int c);        // 空格字符
int iscntrl(int c);        // 控制字符
int isdigit(int c);        // 数字字符
int isgraph(int c);        // 可打印字符（除空格外）
int islower(int c);        // 小写字母
int isprint(int c);        // 可打印字符（包括空格）
int ispunct(int c);        // 标点符号
int isspace(int c);        // 空格字符（空格、换行、回车、制表符等）
int isupper(int c);        // 大写字母
int isxdigit(int c);       // 十六进制数字字符（0-9、a-f、A-F）
// 字符映射函数
int tolower(int c);        // 转换为小写字母
int toupper(int c);        // 转换为大写字母
