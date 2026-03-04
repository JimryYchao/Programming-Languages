#include <stdio.h>
#include <errno.h>
#include <string.h>
#include <math.h>
#include <stdlib.h>

#include "test.h"

void test_errno(void)
{
    errno = 0;
    (void)fopen("non_existent_file.txt", "r");
    printf("ERRNO[%d]: %s\n", errno, strerror(errno));
    // ERRNO[2]: No such file or directory

    errno = 0;
    (void)sqrt(-1.0);
    printf("ERRNO[%d]: %s\n", errno, strerror(errno));
    // ERRNO[33]: Domain error
  
    errno = 0;
    (void)malloc(SIZE_MAX * 2);
    printf("ERRNO[%d]: %s\n", errno, strerror(errno));
    // ERRNO[12]: Not enough space
}