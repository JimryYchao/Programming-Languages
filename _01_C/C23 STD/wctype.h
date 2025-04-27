#pragma once
typedef int _defined_;

typedef _defined_ wint_t;		// 保存与扩展字符集成员相对应的任何值
#define WEOF ((wint_t)(0xFFFF)) // 宽字符流结束标志

int iswalnum(wint_t wc);				// 字母或数字
int iswalpha(wint_t wc);				// 字母
int iswblank(wint_t wc);				// 空格或制表符
int iswspace(wint_t wc);				// 空白字符
int iswdigit(wint_t wc);				// 数字字符
int iswxdigit(wint_t wc);				// 十六进制数字字符
int iswgraph(wint_t wc);				// 图形字符
int iswprint(wint_t wc);				// 可打印字符（包括空格）
int iswpunct(wint_t wc);				// 标点符号
int iswcntrl(wint_t wc);				// 控制字符
int iswlower(wint_t wc);				// 小写字母
int iswupper(wint_t wc);				// 大写字母
typedef _defined_ wctype_t;				// 保存表示本地环境的字符分类的值
wctype_t wctype(const char *property);	// 查找当前 C 本地环境中的字符分类类别
int iswctype(wint_t wc, wctype_t desc); // 按照指定 LC_CTYPE 类别分类宽字符

wint_t towlower(wint_t wc);					 // 转换为小写字母
wint_t towupper(wint_t wc);					 // 转换为大写字母
typedef _defined_ wctrans_t;				 // 保存表示特定于区域设置的字符映射的值
wctrans_t wctrans(const char *property);	 // 查找当前 C 本地环境中的字符映射类别
wint_t towctrans(wint_t wc, wctrans_t desc); // 按照指定的 LC_TYPE 映射分类进行字符映射

// Example: 宽字符分类
void Wctype()
{
	setlocale(LC_ALL, "zh_CN.UTF-8");	// 设置区域为中文
	puts("The character \u6c34 is..."); // 水
	const char *cats[] = {"digit", "alpha", "space", "cntrl"};
	for (int n = 0; n < 4; ++n)
		printf("%s?\t%s\n", cats[n], iswctype(L'\u6c34', wctype(cats[n])) ? "true" : "false");
	/*
		digit?  false
		alpha?  true
		space?  false
		cntrl?  false
	*/
}
