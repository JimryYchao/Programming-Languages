package gostd

import (
	"crypto"
	"crypto/hmac"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"

	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
	_ "golang.org/x/crypto/sha3"

	"hash"
	"io"
	"testing"
)

/*
! Equal 比较两个 MAC 的相等性，而不泄漏定时信息。
! New 使用给定的 hash 和 key 返回一个新的 HMAC 哈希值；h 可以是例如 sha256.New;
	h 在每次调用时都必须返回一个新的 Hash；New 返回 Hash 不实现 encoding.BinaryMarshaler,BinaryUnmarshaler.
*/

func TestHMAC(t *testing.T) {
	key := []byte("Hello World")
	Sum := func(s string, h hash.Hash) {
		io.WriteString(h, "JimryYchao")
		logfln("%20s sum(`JimryYchao`) = %x", s, h.Sum(nil))
	}

	for _, h := range hashs {
		hnew := h.New()
		Sum(h.String(), hnew)

		hmach := hmac.New(h.New, key)
		Sum("hmac+"+h.String(), hmach)
	}
}

var hashs = []crypto.Hash{
	crypto.MD5,
	crypto.SHA1,
	crypto.SHA224,
	crypto.SHA256,
	crypto.SHA384,
	crypto.SHA512,
	crypto.SHA3_224,
	crypto.SHA3_256,
	crypto.SHA3_384,
	crypto.SHA3_512,
	crypto.SHA512_224,
	crypto.SHA512_256,
	crypto.BLAKE2s_256,
	crypto.BLAKE2b_256,
	crypto.BLAKE2b_384,
	crypto.BLAKE2b_512,
}
