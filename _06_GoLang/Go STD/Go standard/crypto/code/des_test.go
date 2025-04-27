package gostd

import (
	"crypto/des"
	"testing"
)

/*
! NewCipher 创建并返回一个 DES 的 cipher.Block
! NewTripleDESCipher 创建并返回一个 TDEA 的 cipher.Block。
*/

func TestNewTripleDESCipher(t *testing.T) {
	// 当需要 EDE2 时，也可以使用 NewTripleDESCipher 来复制 16 字节密钥的前 8 个字节。
	ede2Key := []byte("example key 1234")

	var tripleDESKey []byte
	tripleDESKey = append(tripleDESKey, ede2Key[:16]...)
	tripleDESKey = append(tripleDESKey, ede2Key[:8]...)

	block, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		t.Fatal(err)
	}

	// See crypto/cipher or how to use a cipher.Block
	_ = block
}
