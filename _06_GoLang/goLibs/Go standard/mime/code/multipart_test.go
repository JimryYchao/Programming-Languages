package gostd

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"
)

/*
! NewReader 使用给定的 MIME 边界创建一个新的 multipart Reader; 边界 boundary 通常是从消息的 “Content-Type” 标头的 “boundary” 参数获取的
! multipart.Reader 是 MIME 正文中各部分的迭代器
	NextPart, NextRawPart 返回多部分中的下一个部分或错误，末尾返回 io.EOF
	ReadForm 分析整个多部分消息
! NewWriter 返回一个具有随机边界的新 multipart Writer
	Boundary, Close,
	CreateFormField, CreateFormFile 使用给定的字段名称调用带有标头的 CreatePart。
	CreatePart 使用提供的 header 创建一个新的 multipart section
	FormDataContentType 返回具有此 Writer's Boundary 的 HTTP multipart/form-data 的 Content-Type
	SetBoundary 使用显式值重写 Writer 的默认随机生成的边界分隔符
		在创建任何 part 之前，必须调用 SetBoundary
	WriteField 调用 CreateFormField，然后写入给定值
! multipart.Writer
! File 是用于访问 multipart message 文件的接口
! FileHeader 描述 multipart request 的文件部分
	Open 将打开并返回 FileHeader 的关联文件
! Form 是解析的 multipart 表单
	RemoveAll 删除与表单关联的任何临时文件
! Part 表示 multipart 主体中的单个分部
	Close, Read 读取分部的正文，在 header 之后，下一个分部之前
	FileName 返回分部的 Content-Disposition 标头的 filename 参数
	FormName: 如果 p 的 Content-Disposition 类型为 form-data，则 FormName 返回 name 参数
*/

func TestRead(t *testing.T) {
	msg := &mail.Message{
		Header: map[string][]string{
			"Content-Type": {"multipart/mixed; boundary=foo"},
		},
		Body: strings.NewReader(
			"--foo\r\nFoo: one\r\n\r\nA section\r\n" +
				"--foo\r\nFoo: two\r\n\r\nAnd another\r\n" +
				"--foo--\r\n"),
	}
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			slurp, err := io.ReadAll(p)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Part %q: %q\n", p.Header.Get("Foo"), slurp)
		}
	}
}
