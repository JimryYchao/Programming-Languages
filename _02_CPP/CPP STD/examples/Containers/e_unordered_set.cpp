#include <iostream>
#include <unordered_set>
#include <string>
#include <vector>
#include <algorithm>

module Containers;
using namespace std;


void print(const auto& s, const std::string& name) {
    std::cout << name << ": ";
    for (const auto& num : s) {
        std::cout << num << " ";
    }
    std::cout << std::endl;
};

// std::unordered_set - 无序且唯一的元素集合
void example_unordered_set() {
    std::cout << "\n=== std::unordered_set ===\n";

    // 构造方式
    std::unordered_set<int> us1;
    std::unordered_set<int> us2 = {3, 1, 4, 1, 5, 9};  // 自动去重
    std::unordered_set<int> us3(us2.begin(), us2.end());
    std::unordered_set<int> us4(us2);

    // 插入元素
    us1.insert(10);
    us1.insert({20, 30, 40, 50});
    
    // 插入重复元素无效
    auto result = us1.insert(30);
    if (!result.second) {
        std::cout << "无法插入重复元素 30（已存在）" << std::endl;
    }

    print(us1, "us1");
    print(us2, "us2");

    // 容器属性
    std::cout << "us1 size: " << us1.size() << ", empty: " << std::boolalpha << us1.empty() << std::endl;
    std::cout << "us1 bucket_count: " << us1.bucket_count() << std::endl;

    // 查找元素
    auto it = us1.find(30);
    if (it != us1.end()) {
        std::cout << "在 us1 中找到元素 30" << std::endl;
    }

    // 统计元素数量
    std::cout << "us1 中 40 的数量: " << us1.count(40) << std::endl;

    // 桶操作（unordered_set特有）
    std::cout << "元素 20 所在的桶索引: " << us1.bucket(20) << std::endl;

    // 删除元素
    us1.erase(30);  // 按值删除
    print(us1, "us1 删除 30 后");

    it = us1.find(40);
    if (it != us1.end()) {
        us1.erase(it);  // 按迭代器删除
    }
    print(us1, "us1 删除 40 后");

    // 清空容器
    us3.clear();
    std::cout << "us3 清空后大小: " << us3.size() << std::endl;

    // 交换两个unordered_set
    std::unordered_set<int> setA = {1, 2, 3};
    std::unordered_set<int> setB = {4, 5, 6};
    
    print(setA, "交换前 setA");
    print(setB, "交换前 setB");
    setA.swap(setB);
    print(setA, "交换后 setA");
    print(setB, "交换后 setB");

    // 赋值操作
    us4 = us1;
    print(us4, "us4 赋值后");

    // 重哈希和预留空间
    us1.rehash(100);
    std::cout << "重哈希后桶数量: " << us1.bucket_count() << std::endl;
}

// std::unordered_multiset - 允许重复元素的无序集合
void example_unordered_multiset() {
    std::cout << "\n=== std::unordered_multiset ===\n";

    std::unordered_multiset<int> ums = {3, 1, 4, 1, 5, 9, 5};
    print(ums, "ums");

    std::cout << "ums size: " << ums.size() << ", empty: " << std::boolalpha << ums.empty() << std::endl;
    std::cout << "ums 中 1 的数量: " << ums.count(1) << std::endl;

    // 删除所有值为1的元素
    ums.erase(1);
    std::cout << "删除所有 1 后，1 的数量: " << ums.count(1) << ", ums 大小: " << ums.size() << std::endl;
}

// 高级特性 - 自定义类型和哈希函数
void example_unordered_set_advanced() {
    std::cout << "\n=== std::unordered_set 高级特性 ===\n";

    // 自定义类型
    struct Person {
        std::string name;
        int age;
        
        Person(const std::string& n, int a) : name(n), age(a) {}
        
        // 相等比较运算符
        bool operator==(const Person& other) const {
            return name == other.name && age == other.age;
        }
    };

    // 自定义哈希函数
    struct PersonHash {
        std::size_t operator()(const Person& p) const {
            std::size_t h1 = std::hash<std::string>{}(p.name);
            std::size_t h2 = std::hash<int>{}(p.age);
            return h1 ^ (h2 << 1); // 简单组合哈希值
        }
    };

    // 使用自定义类型和哈希函数的unordered_set
    std::unordered_set<Person, PersonHash> people;
    people.insert(Person("Alice", 30));
    people.insert(Person("Bob", 25));
    people.insert(Person("Charlie", 35));
    people.insert(Person("Alice", 28)); // 不同年龄，视为不同键

    std::cout << "人员集合: " << std::endl;
    for (const auto& p : people) {
        std::cout << "姓名: " << p.name << ", 年龄: " << p.age << std::endl;
    }

    // 大小写不敏感的字符串集合
    struct CaseInsensitiveHash {
        std::size_t operator()(const std::string& s) const {
            std::string lower;
            lower.reserve(s.size());
            for (char c : s) {
                lower.push_back(std::tolower(c));
            }
            return std::hash<std::string>{}(lower);
        }
    };

    struct CaseInsensitiveEqual {
        bool operator()(const std::string& a, const std::string& b) const {
            if (a.length() != b.length()) return false;
            for (size_t i = 0; i < a.length(); ++i) {
                if (std::tolower(a[i]) != std::tolower(b[i])) {
                    return false;
                }
            }
            return true;
        }
    };

    std::unordered_set<std::string, CaseInsensitiveHash, CaseInsensitiveEqual> caseInsensitiveSet;
    caseInsensitiveSet.insert("Hello");
    caseInsensitiveSet.insert("WORLD");
    caseInsensitiveSet.insert("hello"); // 不会被插入，因为被视为重复

    print(caseInsensitiveSet, "大小写不敏感集合");
    std::cout << "'world' 在集合中的数量: " << caseInsensitiveSet.count("world") << std::endl;
}

void test_unordered_set() {
    example_unordered_set();
    example_unordered_multiset();
    example_unordered_set_advanced();
}