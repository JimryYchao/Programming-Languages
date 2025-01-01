package gostd

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"testing"
)

/*
! GenerateKey 使用 rand 生成公钥和私钥对；rand 为 nil 则使用 crypto.rand.Reader;
	此函数的输出是确定性的，等效于从 rand 中读取 SeedSize 字节，并将其传递给 NewKeyFromSeed
! Sign 用私钥对 message 进行签名；如果 len(privateKey) 不是 PrivateKeySize, 则 panic
! Verify 用公钥检查 sig 是否为 message 的有效签名。如果 len(publicKey) 不是 PublicKeySize, 则 panic
! PrivateKey 是 Ed25519 私钥的类型。它实现了 crypto.Signer
	! NewKeyFromSeed 从提供的 seed 计算私钥，如 len(seed)）不是 SeedSize，它会死 panic
	Equal, Public, Seed
	Sign 用 priv.rand 对给定的 msg 进行签名，rand 被忽略，可以为 nil。可以传递 SHA-512 哈希后的 msg_512 和 opts{crypto.SHA512} 进行签名
! PublicKey 是 Ed25519 公钥的类型。
	Equal
!  VerifyWithOptions 使用 pub 和 opts 对  msg 和 sig 进行验证；
	当 VerifyWithOptions 的 opts.Hash 是 crypto.SHA512，则使用预哈希变体 Ed25519ph，并且 message 预计是 SHA-512 哈希，
	否则 opts.Hash 必须是 crypto.Hash(0) 并且 message 不得被哈希，因为 Ed25519 执行两次传递以签名 message。
*/

func TestED25519(t *testing.T) {
	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("The quick brown fox jumps over the lazy dog")

	//? Sign
	sig := ed25519.Sign(pri, msg)

	//? Verify
	if !ed25519.Verify(pub, msg, sig) {
		t.Fail()
	}
}

func TestPriSign(t *testing.T) {
	pub, pri, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("The quick brown fox jumps over the lazy dog")

	//? Sign
	sig, err := pri.Sign(nil, msg, &ed25519.Options{
		Context: "Example_ed25519ctx",
	})
	if err != nil {
		t.Fatal(err)
	}

	//? Verify
	if err := ed25519.VerifyWithOptions(pub, msg, sig, &ed25519.Options{
		Context: "Example_ed25519ctx",
	}); err != nil {
		t.Fatal("invalid signature")
	}
}

func TestED25519_SHA512(t *testing.T) {
	msg_sha512 := sha512.Sum512([]byte("The quick brown fox jumps over the lazy dog"))

	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	//? Sign
	sig, err := pri.Sign(nil, msg_sha512[:], &ed25519.Options{Hash: crypto.SHA512})
	if err != nil {
		t.Fatal(err)
	}

	//? Verify
	if err := ed25519.VerifyWithOptions(pub, msg_sha512[:], sig, &ed25519.Options{Hash: crypto.SHA512}); err != nil {
		t.Fatal(err)
	}
}
