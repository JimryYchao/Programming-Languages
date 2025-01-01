package gostd

/*
! ecdh.Curve
	GenerateKey 生成一个随机 *PrivateKey；大多数程序应使用 crypto/rand.Reader
	NewPrivateKey 检查 key 是否有效并返回一个 *PrivateKey；
		对于 NIST 曲线，将 bytes 解码为固定长度的大端整数，并检查结果是否低于曲线的阶数；零值 key 被拒绝
		对于 X25519 只检查标量的长度
	NewPublicKey 检查 key 是否有效并返回一个 *PublicKey
		对于 NIST 曲线，将 key 解码为一个未压缩的点，压缩编码或无穷点被拒绝
		对于 X25519 只检查 u 坐标长度，对抗性选择的 key 可能导致 ECDH 返回错误
! P256 返回一个实现 NIST P-256 的 Curve，也称为 secp256r1 或 prime256v1
! P384 返回一个实现 NIST P-384 的 Curve，也称为 secp384r1
! P521 返回一个实现 NIST P-521 的 Curve，也称为 secp521r1
! X25519 返回一个实现了 Curve25519  X25519 的 Curve

! ecdh.PrivateKey 是 ECDH 私钥;
		可以用 crypto/x509.ParsePKCS8PrivateKey 解析，并用 crypto/x509.MarshalPKCS8PrivateKey 编码。
		对于 NIST 曲线，需要在解析后使用 crypto/ecdsa.PrivateKey.ECDH 进行转换。
	Bytes 返回私钥的副本
	Curve, Equal
	ECDH 执行 ECDH 交换，并返回共享密钥；PrivateKey 和 PublicKey 必须使用相同的 Curve
		对于 NIST 曲线，结果永远不会是无穷远点。
		对于 X25519，如果结果为全零值，ECDH 将返回错误。
	Public 实现所有标准库私钥的隐式接口; see [crypto.PrivateKey]
! ecdh.PublicKey 是 ECDH 公钥;
		这些密钥可以用 crypto/x509.ParsePKIXPublicKey 解析，并用 crypto/x509.MarshalPKIXPublicKey 编码。
		对于 NIST 曲线，需要在解析后使用 crypto/ecdsa.PublicKey.ECDH 进行转换。
	Bytes, Curve
	Equal: 可以存在具有不同编码的等效公钥，这些公钥将从该 Equal 检查中返回 false，但行为方式与 ECDH 的输入相同
*/
