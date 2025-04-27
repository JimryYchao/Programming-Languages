package gostd

import (
	"crypto"
	md5 "crypto/md5"
	"testing"
)

/*
! RegisterHash 注册一个函数，该函数返回给定哈希函数的新实例；它将从实现哈希函数的包的 init 中调用并注册
! crypto.Hash 标识一个在其他包中实现的加密哈希函数。
	New 使用给定的哈希函数返回一个 hash.Hash 值。如果相应的哈希函数未链接到二进制中，New 将出现异常。
	Available 报告给定的哈希函数是否链接到二进制中。
	HashFunc 仅返回 h 的值，这样的 Hash 就是实现了 SignerOpts
	Size 返回由给定哈希函数生成的摘要的长度（以字节为单位）。它不需要将该哈希函数链接到程序中。
	String 返回对应哈希函数字符串形式
! crypto.PrivateKey 表示未指明算法的私钥；
	Public() crypto.PublicKey
	Equal(x crypto.PrivateKey) bool
! crypto.PublicKey 表示使用未指明算法的公钥
	Equal(x crypto.PublicKey) bool
! crypto.Signer 是不透明私钥的接口，可用于签名操作。例如，保存在硬件模块中的 RSA 密钥。
	Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error)
	   - 使用私钥对 digest 进行签名，可能使用 rand 中的 entropy（）。对于 RSA 密钥，它生成的签名是 PKCS #1 v1.5 或 PSS 签名（由 opts）
		对于 (EC)DSA 密钥，它是一个 DER-serialised ASN.1 的签名结构
	   - Hash 实现了 SignerOpts 接口，因此可以作为 opts 使用的哈希函数
       - 当需要对较大信息的哈希进行签名时，调用方负责对较大信息进行哈希算法，并将 hash 作为 digest 和 opts 传入 Sign
! crypto.Decrypter 是不透明私钥的接口，可用于非对称解密操作。例如，保存在硬件模块中的 RSA 密钥。
	Decrypt(rand io.Reader, msg []byte, opts DecrypterOpts) (plaintext []byte, err error)
		Decrypt 解密 msg，opts 表明使用的原语, 参照各个加密算法实现的文档
*/

func TestHash(t *testing.T) {
	if crypto.MD5.Available() {
		md5Hash := crypto.MD5
		logfln("Register Hash Func: %s, Size:%d", md5Hash, md5Hash.Size())
	}
	// fs := os.DirFS("../")
	for _, msg := range msgs {
		logfln("%11s:%X", msg, md5.Sum(msg))
	}
}
