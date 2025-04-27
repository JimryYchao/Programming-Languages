package gostd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"testing"
)

/*
! MarshalCompressed 将 Curve 上的点转换为压缩形式，如果该点不在曲线上（或者是无穷远处的点），则行为未定义。
! UnmarshalCompressed 将 MarshalCompressed 压缩后的形式还原到 Curve 上的 x,y
! elliptic.Curve 表示 a=-3 的简式 Weierstrass 曲线。除了 P224、P256、P384 和 P521 返回的曲线实现外，不推荐实现 Curve
	! P224, P256, P384, P521 返回的 Curve 作为 crypto/ecdsa.GenerateKey 的 curve
	Params 返回曲线的参数 CurveParams
*/

func TestElliptic(t *testing.T) {
	pri, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader) // P224、P256、P384、P521
	if err != nil {
		t.Fatal(err)
	}
	degist := []byte("Hello World")
	hash := sha512.Sum512(degist)

	// Sign
	sig, err := ecdsa.SignASN1(rand.Reader, pri, hash[:])
	if err != nil {
		t.Fatal(err)
	}

	// Verify
	if !ecdsa.VerifyASN1(&pri.PublicKey, hash[:], sig) {
		t.Fail()
	}
}
