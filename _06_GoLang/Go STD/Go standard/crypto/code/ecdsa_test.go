package gostd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"
)

/*
! Sign, SignASN1 使用私钥对哈希 (对一个大信息的哈希结果) 进行签名; 哈希比私钥的 curve order 长时截断到长度。大多数程序使用 SignASN1 而不是直接处理 Sign 返回的 r,s
! Verify, VerifyASN1 使用公钥验证 hash 中的签名是否有效; 大多数程序应该使用 VerifyASN1
! GenerateKey 为指定的曲线 elliptic.Curve 生成一个新的 ECDSA 私钥
! ecdsa.PrivateKey 表示 ECDSA 私钥
	ECDH 将 k 作为 ecdh.PrivateKey 返回; 如果根据 ecdh.Curve.NewPrivateKey 的定义密钥无效，或者如果 crypto/ecdh 不支持该 Curve 时，则返回错误。
	PublicKey, Public() 返回与 PrivateKey 对应的公钥
	Sign 对 digest 进行签名, opts 参数目前不适用; 仅与 crypto.Signer 接口保持一致
! ecdsa.PublicKey 表示与 ECDSA 私钥对应的公钥
	ECDH 将 k 作为 ecdh.PublicKey 返回; 如果根据 ecdh.Curve.NewPublicKey 的定义密钥无效，或者如果 crypto/ecdh 不支持 Curve，则返回错误。
	PrivateKey.Equal, PublicKey.Equal 报告 p 和 x 是否具有相同的值。
*/

func TestECDSA(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	//? sign
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("signature: %x\n", sig)

	//? verify
	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	fmt.Println("signature verified:", valid)
}
