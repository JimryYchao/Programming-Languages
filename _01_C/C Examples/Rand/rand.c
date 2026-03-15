// 随机数算法, 受限于种子
static unsigned long int seed = 1;
unsigned int rand(void)  // 0~32767
{
    seed = seed * 1103515245 + 12345;
    return (unsigned int) (seed / 65536) % 32768;
}
void srand(unsigned int _seed)   
{
    seed = _seed;
    // 可以借助时钟系统传递种子
    // 例如: srand((unsigned int)time(NULL)); or time(0)
}

// 掷骰子, 一个 n 面的骰子，返回 1~n 之间的随机数
#include <time.h>
#include <stdio.h>
int roll(int sides){  
    if(sides < 2){
        puts("Need at least 2 sides");
        return -1;
    }
    return rand() % sides + 1;
}
int rolln(int n, int sides){  // 掷 n 个骰子, 每个骰子有 sides 个面
    if(n < 1){
        puts("Need at least 1 die and");
        return -1;
    }
    int total = 0;
    while(n--)
        total += roll(sides);
    return total;
}

int main(){
    srand((unsigned int)time(0));
    printf("get a rand number: %d\n", rand());

    for (size_t i = 0; i < 5; i++)
        printf("roll a 6-sides die: %d\n", roll(6));
    
    for (size_t i = 0; i < 30; i++)
        printf("roll 2 3-sides dice: %d\n", rolln(2, 3));
}





