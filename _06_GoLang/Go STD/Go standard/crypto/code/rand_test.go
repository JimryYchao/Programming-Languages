package gostd

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"io"
	"strings"
	"testing"
)

/*
! rand.Reader 是加密安全随机数生成器的全局共享实例。
! Int 返回 [0，max) 中的随机值。
! Prime 返回一个给定位长的数字，该数字是素数的概率很高。如果 rand.Read 返回任何错误或 bits < 2，则 Prime 将返回错误。
! Read 是一个 helper 函数，它使用 io.ReadFull 调用 Reader.Read, 生成一个随机 []byte。返回时当且仅当 err == nil, n == len(B)。
*/

func TestRand(t *testing.T) {
	key := make([]byte, md5.BlockSize)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatal(err)
	}
	h := hmac.New(md5.New, key)
	io.Copy(h, strings.NewReader("Hello World"))
	logfln("%x", h.Sum(nil))
}
