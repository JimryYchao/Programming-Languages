package gostd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

/*
! cipher.AEAD 是一种密码模式，提供与相关数据相关的认证加密。
	NonceSize 返回必须传递给 Seal 和 Open 的随机数大小
	Overhead 返回明文和密文长度之间的最大差异
	Seal 对明文进行加密和身份验证，对 additionData 进行身份验证并将结果追加到 dst，返回更新后的 []byte
		对于给定的密钥，nonce 必须为 NonceSize() 的字节长度，并且在任何时候都是唯一的
		要重用明文存储的加密输出，使用 plaintext[:0] 作为 dst
	Open 对密文进行解密和身份验证，对 additionData 进行身份验证；若成功则将结果明文追加到 dst，返回更新后的 []byte
		nonce 必须为 NonceSize() 的字节长度，并且 nonce 和 addition 都必须于传递给 Seal 的值相匹配
! NewGCM 返回给定的 128 位, 封装在 Galois Counter 模式, 具有标准 nonce 长度的分组密码 AEAD
	一般来说，GCM 的这个实现所执行的 GHASH 操作不是恒定时间的。
! NewGCMWithNonceSize 类似于 NewGCM，但使用给定长度 size 的 nonce。长度不得为零。
! NewGCMWithTagSize 类似于 NewGCM，但生成具有给定长度 tagSize ([12,16] 之间) 的标记
*/

func TestGCMDecrypt(t *testing.T) {
	// 从安全的地方加载您的密钥，并在多个 Seal/Open 调用中重用它。如果要将密码短语转换
	// 为密钥，请使用合适的包，如 bcrypt 或 scrypt。解码后的密钥应该是 16 字节 (AES-128) 或 32 字节 (AES-256)。
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	block, err := aes.NewCipher(key) // 分组密码
	checkErr(err)
	aesgcm, err := cipher.NewGCM(block) // 加密模式
	checkErr(err)
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil) // 解密
	checkErr(err)

	logfln("%s", plaintext)
}

func TestGCMEncrypt(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	plaintext := []byte("exampleplaintext")
	block, _ := aes.NewCipher(key)
	// 对于给定的密钥，不要使用超过 2^32 个随机 nonces，因为有重复的风险。
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		t.Fatal(err.Error())
	}
	aesgcm, _ := cipher.NewGCM(block)
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, []byte("Hello World"))
	logfln("%x", ciphertext)

	// Decrypt
	plaintext, _ = aesgcm.Open(nil, nonce, ciphertext, []byte("Hello World")) // 使用相同的 nonce 和 additionData
	logfln("%s", plaintext)
}

/* 分组密码
! cipher.Block 表示使用给定密钥的分组密码的实现。它提供了加密或解密单个 block 的能力。
	BlockSize cipher's block size.
	Encrypt, Decrypt dst 到 src，它们要么完全重叠要么完全不重叠
! cipher.BlockMode 表示在 block-base 的模式（CBC，ECB 等）下运行的分组密码。
	BlockSize mode's block size
	CryptBlocks 加密或解密一些 block，src 长度必须是 BlockSize 的倍数；dst 和 src 要么完全重叠要么完全不重叠；
		可以接受 len(dst) > len(src)，只会更新 dst[:len(src)] 的部分
! NewCBCEncrypter 返回一个 BlockMode，它使用给定的 Block 在分组密码链接模式下加密。
	iv 的长度必须与 Block.BlockSize 相同
! NewCBCDecrypter 返回一个 BlockMode，它使用给定的 Block 的分组密码链接模式下解密。
	iv 的长度必须与 Block.BlockSize 相同，并且必须匹配用于加密数据 NewCBCEncrypter 的 iv。
*/

func TestCBCDecrypt(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	ciphertext, _ := hex.DecodeString("73c86d43a9d700a253a96c85b0f6b03ac9792e0e757f869cca306bd3cba1c62b")

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	// IV 需要是唯一的，但不安全。因此，通常将其包含在密文的开头。
	if len(ciphertext) < aes.BlockSize {
		t.Fatal("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC 模式总是在整个 Block 中运行
	if len(ciphertext)%aes.BlockSize != 0 { //? src 是 BlockSize 的倍数
		t.Fatal("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// 如果两个参数相同，CryptBlocks 可以相同位置运行。
	mode.CryptBlocks(ciphertext, ciphertext) //? src 和 dst 完全重叠

	// 如果原始明文长度不是 BlockSize 的倍数，则在加密时必须进行填充；解密后删除填充
	// 在解密之前必须对密文进行身份验证 (即通过使用 crypto/hmac)，以避免创建填充 oracle。
	logfln("%s", ciphertext) // exampleplaintext
}

func TestCBCEncrypt(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	// 假设明文是正确填充的，长度是 BlockSize 的倍数
	if len(plaintext)%aes.BlockSize != 0 {
		t.Fatal("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	// IV 需要是唯一的，但不安全。因此，通常将其包含在密文的开头。
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		t.Fatal(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// 为了安全，密文必须经过身份验证 (即通过使用 crypto/hmac) 以及加密。
	logfln("%x", ciphertext)

	// Decrypt
	dmode := cipher.NewCBCDecrypter(block, ciphertext[:aes.BlockSize]) // iv
	dmode.CryptBlocks(ciphertext[aes.BlockSize:], ciphertext[aes.BlockSize:])
	logfln("%s", ciphertext[aes.BlockSize:])
}

/* 流密码
! cipher.Stream 表示流密码
	XORKeyStream 将给定片中的每个字节与密码流中的一个字节; len: dst >= src
! NewCFBDecrypter 返回一个 Stream，它使用给定的 Block 以密码反馈模式进行解密。iv 的长度必须等于 Block 的块大小
! NewCFBEncrypter 返回一个 Stream，它使用给定的 Block 以密码反馈模式进行加密。iv 的长度必须等于 Block 的块大小
! NewCTR 返回一个 Stream，它在计数器模式下进行加密/解密。iv 的长度必须等于 Block 的块大小相同
! NewOFB 返回一个 Stream，它在输出反馈模式下进行加密或解密。iv 的长度必须等于 Block 的块大小
*/

func TestCFBDecrypt(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	ciphertext, _ := hex.DecodeString("7dd015f06bec7f1b8f6559dad89f4131da62261786845100056b353194ad")

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	// IV 需要是唯一的，但不安全。因此，通常将其包含在密文的开头。
	if len(ciphertext) < aes.BlockSize {
		t.Fatal("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// dst 和 src 可以相同
	stream.XORKeyStream(ciphertext, ciphertext)
	logfln("%s", ciphertext)
}

func TestCFBEncrypt(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373") // 32 bits aes
	plaintext := []byte("some plaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		t.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	logfln("%x", ciphertext)
}

func TestCTR(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("some plaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	// IV 需要是唯一的，但不安全。因此，通常将其包含在密文的开头。
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		t.Fatal(err)
	}
	// Encrypt
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Decrypt；CTR 模式对于加密和解密都是相同的，所以可以使用 stream 来解密。
	plaintext2 := make([]byte, len(plaintext))
	stream = cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])
	logfln("%s", plaintext2)
}

func TestOFB(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("some plaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		t.Fatal(err)
	}

	// Encrypt
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Decrypt；OFB 模式对于加密和解密都是相同的，所以可以使用 stream 来解密。
	plaintext2 := make([]byte, len(plaintext))
	stream = cipher.NewOFB(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])
	logfln("%s", plaintext2)
}

// ! cipher.StreamReader 包装 Stream 到 io.Reader 中，它调用 XORKeyStream 来处理 Read 的每个字节切片
// ! cipher.StreamWriter 包装 Stream 到 io.Writer 中；
func TestStreamReader(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	encrypted, _ := hex.DecodeString("cf0495cc6f75dafc23948538e79904a9")
	bReader := bytes.NewReader(encrypted) // bytes Reader

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	//? Encrypt
	// 如果每个密文的 key 都是唯一的，那么可以使用零 IV。
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: bReader}
	// 将输入复制到输出流中，同时进行解密。
	io.Copy(os.Stdout, reader) // some secret text

	//? Decrypt
	bReader2 := bytes.NewReader([]byte("some secret text"))
	stream2 := cipher.NewOFB(block, iv[:])
	var out bytes.Buffer
	writer := &cipher.StreamWriter{S: stream2, W: &out}
	// 将输入复制到输出流中，同时进行加密。
	io.Copy(writer, bReader2)
	logfln("%x", out.Bytes()) // cf0495cc6f75dafc23948538e79904a9
	writer.Close()
}
