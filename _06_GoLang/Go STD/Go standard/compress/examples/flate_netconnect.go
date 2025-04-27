package examples

import (
	"compress/flate"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

// DEFLATE 适用于通过网络传输压缩数据。
func ExampleNetConnect() {
	var wg sync.WaitGroup
	defer wg.Wait()

	// 使用 io.Pipe 模拟一个 network 连接。真实的网络应用程序应该注意正确关闭底层连接。
	rp, wp := io.Pipe()

	// 启动一个程序作为发送端。
	wg.Add(1)
	go func() {
		defer wg.Done()
		zw, err := flate.NewWriter(wp, flate.BestSpeed)
		if err != nil {
			log.Fatal(err)
		}

		b := make([]byte, 256)
		for _, m := range strings.Fields("A long time ago in a galaxy far, far away...") {
			// 使用一个简单的帧格式，其中第一个字节是消息长度，后面跟着消息本身。
			b[0] = uint8(copy(b[1:], m))
			if _, err := zw.Write(b[:1+len(m)]); err != nil {
				log.Fatal(err)
			}
			// Flush 确保接收端可以读取到目前为止发送的所有数据。
			if err := zw.Flush(); err != nil {
				log.Fatal(err)
			}
		}
		if err := zw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// 启动一个程序来充当接收端。
	wg.Add(1)
	go func() {
		defer wg.Done()
		zr := flate.NewReader(rp)
		b := make([]byte, 256)
		for {
			// 读取消息长度。这保证在发送端每次相应的 Flush 和 Close 都返回。
			if _, err := io.ReadFull(zr, b[:1]); err != nil {
				if err == io.EOF {
					break // 发送端关闭了信号流
				}
				log.Fatal(err)
			}

			// 读取消息内容
			n := int(b[0])
			if _, err := io.ReadFull(zr, b[:n]); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Received %d bytes: %s\n", n, b[:n])
		}
		if err := zr.Close(); err != nil {
			log.Fatal(err)
		}
	}()
}
