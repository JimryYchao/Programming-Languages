### Go Standard libraries     <a href="https://pkg.go.dev/std" target="_blank"><img src="./_rsc/link-src.drawio.png"/></a>

- 官网链接 <img src="./_rsc/link-src.drawio.png"/> 
- 补充说明  <img  src="./_rsc/link-others.drawio.png"/>
- 代码  <img src="./_rsc/link-code.drawio.png"/>
- 示例  <img src="./_rsc/link-exam.drawio.png"/>

---

- [x] bufio 包装了 `io.Reader` 和 `io.Writer` 并提供了缓冲和一些文本 I/O 帮助。<a href="https://pkg.go.dev/bufio" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./bufio/code/bufio_test.go"   ><img src="./_rsc/link-code.drawio.png"   
  id="exam"/></a><a href="bufio/bufio.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

- [x] bytes 实现了操作字节切片的函数。它类似于 strings 包的功能。       <a href="https://pkg.go.dev/bytes" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./bytes/code/bytes_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] cmp 提供了与比较有序值相关的类型和函数。        <a href="https://pkg.go.dev/cmp" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./cmp/cmp.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./cmp/code/cmp_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] compress 提供一些压缩和解压缩的工具包。 <a href="https://pkg.go.dev/compress" target="_blank"><img src="./_rsc/link-src.drawio.png" id="other"/></a><a href="compress/compress.md"  ><img  src="./_rsc/link-others.drawio.png" /></a>
 
  - [x] bzip2 实现 bzip2 解压缩。        <a href="https://pkg.go.dev/compress/bzip2" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="compress/code/bzip2_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] flate 实现了 [RFC 1951](https://rfc-editor.org/rfc/rfc1951.html) 中描述的 DEFLATE 压缩数据格式。        <a href="https://pkg.go.dev/compress/flate" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="compress/code/flate_test.go"   ><img src="./_rsc/link-code.drawio.png"   
  id="exam"/></a><a href="compress/examples/flate_netconnect.go"   ><img src="./_rsc/link-exam.drawio.png"  /></a>
  
  - [x] gzip 包实现了对 gzip 格式压缩文件的读取和写入。        <a href="https://pkg.go.dev/compress/gzip" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="compress/code/gzip_test.go"   ><img src="./_rsc/link-code.drawio.png"  
  id="exam"/></a><a href="compress/examples/gzip_httpSend.go"   ><img src="./_rsc/link-exam.drawio.png"  /></a>
  
  - [x] zlib 实现了对 zlib 格式压缩数据的读取和写入。        <a href="https://pkg.go.dev/compress/zlib" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="compress/code/zlib_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
  
  - [x] lzw 实现 Lempel-Ziv-Welch 压缩数据格式。      <a href="https://pkg.go.dev/compress/lzw" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="compress/code/lzw_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] container 提供一些容器数据类型。        <a href="https://pkg.go.dev/container" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="container/container.md"  ><img  src="./_rsc/link-others.drawio.png" /></a>

  - [x] heap 为实现 `heap.Interface` 的任何类型提供堆操作。        <a href="https://pkg.go.dev/container/heap" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="container/code/heap_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="container/container.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

  - [x] list 实现了双向链表。        <a href="https://pkg.go.dev/container/list" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="container/code/list_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] ring 实现了循环列表上的操作。        <a href="https://pkg.go.dev/container/ring" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="container/code/ring_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] context 定义了 Context 上下文类型，它携带 *deadlines*、*cancellation signals* 和跨 API 边界和进程之间的其他 *request-scoped  values*。       <a href="https://pkg.go.dev/context" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./context/context.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./context/code/context_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./context/context.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

- [x] crypto 包含一些通用的加密常量和算法。        <a href="https://pkg.go.dev/crypto" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="crypto/crypto.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="crypto/code/crypto_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  /></a>

  - [x] aes 实现 AES 加密。        <a href="https://pkg.go.dev/crypto/aes" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/aes_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./context/context.md#exam"   ><img src="./_rsc/link-exam.drawio.png"    /></a>

  - [x] cipher 实现了标准的分组密码模式，这些模式可以封装在低级分组密码的实现中。        <a href="https://pkg.go.dev/crypto/cipher" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/cipher_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] des 实现 Data Encryption Standard (DES) 和 Triple Data Encryption Algorithm (TDEA)。        <a href="https://pkg.go.dev/crypto/des" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/des_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] ecdh 实现了基于 NIST 曲线和 Curve25519 的 Elliptic Curve Diffie-Hellman (ECDH)。       <a href="https://pkg.go.dev/crypto/ecdh" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/ecdh_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] ecdsa 实现了在 FIPS 186-4 和 SEC 1 2.0 中定义的椭圆曲线数字签名算法。      <a href="https://pkg.go.dev/crypto/ecdsa" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/ecdsa_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] ed25519 实现了 [Ed25519](https://ed25519.cr.yp.to/) 签名算法。        <a href="https://pkg.go.dev/crypto/ed25519" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/ed25519_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] elliptic 实现了素数域上的标准 NIST P-224、P-256、P-384 和 P-521 椭圆曲线。        <a href="https://pkg.go.dev/crypto/elliptic" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/elliptic_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] hmac 实现了 Keyed-Hash Message Authentication Code (HMAC)。       <a href="https://pkg.go.dev/crypto/hmac" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/hmac_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] md5 实现了 [RFC 1321](https://www.rfc-editor.org/rfc/rfc1321.html) 中定义的 MD5 哈希算法。         <a href="https://pkg.go.dev/crypto/md5" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/md5_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] rand 实现了一个加密安全的随机数生成器。       <a href="https://pkg.go.dev/crypto/rand" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/rand_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] rc4 实现了 RC4加密。        <a href="https://pkg.go.dev/crypto/rc4" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/rc4_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] rsa 实现了 PKCS #1 和 [RFC 8017](https://www.rfc-editor.org/rfc/rfc8017.html) 中指定的 RSA 加密。       <a href="https://pkg.go.dev/crypto/rsa" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/rsa_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] sha1 实现了 [RFC 3174](https://www.rfc-editor.org/rfc/rfc3174.html) 中定义的 SHA-1 哈希算法。        <a href="https://pkg.go.dev/crypto/sha1" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/sha1_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] sha256 实现了 FIPS 180-4 中定义的 SHA224 和 SHA256 哈希算法。        <a href="https://pkg.go.dev/crypto/sha256" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/sha256_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] sha512 实现了 FIPS 180-4 中定义的 SHA-384、SHA-512、SHA-512/224 和 SHA-512/256 哈希算法。        <a href="https://pkg.go.dev/crypto/sha512" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/sha512_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] subtle 实现了在加密代码中经常有用但需要仔细考虑才能正确使用的函数。        <a href="https://pkg.go.dev/crypto/subtle" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/subtle_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [ ] tls 部分实现了 [RFC 5246](https://www.rfc-editor.org/rfc/rfc5246.html) 中指定的 TLS 1.2 和 [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446.html) 中指定的 TLS 1.3。      <a href="https://pkg.go.dev/crypto/tls" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/tls_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [ ] x509 实现了 X.509 标准的一个子集。        <a href="https://pkg.go.dev/crypto/x509" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/x509_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [ ] x509/pkix 包含用于 ASN.1 解析和序列化 X.509 证书、CRL 和 OCSP 的共享低级结构。        <a href="https://pkg.go.dev/crypto/x509/pkix" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="crypto/code/pkix_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] debug 包含一些有关调试的包。       <a href="https://pkg.go.dev/debug" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="debug/debug.md"  ><img  src="./_rsc/link-others.drawio.png"   /></a>

  - [x] buildinfo 提供了访问 Go 二进制文件中嵌入的关于如何构建它的信息。       <a href="https://pkg.go.dev/debug/buildinfo" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/buildinfo_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [ ] dwarf 提供对从可执行文件加载的 DWARF 调试信息的访问       <a href="https://pkg.go.dev/debug/dwarf" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/dwarf_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
  
  - [ ] elf 实现了对 ELF 对象文件的访问。       <a href="https://pkg.go.dev/debug/elf" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/elf_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
 
  - [ ] gosym 包实现了对由 gc 编译器生成的 Go 二进制文件中嵌入的 Go 符号和行号表的访问。       <a href="https://pkg.go.dev/debug/gosym" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/gosym_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>  
 
  - [ ] macho 实现对 Mach-O 对象文件的访问。       <a href="https://pkg.go.dev/debug/macho" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/macho_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
 
  - [ ] pe 实现对 PE（Microsoft Windows 可移植可执行文件）文件的访问。       <a href="https://pkg.go.dev/debug/pe" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/pe_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
 
  - [ ] plan9obj 实现了对 Plan9 a.out 对象文件的访问。       <a href="https://pkg.go.dev/debug/plan9obj" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="debug/code/plan9obj_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] embed 提供了对嵌入在运行的 Go 程序中的文件的访问。      <a href="https://pkg.go.dev/embed" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./embed/embed.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./embed/code/embed_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] errors 实现一些函数来处理错误。       <a href="https://pkg.go.dev/errors" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./errors/code/errors_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] expvar 提供了一个公共变量的标准化接口。        <a href="https://pkg.go.dev/expvar" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="expvar/expvar.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="expvar/code/expvar_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>


- [x] flag 实现命令行标志解析。     <a href="https://pkg.go.dev/flag" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./flag/flag.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./flag/code/flag_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] fmt 使用类 C 的 `printf` 和 `scanf` 的函数实现格式化 I/O。“*verbs*” 格式从 C 派生的。       <a href="https://pkg.go.dev/" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./fmt/fmt.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./fmt/code/fmt_test.go"><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] index/suffixarray 使用内存中的后缀数组，在对数时间内实现子字符串搜索        <a href="https://pkg.go.dev/index/suffixarray" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="#"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="index/code/suffixarray_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] io 提供 I/O 原语的基本接口。       <a href="https://pkg.go.dev/io" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./io/code/io_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] fs 定义了文件系统的基本接口。<a href="https://pkg.go.dev/io/fs" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./io/code/fs_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] log 实现了一个简单的日志记录包。       <a href="https://pkg.go.dev/log" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./log/code/log_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./log/log.md#exam.md"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

  - [x] syslog 为系统日志服务提供了一个简单的接口。它可以使用 UNIX 域套接字、UDP 或 TCP 向系统日志守护程序发送消息。       <a href="https://pkg.go.dev/log/syslog"  target="_blank"><img src="./_rsc/link-src.drawio.png" /></a>

  - [x] slog 提供结构化日志记录，其中日志记录包括消息、严重性级别和以键值对表示的各种其他属性。       <a href="https://pkg.go.dev/log/slog" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./log/slog.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./log/code/slog_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./log/slog.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

- [x] maps 定义了各种对任何类型的映射的辅助函数。      <a href="https://pkg.go.dev/maps"  target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./maps/code/maps_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] math 提供基本常量和数学函数。        <a href="https://pkg.go.dev/math" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./math/code/math_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] big 实现了任意精度的算术（大数字）。        <a href="https://pkg.go.dev/math/big" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./math/big.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./math/code/big_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./math/big.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

  - [x] bits 为预先声明的无符号整数类型实现位计数和操作函数。        <a href="https://pkg.go.dev/math/bits" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./math/code/bits_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] cmplx 为复数提供基本常数和数学函数。        <a href="https://pkg.go.dev/math/cmplx" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./math/code/cmplx_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] rand 实现了适合模拟等任务的伪随机数生成器。       <a href="https://pkg.go.dev/math/rand" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./math/rand.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./math/code/rand_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
      - [x] rand/v2 实现了适合模拟等任务的伪随机数生成器。       <a href="https://pkg.go.dev/math/rand/v2" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./math/code/randv2_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] mime 实现 MIME 规范的一部分。        <a href="https://pkg.go.dev/mime" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="mime/code/mime_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] multipart 实现 RFC 2046 中定义的 MIME 多分部分析        <a href="https://pkg.go.dev/mime/multipart" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="mime/code/multipart_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] quotedprintable 实现 RFC 2045 指定的 quoted-printable 编码。        <a href="https://pkg.go.dev/mime/quotedprintable" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="mime/code/quotedprintable_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] os 为操作系统功能提供独立于平台的接口。        <a href="https://pkg.go.dev/os" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="os/code/os_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] exec 运行外部命令。        <a href="https://pkg.go.dev/os/exec" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="os/exec.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="os/code/exec_test.go"   ><img src="./_rsc/link-code.drawio.png" /></a>

  - [x] signal 实现对输入信号的访问。        <a href="https://pkg.go.dev/os/signal" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="os/signal.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="os/code/signal_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] user 允许按 name 或 id 查找用户帐户。        <a href="https://pkg.go.dev/os/user" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="os/code/user_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] path 实现用于操作斜杠分隔路径的实用程序例程。        <a href="https://pkg.go.dev/path" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="path/code/path_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] filepath 实现了用于以与目标操作系统定义的文件路径兼容的方式，操作文件名路径的实用程序例程。        <a href="https://pkg.go.dev/path/filepath" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="path/code/filepath_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] plugin 实现 Go 插件的加载和 symbol 符号解析。        <a href="https://pkg.go.dev/plugin" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="plugin/plugin.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="plugin/code/plugin_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  
<!-- TODO -------------------------------------------------------------------------------------------------------------------->

- [x] reflect 实现运行时反射，允许程序操作具有任意类型的对象。        <a href="https://pkg.go.dev/reflect" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="reflect/reflect.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./reflect/code/reflect_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="reflect/reflect.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

- [x] slices 定义了对任何类型的切片的辅助函数。      <a href="https://pkg.go.dev/slices"  target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./slices/code/slices_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] sort 提供了用于对切片和用户定义的集合进行排序的原语。       <a href="https://pkg.go.dev/sort" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="sort/code/sort_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="sort/sort.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>


- [x] strconv 实现了基本数据类型的字符串表示形式之间的转换。        <a href="https://pkg.go.dev/strconv" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./strconv/strconv.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./strconv/code/strconv_test.go"   ><img src="./_rsc/link-code.drawio.png"
  id="exam"/></a><a href="strconv/strconv.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>



- [x] strings 实现了一些函数来操作 UTF-8 编码的字符串。      <a href="https://pkg.go.dev/strings "  target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./strings/strings.md"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] sync 提供基本的同步原语。        <a href="https://pkg.go.dev/sync" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./sync/code/sync_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./sync/sync.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>

  - [x] atomic 提供了用于实现同步算法的低级原子内存原语。        <a href="https://pkg.go.dev/sync/atomic" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./sync/atomic.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./sync/code/atomic_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [ ] syscall 包含一个到低级操作系统原语的接口。参见 [sys](../_03_Go%20thrid-party/sys/sys.md)       <a href="https://pkg.go.dev/syscall" target="_blank"><img src="./_rsc/link-src.drawio.png" /></a>

  - [ ] TODO


- [x] testing 为 Go 包提供自动化测试支持。<a href="https://pkg.go.dev/testing" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./testing/testing.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./testing/code/testing_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
  
  - [x] fstest 实现了对测试实现和文件系统用户的支持。      <a href="https://pkg.go.dev/testing/fstest" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./testing/code/fstest_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>
  
  - [x] iotest 实现了主要用于测试的 Readers 和 Writers。        <a href="https://pkg.go.dev/testing/iotest" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./testing/code/iotest_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] quick 实现了一些实用函数来帮助进行黑盒测试。        <a href="https://pkg.go.dev/testing/quick" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./testing/code/quick_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] slogtest 实现了对 `log/slog.Handler` 的测试实现的支持。       <a href="https://pkg.go.dev/testing/slogtest" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./testing/code/slogtest_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>


- [x] time 提供测量和显示时间的功能。      <a href="https://pkg.go.dev/time" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./time/time.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./time/code/time_test.go"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="./time/time.md#exam"   ><img src="./_rsc/link-exam.drawio.png"  /></a>


- [x] unsafe 包含绕过 Go 程序类型安全的操作。       <a href="https://pkg.go.dev/unsafe" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="./unsafe/unsafe.md"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="./unsafe/code/unsafe_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

- [x] unicode 提供一些数据和函数来测试 Unicode 码位的某些属性。      <a href="https://pkg.go.dev/unicode" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="exam"/></a><a href="./unicode/code/unicode_test.go"   ><img src="./_rsc/link-exam.drawio.png"  /></a>
  
  - [x] utf8 实现了支持以 UTF-8 编码的文本函数和常量。     <a href="https://pkg.go.dev/unicode/utf8" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./unicode/code/utf8_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

  - [x] utf16 实现了 UTF-16 序列的编码和解码。       <a href="https://pkg.go.dev/unicode/utf16" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="code"/></a><a href="./unicode/code/utf16_test.go"   ><img src="./_rsc/link-code.drawio.png"   /></a>

<!-- TODO -------------------------------------------------------------------------------------------------------------------->

<!-- 

- [ ]         <a href="https://pkg.go.dev/#" target="_blank"><img src="./_rsc/link-src.drawio.png" 
  id="other"/></a><a href="#"  ><img  src="./_rsc/link-others.drawio.png" 
  id="code"/></a><a href="#"   ><img src="./_rsc/link-code.drawio.png" 
  id="exam"/></a><a href="#exam"   ><img src="./_rsc/link-exam.drawio.png"
  /></a>

--> 



---
### Go supplemental libraries       <a href="https://pkg.go.dev/std" target="_blank"><img src="./_rsc/link-src.drawio.png"/></a>


---