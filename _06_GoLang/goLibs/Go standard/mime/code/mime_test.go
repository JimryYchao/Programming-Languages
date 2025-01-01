package gostd

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"testing"
)

/*
! AddExtensionType 将与扩展名 ext 关联的 MIME 类型设置为 typ
! ExtensionsByType 返回已知与 MIME 类型类型关联的扩展
! TypeByExtension 返回与文件扩展名 ext 关联的 MIME 类型; 在 Windows 上，MIME 类型是从注册表中提取的。
	默认情况下，文本类型的 charset 参数设置为“utf-8”
! FormatMediaType 将 media type t 和参数 param 序列化为符合 RFC 2045 和 RFC 2616 的媒体类型
! ParseMediaType 根据 RFC 1521 分析 media type 值和任何可选参数 param
*/

func TestExtByType(t *testing.T) {
	mime.AddExtensionType(".xswl", "image/bmp") // 设置关联
	mtyps := []string{
		"image/gif",
		"image/jpeg",
		"image/png",
		"image/tiff",
		"image/bmp",
		"image/svg+xml",
	}
	for _, typ := range mtyps {
		texts, err := mime.ExtensionsByType(typ)
		if err != nil {
			t.Fatal(err)
		}
		logfln(">>> %s <<<", typ)
		for _, t := range texts {
			log(t)
		}
	}

}
func TestFormat_Parse_MIME(t *testing.T) {
	mediatype := "text/html"
	params := map[string]string{
		"charset": "utf-8",
	}
	result := mime.FormatMediaType(mediatype, params)
	log(result)

	mtye, params, _ := mime.ParseMediaType(result)
	log(" type:", mtye, "\n", "charset:", params["charset"])
}
func TestTypeByExt(t *testing.T) {
	exts := []string{
		"doc", "docx", "xls", "xlsx", "ppt", "pptx", "gz", "gzip", "zip", "7zip", "rar", "tar", "taz",
		"pdf", "rtf", "gif", "jpg", "jpeg", "jpg2", "png", "tif", "tiff", "bmp", "svg", "svgz", "webp", "ico",
		"wps", "et", "dps", "psd", "cdr", "swf", "txt", "js", "css", "htm", "html", "shtml", "xht", "xhtml",
		"vcf", "php", "php3", "php4", "phtml", "jar", "apk", "exe", "crt", "pem", "mp3", "mid", "midi", "wav", "m3u",
		"m4a", "ogg", "ra", "mp4", "mpg", "mpe", "mpeg", "qt", "mov", "m4v", "wmv", "avi", "webm", "flv", "evg", "fif",
		"spl", "hta", "acx", "hqx", "dot", "*", "bin", "class", "dms", "exe", "lha", "lzh", "oda",
		"axs", //...
	}
	for _, ext := range exts {
		logfln("%10s >> %s ", "*."+ext, mime.TypeByExtension("."+ext))
	}
}

/*
! mime.WordDecoder 解码包含 RFC 2047 编码字的 MIME 标头
	Decode 解码 RFC 2047 编码字
	DecodeHeader 解码给定字符串的所有编码字
! mime.WordEncoder 是 RFC 2047 编码的字编码器
	Encode 返回 s 的编码字形式。如果 s 是不带特殊字符的 ASCII，则返回原封不动。提供的字符集是 s 的 IANA 字符集名称。它不区分大小写。
*/

func TestWordDecode(t *testing.T) {
	dec := new(mime.WordDecoder)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "x-case":
			content, err := io.ReadAll(input)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(bytes.ToUpper(content)), nil
		default:
			return nil, fmt.Errorf("unhandled charset %q", charset)
		}
	}

	header, err := dec.Decode("=?utf-8?q?=C2=A1Hola,_se=C3=B1or!?=")
	if err != nil {
		panic(err)
	}
	log(header)

	header, _ = dec.Decode("=?x-case?q?hello!?=")
	log(header)

	//? DecodeHeader
	header, _ = dec.DecodeHeader("=?utf-8?q?=C3=89ric?= <eric@example.org>, =?utf-8?q?Ana=C3=AFs?= <anais@example.org>")
	log(header)

	header, _ = dec.DecodeHeader("=?utf-8?q?=C2=A1Hola,?= =?utf-8?q?_se=C3=B1or!?=")
	log(header)

	header, _ = dec.DecodeHeader("=?x-case?q?hello_?= =?x-case?q?world!?=")
	log(header)
}

func TestWordEncode(t *testing.T) {
	words := []struct{ charset, mess string }{
		{"utf-8", "¡Hola, señor!"},
		{"utf-8", "Hello!"},
		{"UTF-8", "¡Hola, señor!"},
		{"ISO-8859-1", "Caf\xE9"},
		{"utf-16", "大家好"},
	}

	for _, w := range words {
		logfln("%s QEncode : %s", w.mess, mime.QEncoding.Encode(w.charset, w.mess))
		logfln("%s BEncode : %s", w.mess, mime.BEncoding.Encode(w.charset, w.mess))
	}
}
