package gostd

import (
	"crypto/rc4"
	"testing"
)

/*
! NewCipher 创建并返回新的 Cipher。key 参数应该是 RC4 键，最少 1 个字节，最多 256 个字节
! rc4.Cipher 是使用特定密钥的 RC4 实例。
	XORKeyStream 将 dst 设置为 src 与密钥流进行异或的结果。Dst 和 src 必须完全重叠或完全不重叠。
*/

func TestRC4(t *testing.T) {

	plaintext := []byte("Hello World")

	key, c := RC4Encrypt(t, plaintext)

	plaintext = RC4Decrypt(t, key, c)
	logfln("%s", plaintext)
}

func RC4Encrypt(t *testing.T, plaintext []byte) (key, ciphertext []byte) {
	if len(plaintext) == 0 {
		return nil, nil
	}
	var l = len(plaintext)
	if len(plaintext) > 256 {
		l = 256
	}

	key = make([]byte, l)
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	ciphertext = make([]byte, l)
	cipher.XORKeyStream(ciphertext, plaintext)
	return
}
func RC4Decrypt(t *testing.T, key, ciphertext []byte) (plaintext []byte) {
	if len(key) < 1 || len(key) > 256 {
		t.Fatal("key is invalid")
	}
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	plaintext = make([]byte, len(ciphertext))
	cipher.XORKeyStream(plaintext, ciphertext)
	return
}
