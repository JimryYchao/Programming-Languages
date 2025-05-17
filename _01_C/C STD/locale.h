#pragma once

struct lconv // localeconv() return
{
    // 非货币数值格式化参数
    char *decimal_point; // "." // 用于格式化非货币量的小数点字符
    char *thousands_sep; // ""  // 在格式化的非货币量中，用于分隔小数点字符前的数字组的字符
    char *grouping;      // ""  // 一个字符串，其元素指示格式化的非货币量中每组数字的大小
    // 货币数值格式化参数
    char *mon_decimal_point; // ""  // 用于格式化货币数量的小数点
    char *mon_thousands_sep; // ""  // 格式化货币量中小数点前的数字组的分隔符
    char *mon_grouping;      // ""  // 一个字符串，其元素指示格式化货币量中每组数字的大小
    char *positive_sign;     // ""  // 用于指示非负值格式货币量的字符串
    char *negative_sign;     // ""  // 用于指示负值格式货币量的字符串
    // 本地货币数值格式化参数
    char *currency_symbol; // ""  		 // 适用于当前区域设置的本地货币符号
    char frac_digits;      // CHAR_MAX  // 要在本地格式化货币量中显示的小数位数（小数点后的数字）
    char p_cs_precedes;    // CHAR_MAX  // 若 currency_symbol 置于非负本地格式货币量的值之前则为 1，于其后则为 ​0​
    char n_cs_precedes;    // CHAR_MAX  // 若 currency_symbol 置于负本地格式货币量的值之前则为 1，于其后则为 ​0​
    char p_sep_by_space;   // CHAR_MAX  // 设置为一个值，指示 currency_symbol、positive_sign 及非负货币量的分隔
    char n_sep_by_space;   // CHAR_MAX  // 设置为一个值，指示 currency_symbol、positive_sign 及负货币量的分隔
    char p_sign_posn;      // CHAR_MAX  // 设置为一个值，指示非负货币量中 positive_sign 的位置
    char n_sign_posn;      // CHAR_MAX  // 设置为一个值，指示负货币量中 negative_sign 的位置
    // 国际货币数值格式化参数
    char *int_curr_symbol;   // ""  		 // 适用于当前区域设置的国际货币符号
    char int_frac_digits;    // CHAR_MAX  // 在国际格式的货币量中显示的小数位数（小数点后的数字）
    char int_p_cs_precedes;  // CHAR_MAX	 // 如果 int_curr_symbol 置于非负国际格式货币量值的前面或后面，则设置为 1 或 0
    char int_n_cs_precedes;  // CHAR_MAX  // 如果 int_curr_symbol 置于负国际格式货币量值的前面或后面，则设置为 1 或 0
    char int_p_sep_by_space; // CHAR_MAX  // 设置为一个值，该值指示 int_curr_symbol、符号字符串和非负国际格式货币数量的值的分隔
    char int_n_sep_by_space; // CHAR_MAX  // 设置为一个值，该值指示 int_curr_symbol、符号字符串和负国际格式货币数量的值的分隔
    char int_p_sign_posn;    // CHAR_MAX  // 设置为一个值，该值指示非负国际格式货币量的 positive_sign 的位置
    char int_n_sign_posn;    // CHAR_MAX  // 设置为一个值，该值指示负数国际格式货币量的 negative_sign 的位置
}

#define NULL        // 空指针常量
#define LC_COLLATE  // 字符串排序类别
#define LC_CTYPE    // 字符分类类别
#define LC_MONETARY // 货币格式类别
#define LC_NUMERIC  // 数字格式类别
#define LC_TIME     // 时间格式类别
#define LC_ALL (LC_COLLATE | LC_CTYPE | LC_MONETARY | LC_NUMERIC | LC_TIME)

char * setlocale(int category, const char *locale); // 设置或获取当前区域设置。locale = "" 应用本地区域设置；locale = "C" 应用最小环境设置
struct lconv *localeconv(void);              // 获取区域设置的详细信息

// Example: Setlocale
void LocaleTest(void)
{
    char *locale = setlocale(LC_ALL, NULL); // 获取当前区域设置
    setlocale(LC_ALL, "zh_CN.UTF-8");       // 设置区域设置为中文（中国）
    struct lconv *lc = localeconv();
    printf("Currency symbol: %s\n", lc->int_curr_symbol); // CNY
    setlocale(LC_MONETARY, "en_IN.utf8");
    lc = localeconv();                                 // 更新 locale 需要同时更新 lconv 对象
    printf("本地货币符号: %s\n", lc->currency_symbol); // ₹
    printf("国际货币符号: %s\n", lc->int_curr_symbol); // INR
}
