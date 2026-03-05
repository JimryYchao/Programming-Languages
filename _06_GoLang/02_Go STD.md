### Go STD (Go 1.26.1)

| Package | Description | Example |
 | :--- | :--- | :--- |
| archive | 压缩存储访问 | 
| &nbsp;&nbsp;&nbsp;&nbsp;[tar](https://pkg.go.dev/archive/tar) | 对 tar 压缩包的访问 | [[↗]](./Go%20Libs/GoSTD_examples/archive/tar_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[zip](https://pkg.go.dev/archive/zip) | zip 压缩包的读写 | [[↗]](./Go%20Libs/GoSTD_examples/archive/zip_test.go) |
| [bufio](https://pkg.go.dev/bufio) | 实现缓冲 I/O，包装 io.Reader 或 io.Writer 对象 | [[↗]](./Go%20Libs/GoSTD_examples/bufio/bufio_test.go) |
| [builtin](https://pkg.go.dev/builtin) | 提供 Go 预声明标识符 |
| [bytes](https://pkg.go.dev/bytes) | 字节切片操作 | [[↗]](./Go%20Libs/GoSTD_examples/bytes/bytes_test.go) |
| [cmp](https://pkg.go.dev/cmp) | 有序值类型相关 |
| [compress](https://pkg.go.dev/compress) | 解压缩相关 |
| &nbsp;&nbsp;&nbsp;&nbsp;[bzip2](https://pkg.go.dev/compress/bzip2) | bzip2 解压 | [[↗]](./Go%20Libs/GoSTD_examples/compress/bzip2_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[flate](https://pkg.go.dev/compress/flate) | DEFLATE 压缩数据格式 | [[↗]](./Go%20Libs/GoSTD_examples/compress/flate_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[gzip](https://pkg.go.dev/compress/gzip) | gzip 压缩数据格式 | [[↗]](./Go%20Libs/GoSTD_examples/compress/gzip_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[lzw](https://pkg.go.dev/compress/lzw) | lzw 压缩数据格式 | [[↗]](./Go%20Libs/GoSTD_examples/compress/lzw_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[zlib](https://pkg.go.dev/compress/zlib) | zlib 压缩数据格式 | [[↗]](./Go%20Libs/GoSTD_examples/compress/zlib_test.go) |
| [container](https://pkg.go.dev/container) | 容器相关 | [[↗]](./Go%20Libs/GoSTD_examples/container/container_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[heap](https://pkg.go.dev/container/heap) | 堆 |
| &nbsp;&nbsp;&nbsp;&nbsp;[list](https://pkg.go.dev/container/list) | 双向链表 |
| &nbsp;&nbsp;&nbsp;&nbsp;[ring](https://pkg.go.dev/container/ring) | 循环链表 |
| [context](https://pkg.go.dev/context) | 定义 Context 类型，携带截止时间、取消信号和其他请求范围的值 | [[↗]](./Go%20Libs/GoSTD_examples/context/context_test.go) |
| [crypto](https://pkg.go.dev/crypto) | 收集常见的加密常量 | [[↗]](./Go%20Libs/GoSTD_examples/crypto/crypto_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[aes](https://pkg.go.dev/crypto/aes) | AES 加密 |
| &nbsp;&nbsp;&nbsp;&nbsp;[cipher](https://pkg.go.dev/crypto/cipher) | 标准分组加密 |
| &nbsp;&nbsp;&nbsp;&nbsp;[des](https://pkg.go.dev/crypto/des) | 数据加密标准 (DES) 和三重数据加密算法 (TDEA) |
| &nbsp;&nbsp;&nbsp;&nbsp;[dsa](https://pkg.go.dev/crypto/dsa) | DSA 数字签名算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[ecdh](https://pkg.go.dev/crypto/ecdh) |  NIST 曲线和 Curve25519 椭圆曲线 |
| &nbsp;&nbsp;&nbsp;&nbsp;[ecdsa](https://pkg.go.dev/crypto/ecdsa) | FIPS 186-5 椭圆曲线数字签名算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[ed25519](https://pkg.go.dev/crypto/ed25519) |  Ed25519 数字签名算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[elliptic](https://pkg.go.dev/crypto/elliptic) | NIST P-224、P-256、P-384 和 P-521 椭圆曲线 |
| &nbsp;&nbsp;&nbsp;&nbsp;[fips140](https://pkg.go.dev/crypto/fips140) | 本地程序是否启用 FIPS 140-3 模式 |
| &nbsp;&nbsp;&nbsp;&nbsp;[hkdf](https://pkg.go.dev/crypto/hkdf) | 基于 HMAC 的提取与展开密钥派生函数（HKDF） |
| &nbsp;&nbsp;&nbsp;&nbsp;[hmac](https://pkg.go.dev/crypto/hmac) | 密钥哈希消息认证码（HMAC） |
| &nbsp;&nbsp;&nbsp;&nbsp;[hpke](https://pkg.go.dev/crypto/hpke) | 混合公钥加密（HPKE） |
| &nbsp;&nbsp;&nbsp;&nbsp;[md5](https://pkg.go.dev/crypto/md5) | MD5 哈希算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[mlkem](https://pkg.go.dev/crypto/mlkem) | 实现 ML-KEM 量子抗性密钥封装方法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[mlkem/mlkemtest](https://pkg.go.dev/crypto/mlkem/mlkemtest) |为 ML-KEM 算法提供测试函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[pbkdf2](https://pkg.go.dev/crypto/pbkdf2) | 实现 PBKDF2 密钥派生函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rand](https://pkg.go.dev/crypto/rand) | 实现加密安全的随机数生成器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rc4](https://pkg.go.dev/crypto/rc4) | 实现 RC4 加密 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rsa](https://pkg.go.dev/crypto/rsa) | 实现 RSA 加密 |
| &nbsp;&nbsp;&nbsp;&nbsp;[sha1](https://pkg.go.dev/crypto/sha1) | 实现 SHA-1 哈希算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[sha256](https://pkg.go.dev/crypto/sha256) | 实现 SHA-224 和 SHA-256 哈希算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[sha3](https://pkg.go.dev/crypto/sha3) | 实现 SHA-3 哈希算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[sha512](https://pkg.go.dev/crypto/sha512) | 实现 SHA-384、SHA-512 等哈希算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[subtle](https://pkg.go.dev/crypto/subtle) | 实现加密代码中常用的函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[tls](https://pkg.go.dev/crypto/tls) | 部分实现 TLS 1.2 和 TLS 1.3 |
| &nbsp;&nbsp;&nbsp;&nbsp;[x509](https://pkg.go.dev/crypto/x509) | 实现 X.509 标准的子集 |
| &nbsp;&nbsp;&nbsp;&nbsp;[x509/pkix](https://pkg.go.dev/crypto/x509/pkix) | 用于 ASN.1 解析和序列化 X.509 证书、CRL 和 OCSP 的共享低层结构 |
| [database](https://pkg.go.dev/database) | 数据库相关包 | [[↗]](./Go%20Libs/GoSTD_examples/database/database_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[sql](https://pkg.go.dev/database/sql) | 为 SQL（或类 SQL）数据库提供通用接口 |
| &nbsp;&nbsp;&nbsp;&nbsp;[sql/driver](https://pkg.go.dev/database/sql/driver) | 定义数据库驱动程序要实现的接口 |
| [debug](https://pkg.go.dev/debug) | 调试相关包 | [[↗]](./Go%20Libs/GoSTD_examples/debug/debug_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[buildinfo](https://pkg.go.dev/debug/buildinfo) | 提供对嵌入在 Go 二进制文件中的构建信息的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[dwarf](https://pkg.go.dev/debug/dwarf) | 提供对从可执行文件加载的 DWARF 调试信息的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[elf](https://pkg.go.dev/debug/elf) | 实现对 ELF 对象文件的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[gosym](https://pkg.go.dev/debug/gosym) | 实现对 Go 二进制文件中嵌入的 Go 符号和行号表的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[macho](https://pkg.go.dev/debug/macho) | 实现对 Mach-O 对象文件的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[pe](https://pkg.go.dev/debug/pe) | 实现对 PE（Microsoft Windows 可移植可执行文件）文件的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[plan9obj](https://pkg.go.dev/debug/plan9obj) | 实现对 Plan 9 a.out 对象文件的访问 |
| [embed](https://pkg.go.dev/embed) | 提供对嵌入在运行中的 Go 程序中的文件的访问 | [[↗]](./Go%20Libs/GoSTD_examples/embed/embed_test.go) |
| [encoding](https://pkg.go.dev/encoding) | 定义其他包共享的接口，用于在字节级和文本表示之间转换数据 | [[↗]](./Go%20Libs/GoSTD_examples/encoding/encoding_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[ascii85](https://pkg.go.dev/encoding/ascii85) | 实现 ascii85 数据编码 |
| &nbsp;&nbsp;&nbsp;&nbsp;[asn1](https://pkg.go.dev/encoding/asn1) | 实现 DER 编码的 ASN.1 数据结构的解析 |
| &nbsp;&nbsp;&nbsp;&nbsp;[base32](https://pkg.go.dev/encoding/base32) | 实现 base32 编码 |
| &nbsp;&nbsp;&nbsp;&nbsp;[base64](https://pkg.go.dev/encoding/base64) | 实现 base64 编码 |
| &nbsp;&nbsp;&nbsp;&nbsp;[binary](https://pkg.go.dev/encoding/binary) | 实现数字和字节序列之间的简单转换 |
| &nbsp;&nbsp;&nbsp;&nbsp;[csv](https://pkg.go.dev/encoding/csv) | 读写逗号分隔值 (CSV) 文件 |
| &nbsp;&nbsp;&nbsp;&nbsp;[gob](https://pkg.go.dev/encoding/gob) | 管理 gob 流 - 在 Encoder（发送器）和 Decoder（接收器）之间交换的二进制值 |
| &nbsp;&nbsp;&nbsp;&nbsp;[hex](https://pkg.go.dev/encoding/hex) | 实现十六进制编码和解码 |
| &nbsp;&nbsp;&nbsp;&nbsp;[json](https://pkg.go.dev/encoding/json) | 实现 JSON 的编码和解码 |
| &nbsp;&nbsp;&nbsp;&nbsp;[json/jsontext](https://pkg.go.dev/encoding/json/jsontext) | 实现 JSON 的语法处理 |
| &nbsp;&nbsp;&nbsp;&nbsp;[json/v2](https://pkg.go.dev/encoding/json/v2) | 实现 JSON 的语义处理 |
| &nbsp;&nbsp;&nbsp;&nbsp;[pem](https://pkg.go.dev/encoding/pem) | 实现 PEM 数据编码 |
| &nbsp;&nbsp;&nbsp;&nbsp;[xml](https://pkg.go.dev/encoding/xml) | 实现一个简单的 XML 1.0 解析器 |
| [errors](https://pkg.go.dev/errors) | 实现操作错误的函数 | [[↗]](./Go%20Libs/GoSTD_examples/errors/errors_test.go) |
| [expvar](https://pkg.go.dev/expvar) | 为公共变量（如服务器中的操作计数器）提供标准化接口 | [[↗]](./Go%20Libs/GoSTD_examples/expvar/expvar_test.go) |
| [flag](https://pkg.go.dev/flag) | 实现命令行标志解析 | [[↗]](./Go%20Libs/GoSTD_examples/flag/flag_test.go) |
| [fmt](https://pkg.go.dev/fmt) | 实现格式化 I/O，具有类似于 C 的 printf 和 scanf 的函数 | [[↗]](./Go%20Libs/GoSTD_examples/fmt/fmt_test.go) |
| [go](https://pkg.go.dev/go) | Go 语言相关包 |
| &nbsp;&nbsp;&nbsp;&nbsp;[ast](https://pkg.go.dev/go/ast) | 声明用于表示 Go 包语法树的类型 |
| &nbsp;&nbsp;&nbsp;&nbsp;[build](https://pkg.go.dev/go/build) | 收集有关 Go 包的信息 |
| &nbsp;&nbsp;&nbsp;&nbsp;[build/constraint](https://pkg.go.dev/go/build/constraint) | 实现构建约束行的解析和评估 |
| &nbsp;&nbsp;&nbsp;&nbsp;[constant](https://pkg.go.dev/go/constant) | 实现表示无类型 Go 常量的值及其相应的操作 |
| &nbsp;&nbsp;&nbsp;&nbsp;[doc](https://pkg.go.dev/go/doc) | 从 Go AST 中提取源代码文档 |
| &nbsp;&nbsp;&nbsp;&nbsp;[doc/comment](https://pkg.go.dev/go/doc/comment) | 实现 Go 文档注释的解析和重新格式化 |
| &nbsp;&nbsp;&nbsp;&nbsp;[format](https://pkg.go.dev/go/format) | 实现 Go 源代码的标准格式化 |
| &nbsp;&nbsp;&nbsp;&nbsp;[importer](https://pkg.go.dev/go/importer) | 提供对导出数据导入器的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[parser](https://pkg.go.dev/go/parser) | 实现 Go 源文件的解析器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[printer](https://pkg.go.dev/go/printer) | 实现 AST 节点的打印 |
| &nbsp;&nbsp;&nbsp;&nbsp;[scanner](https://pkg.go.dev/go/scanner) | 实现 Go 源文本的扫描器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[token](https://pkg.go.dev/go/token) | 定义表示 Go 编程语言词法标记的常量 |
| &nbsp;&nbsp;&nbsp;&nbsp;[types](https://pkg.go.dev/go/types) | 声明数据类型并实现 Go 包类型检查的算法 |
| &nbsp;&nbsp;&nbsp;&nbsp;[version](https://pkg.go.dev/go/version) | 提供对 Go 版本的操作 |
| [hash](https://pkg.go.dev/hash) | 提供哈希函数的接口 | [[↗]](./Go%20Libs/GoSTD_examples/hash/hash_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[adler32](https://pkg.go.dev/hash/adler32) | 实现 Adler-32 校验和 |
| &nbsp;&nbsp;&nbsp;&nbsp;[crc32](https://pkg.go.dev/hash/crc32) | 实现 32 位循环冗余校验 (CRC-32) |
| &nbsp;&nbsp;&nbsp;&nbsp;[crc64](https://pkg.go.dev/hash/crc64) | 实现 64 位循环冗余校验 (CRC-64) |
| &nbsp;&nbsp;&nbsp;&nbsp;[fnv](https://pkg.go.dev/hash/fnv) | 实现 FNV-1 和 FNV-1a 非加密哈希函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[maphash](https://pkg.go.dev/hash/maphash) | 提供对字节序列和可比较值的哈希函数 |
| [html](https://pkg.go.dev/html) | 提供 HTML 文本的转义和反转义函数 | [[↗]](./Go%20Libs/GoSTD_examples/html/html_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[template](https://pkg.go.dev/html/template) | 实现数据驱动的模板，用于生成可防止代码注入的 HTML 输出 |
| [image](https://pkg.go.dev/image) | 实现基本的 2D 图像库 | [[↗]](./Go%20Libs/GoSTD_examples/image/image_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[color](https://pkg.go.dev/image/color) | 实现基本的颜色库 |
| &nbsp;&nbsp;&nbsp;&nbsp;[color/palette](https://pkg.go.dev/image/color/palette) | 提供标准颜色调色板 |
| &nbsp;&nbsp;&nbsp;&nbsp;[draw](https://pkg.go.dev/image/draw) | 提供图像合成函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[gif](https://pkg.go.dev/image/gif) | 实现 GIF 图像解码器和编码器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[jpeg](https://pkg.go.dev/image/jpeg) | 实现 JPEG 图像解码器和编码器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[png](https://pkg.go.dev/image/png) | 实现 PNG 图像解码器和编码器 |
| [index](https://pkg.go.dev/index) | 索引相关包 | [[↗]](./Go%20Libs/GoSTD_examples/index/index_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[suffixarray](https://pkg.go.dev/index/suffixarray) | 实现使用内存后缀数组的对数时间子串搜索 |
| [io](https://pkg.go.dev/io) | I/O 原语接口 | [[↗]](./Go%20Libs/GoSTD_examples/io/io_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[fs](https://pkg.go.dev/io/fs) | 定义文件系统接口 |
| &nbsp;&nbsp;&nbsp;&nbsp;[ioutil](https://pkg.go.dev/io/ioutil) | 提供 I/O 实用函数（已弃用，建议使用 os 和 io 包） |
| [log](https://pkg.go.dev/log) | 实现简单的日志记录包 | [[↗]](./Go%20Libs/GoSTD_examples/log/log_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[slog](https://pkg.go.dev/log/slog) | 结构化日志记录 |
| &nbsp;&nbsp;&nbsp;&nbsp;[syslog](https://pkg.go.dev/log/syslog) | 实现系统日志服务的接口 |
| [maps](https://pkg.go.dev/maps) | 提供对任意类型映射的操作函数 |
| [math](https://pkg.go.dev/math) | 数学函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[big](https://pkg.go.dev/math/big) | 实现大整数和有理数以任意精度 |
| &nbsp;&nbsp;&nbsp;&nbsp;[bits](https://pkg.go.dev/math/bits) | 提供位操作函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[cmplx](https://pkg.go.dev/math/cmplx) | 提供复数的基本常量和数学函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rand](https://pkg.go.dev/math/rand) | 实现伪随机数生成器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rand/v2](https://pkg.go.dev/math/rand/v2) | 实现伪随机数生成器（v2 版本） |
| [mime](https://pkg.go.dev/mime) | MIME 相关包 | [[↗]](./Go%20Libs/GoSTD_examples/mime/mime_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[multipart](https://pkg.go.dev/mime/multipart) | 实现 MIME 多部分消息的解析 |
| &nbsp;&nbsp;&nbsp;&nbsp;[quotedprintable](https://pkg.go.dev/mime/quotedprintable) | 实现 quoted-printable 编码，如 RFC 2045 中所述 |
| [net](https://pkg.go.dev/net) | 网络相关包 | [[↗]](./Go%20Libs/GoSTD_examples/net/net_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[http](https://pkg.go.dev/net/http) | 提供 HTTP 客户端和服务器实现 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/cgi](https://pkg.go.dev/net/http/cgi) | 实现 CGI 规范 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/cookiejar](https://pkg.go.dev/net/http/cookiejar) | 提供 HTTP 客户端的 Cookie 管理 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/fcgi](https://pkg.go.dev/net/http/fcgi) | 实现 FastCGI 规范 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/httptest](https://pkg.go.dev/net/http/httptest) | 提供 HTTP 测试工具 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/httptrace](https://pkg.go.dev/net/http/httptrace) | 提供 HTTP 客户端请求跟踪功能 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/httputil](https://pkg.go.dev/net/http/httputil) | 提供 HTTP 工具函数 |
| &nbsp;&nbsp;&nbsp;&nbsp;[http/pprof](https://pkg.go.dev/net/http/pprof) | 提供性能分析工具 |
| &nbsp;&nbsp;&nbsp;&nbsp;[netip](https://pkg.go.dev/net/netip) | 提供 IP 地址和端口的处理功能 |
| &nbsp;&nbsp;&nbsp;&nbsp;[mail](https://pkg.go.dev/net/mail) | 实现邮件地址解析 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rpc](https://pkg.go.dev/net/rpc) | 提供对通过网络或其他 I/O 连接导出的对象方法的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[rpc/jsonrpc](https://pkg.go.dev/net/rpc/jsonrpc) | 实现 JSON-RPC 协议 |
| &nbsp;&nbsp;&nbsp;&nbsp;[smtp](https://pkg.go.dev/net/smtp) | 实现 SMTP 客户端 |
| &nbsp;&nbsp;&nbsp;&nbsp;[textproto](https://pkg.go.dev/net/textproto) | 实现文本协议 |
| &nbsp;&nbsp;&nbsp;&nbsp;[url](https://pkg.go.dev/net/url) | 实现 URL 解析 |
| [os](https://pkg.go.dev/os) | 操作系统功能接口 | [[↗]](./Go%20Libs/GoSTD_examples/os/os_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[exec](https://pkg.go.dev/os/exec) | 运行外部命令 |
| &nbsp;&nbsp;&nbsp;&nbsp;[signal](https://pkg.go.dev/os/signal) | 实现对输入信号的访问 |
| &nbsp;&nbsp;&nbsp;&nbsp;[user](https://pkg.go.dev/os/user) | 提供对当前用户和用户组的访问 |
| [path](https://pkg.go.dev/path) | 路径操作 | [[↗]](./Go%20Libs/GoSTD_examples/path/path_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[filepath](https://pkg.go.dev/path/filepath) | 实现对操作系统路径的操作 |
| [plugin](https://pkg.go.dev/plugin) | 支持 Go 插件系统 |
| [reflect](https://pkg.go.dev/reflect) | 实现运行时反射，允许程序操作任意类型的对象 |
| [regexp](https://pkg.go.dev/regexp) | 正则表达式 |
| &nbsp;&nbsp;&nbsp;&nbsp;[syntax](https://pkg.go.dev/regexp/syntax) | 正则表达式语法分析 |
| [runtime](https://pkg.go.dev/runtime) | 运行时支持 | [[↗]](./Go%20Libs/GoSTD_examples/runtime/runtime_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[cgo](https://pkg.go.dev/runtime/cgo) | 支持 CGO 功能 |
| &nbsp;&nbsp;&nbsp;&nbsp;[coverage](https://pkg.go.dev/runtime/coverage) | 代码覆盖率工具 |
| &nbsp;&nbsp;&nbsp;&nbsp;[debug](https://pkg.go.dev/runtime/debug) | 提供调试功能 |
| &nbsp;&nbsp;&nbsp;&nbsp;[metrics](https://pkg.go.dev/runtime/metrics) | 提供运行时指标 |
| &nbsp;&nbsp;&nbsp;&nbsp;[pprof](https://pkg.go.dev/runtime/pprof) | 提供性能分析功能 |
| &nbsp;&nbsp;&nbsp;&nbsp;[race](https://pkg.go.dev/runtime/race) | 支持竞态检测 |
| &nbsp;&nbsp;&nbsp;&nbsp;[secret](https://pkg.go.dev/runtime/secret) | 密钥管理功能 |
| &nbsp;&nbsp;&nbsp;&nbsp;[trace](https://pkg.go.dev/runtime/trace) | 提供运行时追踪功能 |
| [simd](https://pkg.go.dev/simd) | 提供 SIMD 加速 |
| &nbsp;&nbsp;&nbsp;&nbsp;[archsimd](https://pkg.go.dev/simd/archsimd) | 提供架构特定的 SIMD 实现 |
| [slices](https://pkg.go.dev/slices) | 提供对切片的操作函数 |
| [sort](https://pkg.go.dev/sort) | 排序算法 | [[↗]](./Go%20Libs/GoSTD_examples/sort/sort_test.go) |
| [strconv](https://pkg.go.dev/strconv) | 字符串转换 | [[↗]](./Go%20Libs/GoSTD_examples/strconv/strconv_test.go) |
| [strings](https://pkg.go.dev/strings) | 字符串操作 | [[↗]](./Go%20Libs/GoSTD_examples/strings/strings_test.go) |
| [structs](https://pkg.go.dev/structs) | 提供对结构体的操作函数 |
| [sync](https://pkg.go.dev/sync) | 同步原语 | [[↗]](./Go%20Libs/GoSTD_examples/sync/sync_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[atomic](https://pkg.go.dev/sync/atomic) | 提供原子操作 |
| [syscall](https://pkg.go.dev/syscall) | 系统调用接口 |
| &nbsp;&nbsp;&nbsp;&nbsp;[js](https://pkg.go.dev/syscall/js) | 提供 JavaScript 接口 |
| [testing](https://pkg.go.dev/testing) | 提供 Go 自动化测试 | [[↗]](./Go%20Libs/GoSTD_examples/testing/testing_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[cryptotest](https://pkg.go.dev/testing/cryptotest) | 提供确定性随机源测试 |[[↗]](./Go%20Libs/GoSTD_examples/testing/cryptotest_test.go)
| &nbsp;&nbsp;&nbsp;&nbsp;[fstest](https://pkg.go.dev/testing/fstest) | 提供文件系统测试工具 | [[↗]](./Go%20Libs/GoSTD_examples/testing/fstest_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[iotest](https://pkg.go.dev/testing/iotest) | 提供 I/O 测试工具 | [[↗]](./Go%20Libs/GoSTD_examples/testing/iotest_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[quick](https://pkg.go.dev/testing/quick) | 提供快速测试功能 | [[↗]](./Go%20Libs/GoSTD_examples/testing/quick_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[slogtest](https://pkg.go.dev/testing/slogtest) | 提供对 log/slog 包的测试工具 | [[↗]](./Go%20Libs/GoSTD_examples/testing/slogtest_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[synctest](https://pkg.go.dev/testing/synctest) | 测试并发代码 | [[↗]](./Go%20Libs/GoSTD_examples/testing/synctest_test.go) |
| [text](https://pkg.go.dev/text) | 文本处理 |
| &nbsp;&nbsp;&nbsp;&nbsp;[scanner](https://pkg.go.dev/text/scanner) | 实现词法扫描器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[tabwriter](https://pkg.go.dev/text/tabwriter) | 实现写表格的过滤器 |
| &nbsp;&nbsp;&nbsp;&nbsp;[template](https://pkg.go.dev/text/template) | 实现数据驱动的模板 |
| &nbsp;&nbsp;&nbsp;&nbsp;[template/parse](https://pkg.go.dev/text/template/parse) | 实现模板解析 |
| [time](https://pkg.go.dev/time) | 时间和日期 | [[↗]](./Go%20Libs/GoSTD_examples/time/time_test.go) |
| &nbsp;&nbsp;&nbsp;&nbsp;[tzdata](https://pkg.go.dev/time/tzdata) | 提供时区数据 |
| [unicode](https://pkg.go.dev/unicode) | Unicode 支持 |
| &nbsp;&nbsp;&nbsp;&nbsp;[utf8](https://pkg.go.dev/unicode/utf8) | 实现 UTF-8 编码的支持 |
| &nbsp;&nbsp;&nbsp;&nbsp;[utf16](https://pkg.go.dev/unicode/utf16) | 实现 UTF-16 编码的支持 |
| [unique](https://pkg.go.dev/unique) | 提供唯一性保证 |
| [unsafe](https://pkg.go.dev/unsafe) | 提供绕过 Go 类型安全的操作 |
| [weak](https://pkg.go.dev/weak) | 提供弱引用功能 |
