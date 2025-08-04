#include "test.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>


// 演示格式化 IO
static void example_formatted_io(const char* str)
{
	puts("\n[Formatted I/O]");
	int len = snprintf(NULL, NULL, "Receive: %s", str) + 1;
	char* buf = (char*)calloc(1, len * sizeof(char));
	int cn = snprintf(buf, len, "Receive: %s", str);
	printf("Output: '%s' (%d characters needed)\n", buf, cn);
	free(buf);
}

// 演示文件操作
static void example_file_operations(void)
{
	puts("\n[File Operations]");
	puts("press '#' to exit...");

	// 写入文件
	FILE* fp = fopen("test.txt", "w");
	if (!fp) {
		printf("Error opening file: %s\n", strerror(errno));
		return;
	}

	char buf[128];
	while (1) {
		if (fscanf_s(stdin, "%99[^#]", &buf, 128) == 1)  // 直至读取到 # 中止
			fprintf_s(fp, buf);
		else break;
	}
	fclose(fp);

	// 读取文件
	if (fopen_s(&fp, "test.txt", "r") == 0) {
		puts("Read test.txt :");

		int c;
		while ((c = fgetc(fp))!= EOF)
			putc(c, stdout);

		if (ferror(fp)) {
			puts("\n\nI/O error when reading");
			exit(EXIT_FAILURE);
		}
		else if (feof(fp))
			puts("\n\nEnd of file is reached successfully");
		fclose(fp);
	}
	remove("test.txt");
}

// 演示二进制 I/O
static void example_binary_io(void)
{
	puts("\n[Binary I/O]");
	struct Data {

		int id;
		double value;
		char tag[4];
	} data = { 123, 4.56, "XYZ" }, read_data;

	// 二进制写入
	FILE* fp = fopen("data.bin", "wb");
	if (fp) {
		fwrite(&data, sizeof(data), 1, fp);  // 逐字节写入
		fclose(fp);
	}

	// 二进制读取
	if ((fp = fopen("data.bin", "rb"))) {
		fread(&read_data, sizeof(read_data), 1, fp);  // 逐字节读取
		printf("Read binary: %d, %.2f, %s\n",
			read_data.id, read_data.value, read_data.tag);
		fclose(fp);
	}

	remove("data.bin");
}

// 演示流定位
static void example_stream_positioning(void)
{
	puts("\n[Stream Positioning]");

	FILE* fp = tmpfile();
	if (!fp) return;

	fprintf(fp, "1234567890");
	rewind(fp);

	// fgetpos/fsetpos
	fpos_t pos;
	fgetpos(fp, &pos);
	fseek(fp, 3, SEEK_SET);

	char c = fgetc(fp);
	printf("Character '%c' at position 3\n", c);

	fsetpos(fp, &pos);
	c = fgetc(fp);
	printf("After reset: '%c'\n", c);

	fclose(fp);
}

void test_stdio(void)
{
	example_formatted_io("Hello World");
	example_file_operations();
	example_binary_io();
	example_stream_positioning();
}