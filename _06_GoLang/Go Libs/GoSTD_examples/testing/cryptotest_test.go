package gostd_testing

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"testing/cryptotest"
)

/*
! cryptotest.SetGlobalRandom
cryptotest 在加密测试中的实际应用场景：
- 确定性 RSA 密钥生成测试
- 确定性 AES 加密测试
- 确定性随机字节生成测试
- 可重复的加密算法验证
*/

// FuzzRSAKeyGeneration 测试确定性 RSA 密钥生成, 使用相同的种子应该生成相同的密钥对
func FuzzRSAKeyGeneration(f *testing.F) {
	f.Fuzz(func(t *testing.T, seed uint64) {
		// 生成第一个密钥对
		cryptotest.SetGlobalRandom(t, seed)
		key1, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("Failed to generate first key pair: %v", err)
		}

		// 生成第二个密钥对
		cryptotest.SetGlobalRandom(t, seed)
		key2, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("Failed to generate second key pair: %v", err)
		}

		// 验证密钥是否相同
		if !key1.Equal(key2) {
			t.Fatal("Key pairs should be identical with same seed")
		}
		if !key1.PublicKey.Equal(&key2.PublicKey) {
			t.Fatal("Public keys should be identical with same seed")
		}
	})
}
