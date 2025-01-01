package gostd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"testing"
)

/*
! NewCipher 创建并返回一个基于 ASE 密钥的 cipher.Block
	key 参数应该是 AES 密钥，可以是 16、24 或 32 字节，以选择 AES-128、AES-192 或 AES-256
*/

func TestAES(t *testing.T) {
	ciphertexts := make([][]byte, len(msgs))
	// Encrypt
	for i, msg := range msgs {
		c, err := OFBEncrypt(msg)
		if err != nil {
			t.Fatal(err)
		}
		ciphertexts[i] = c
	}

	// Decrypt
	for _, c := range ciphertexts {
		p, err := OFBDecrypt(c)
		if err != nil {
			t.Fatal(err)
		}
		logfln("%s", p)
	}
}

var pkeyString = "6368616e676520746869732070617373"

func OFBEncrypt(plaintext []byte) ([]byte, error) {
	key, _ := hex.DecodeString(pkeyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	encrypter := cipher.NewOFB(block, iv)
	encrypter.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func OFBDecrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	key, _ := hex.DecodeString(pkeyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	decrypter := cipher.NewOFB(block, ciphertext[:aes.BlockSize])
	decrypter.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return plaintext, nil
}
