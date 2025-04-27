package examples

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

// Scan File
func TestScannerFile(t *testing.T) {
	file, err := os.Open("readline.file")
	if err != nil {
		t.Fatal(err)
	}
	scr := bufio.NewScanner(file)
	scr.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		var c = 0
		for len(data) > 0 { // 跳过空行
			if data[0] == '\r' || data[0] == '\n' {
				c++
				data = data[1:]
			} else {
				advance, token, err = bufio.ScanLines(data, atEOF)
				return advance + c, token, err
			}
		}
		return 0, nil, io.EOF
	})

	for scr.Scan() {
		fmt.Printf("scan line: %s\n", scr.Text())
	}
	if err = scr.Err(); err != nil {
		t.Fatal(err)
	}
}
