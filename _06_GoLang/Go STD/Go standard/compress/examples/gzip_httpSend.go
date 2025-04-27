package examples

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func ExampleHttpSend() {
	// 这是一个编写压缩/读取器的示例。这对于 HTTP 客户端主体非常有用，如下所示。
	const testdata = "the data to be compressed"

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		zr, err := gzip.NewReader(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		// 输出示例的数据。
		if _, err := io.Copy(os.Stdout, zr); err != nil {
			log.Fatal(err)
		}
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// 要压缩的数据，作为 io.Reader
	dataReader := strings.NewReader(testdata)

	// bodyReader 作为 io.Reader 是 HTTP 请求的主体。
	// httpWriter 作为 io.Writer 是 HTTP 请求的主体。
	bodyReader, httpWriter := io.Pipe()

	//确保 bodyReader 始终处于关闭状态，以便下面的 goroutine 将永远退出。
	defer bodyReader.Close()

	// gzipWriter 将数据压缩到 httpWriter。
	gzipWriter := gzip.NewWriter(httpWriter)

	// errch 从写入 goroutine 中收集任何错误。
	errch := make(chan error, 1)

	go func() {
		defer close(errch)
		sentErr := false
		sendErr := func(err error) {
			if !sentErr {
				errch <- err
				sentErr = true
			}
		}
		// 复制数据到 gzipWriter, gzipWriter 会将数据压缩并发送到 bodyReader。
		if _, err := io.Copy(gzipWriter, dataReader); err != nil && err != io.ErrClosedPipe {
			sendErr(err)
		}
		if err := gzipWriter.Close(); err != nil && err != io.ErrClosedPipe {
			sendErr(err)
		}
		if err := httpWriter.Close(); err != nil && err != io.ErrClosedPipe {
			sendErr(err)
		}
	}()

	// 向测试服务器发送 HTTP 请求。
	req, err := http.NewRequest("PUT", ts.URL, bodyReader)
	if err != nil {
		log.Fatal(err)
	}

	// 将 request 传递给 http.Client.Do 将保证关闭 body，在本例中是 bodyReader。
	resp, err := ts.Client().Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// 检查压缩数据是否出错。
	if err := <-errch; err != nil {
		log.Fatal(err)
	}

	// 在这个例子中，不关心响应。
	resp.Body.Close()
}
