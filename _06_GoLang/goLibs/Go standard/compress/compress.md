<a id="TOP"></a>

## compress

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<!-- <a href="./code/XXX_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a> -->
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/compress"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `bzip2` 实现 bzip2 解压缩。[[↗]](code/bzip2_test.go)

包 `flate` 实现了 [RFC 1951](https://rfc-editor.org/rfc/rfc1951.html) 中描述的 DEFLATE 压缩数据格式。`gzip` 和 `zlib` 包实现了对基于 DEFLATE 格式文件的访问。[[↗]](code/flate_test.go)

包 `gzip` 实现了对 `gzip` 格式压缩文件的读取和写入。[[↗]](code/gzip_test.go)

包 `zlib` 实现了对 `zlib` 格式压缩数据的读取和写入。[[↗]](code/zlib_test.go)

包 `lzw` 实现 Lempel-Ziv-Welch 压缩数据格式，它实现了 GIF 和 PDF 文件格式所使用的 LZW。TIFF 文件格式使用类似但不兼容的 LZW 算法版本，参照 golang.org/x/image/tiff/lzw 包。[[↗]](code/lzw_test.go)


---
<a id="exam" ><a>

### Examples

- [`flate` 模拟网络传输压缩数据。](examples/flate_netconnect.go)

- [`gzip` 网络压缩传输](examples/gzip_httpSend.go)
---