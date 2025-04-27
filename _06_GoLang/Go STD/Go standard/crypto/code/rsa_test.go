package gostd

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"strings"
	"testing"
)

/*
! GenerateKey 生成给定位大小的随机 RSA 私钥;
	bit 的大小不能低于 (len(msg)+hash.Size*2 +2)*8 以至少保证 msg 不会超过 pub 公钥模数 (PublicKey.Size) 的长度减去 hash 长度的两倍，再减去 2
! rsa.PrivateKey 代表 RSA 私钥；
	Equal, Public
	Precompute 进行一些预计算，以加快未来的私钥操作; 调用后 PrivateKey.Precomputed 不可更改，它由 Precompute() 生成
	Decrypt: opts 为 nil 或 *PKCS1v15DecryptOptions 时，执行 PKCS #1 v1.5 解密；否则 opts 必须具有类型 *OAEPOptions 并进行 OAEP 解密
		rand 仅在 opts 为 nil 或 *OAEPOptions 时可以忽略
	Sign: opts 是 *PSSOptions，使用 PSS 算法; 否则使用 PKCS #1 v1.5。digest 必须是使用 opts.HashFunc 对 message 哈希后的结果
	Validate 对私钥执行基本的健全性检查，检查 priv 是否有效
! rsa.PublicKey 代表 RSA 公钥
	Equal
	Size 返回以字节为单位的模数 pub.N 大小。原始签名和密文的大小与此公钥相同。
! EncryptOAEP 使用 RSA-OAEP 加密给定的 msg。rand 参数中作为熵的来源，确保对同一 msg 加密两次不会产生相同的密文
	label 参数可以包含任意数据，这些数据不会被加密，但会为 msg 提供重要的上下文。例如，如果给定的公钥用于加密两种类型的消息，
	则可以使用不同的 label 来确保用于一种目的的密文不能被攻击者用于另一种目的。如果不需要，可以为空。
	消息的长度不能超过公钥模数 (PublicKey.Size) 的长度减去 hash 长度的两倍，再减去 2。
! DecryptOAEP 使用 RSA-OAEP 解密密文; 给定消息的加密和解密必须使用相同的哈希函数, 例如 sha256.New
	参数 random 可以忽略；label 参数必须与加密时给定的值匹配
! OAEPOptions 是一个接口，用于使用 crypto.Decrypter 接口向 OAEP 解密传递 opts。
*/

func TestOAEP(t *testing.T) {
	secretMessage := []byte("send reinforcements, we're going to advance")
	label := []byte("orders")
	pri, err := rsa.GenerateKey(rand.Reader, (len(secretMessage)+sha256.Size*2+2)*8)
	if err != nil {
		t.Fatal(err)
	}
	// Encrypt
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pri.PublicKey, secretMessage, label)
	if err != nil {
		t.Fatalf("Error from encryption: %s", err)
	}
	logfln("Ciphertext: %x", ciphertext)

	// rsa.DecryptOAEP Decrypt
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pri, ciphertext, label)
	if err != nil {
		t.Fatal(err)
	}
	logfln("rsa.DecryptOAEP Plaintext: %s", plaintext)

	// pri.Decrypt( OAEPOptions )
	plaintext2, err := pri.Decrypt(rand.Reader, ciphertext, &rsa.OAEPOptions{Label: label, Hash: crypto.SHA256})
	if err != nil {
		t.Fatal(err)
	}
	logfln("pri.Decrypt(OAEPOptions) Plaintext: %s", plaintext2)
}

/*
! EncryptPKCS1v15 使用 RSA 和 PKCS #1 v1.5 中的填充方案对给定消息进行加密。消息长度不能超过 pub 模数减去 11 字节。rand 作为熵
	加密：使用此函数加密 session keys 密钥以外的明文是危险的。在新协议中使用 RSA OAEP。
! DecryptPKCS1v15, DecryptPKCS1v15SessionKey 使用 RSA 和 PKCS #1 v1.5 中的填充方案解密 ciphertext;
	random 可以忽略；DecryptPKCS1v15 无论是否返回错误，都会泄露机密信息; see DecryptPKCS1v15SessionKey
	DecryptPKCS1v15SessionKey 使用一个 session key 来保存解密后的明文；如果密文的长度过长则返回错误
		用户可以预先生成一个随机 session key，如果 key 太小，攻击者可鞥会进行暴力破解，
		至少使用 16 字节的 key 可以防止这种攻击。
*/

func TestPKCS1v15(t *testing.T) {
	secretMessage := []byte("send reinforcements, we're going to advance")
	pri, _ := rsa.GenerateKey(rand.Reader, 1024)
	if pri.Validate() != nil {
		t.Fatal("privateKey is invalid")
	}

	// EncryptPKCS1v15
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &pri.PublicKey, secretMessage)
	if err != nil {
		t.Fatal(err)
	}
	logfln("%x", ciphertext)

	// DecryptPKCS1v15
	plaintext, _ := rsa.DecryptPKCS1v15(rand.Reader, pri, ciphertext)
	logfln("%s", plaintext)

	// DecryptPKCS1v15SessionKey
	sessionKey := make([]byte, len(secretMessage))
	rand.Read(sessionKey)
	err = rsa.DecryptPKCS1v15SessionKey(rand.Reader, pri, ciphertext, sessionKey)
	if err != nil {
		t.Fatal(err)
	}
	logfln("%s", sessionKey)

	// pri.Decrypt( PKCS1v15DecryptOptions )
	pri.Precompute()
	mess, _ := pri.Decrypt(rand.Reader, ciphertext, &rsa.PKCS1v15DecryptOptions{SessionKeyLen: len(sessionKey)})
	logfln("%s", mess)
}

/*
! SignPKCS1v15 为给定 hashed 签名；hashed 必须是由 Hash 函数运算的结果；random 可以忽略
! VerifyPKCS1v15 验证 RSA PKCS #1 v1.5 签名
! SignPSS 为 digest 签名；digest 必须是由 Hash 函数运算的结果；opts 可以为 nil
! VerifyPSS 验证 PSS 签名
! PSSOptions 包含用于创建和验证 PSS 签名的选项。
*/
func TestSignPKCS1v15(t *testing.T) {
	message := []byte("send reinforcements, we're going to advance")
	pri, _ := rsa.GenerateKey(rand.Reader, 1024)

	hashed := sha256.Sum256(message)
	logfln("message hashed : %x", hashed)

	// Sign
	sig, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed[:])
	if err != nil {
		t.Fatal(err)
	}
	logfln("sig : %x", sig)

	// Verify
	err = rsa.VerifyPKCS1v15(&pri.PublicKey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		logfln("Verify failed")
	}
}

func TestSignPSS(t *testing.T) {
	message := []byte("send reinforcements, we're going to advance")
	pri, _ := rsa.GenerateKey(rand.Reader, 1024)

	h := sha256.New()
	io.Copy(h, strings.NewReader(string(message)))
	hashed := h.Sum(nil)
	logfln("message hashed : %x", hashed)

	// Sign
	sig, err := rsa.SignPSS(rand.Reader, pri, crypto.SHA256, hashed, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		t.Fatal(err)
	}
	logfln("sig : %x", sig)

	// Verify
	err = rsa.VerifyPSS(&pri.PublicKey, crypto.SHA256, hashed, sig, nil)
	if err != nil {
		logfln("Verify failed")
	}
}
