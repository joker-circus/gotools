package gotools

import (
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// 是否 UTF8 编码
func IsUTF8(s string) bool {
	return utf8.ValidString(s)
}

// 是否 GBK 编码
func IsGBK(s string) bool {
	if IsUTF8(s) {
		return false
	}
	data := []byte(s)
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			i++
			continue
		}

		if data[i] >= 0x81 &&
			data[i] <= 0xfe &&
			data[i+1] >= 0x40 &&
			data[i+1] <= 0xfe &&
			data[i+1] != 0xf7 {
			i += 2
			continue
		}

		return false
	}
	return true
}

// GBK 转 UTF8。对 GBK 字符解码。
//
// GB2312 使用单字节和双字节来进行编码，
// 单个字节编码与ASCII完全一样，兼容了ASCII码，
// 双字节编码高低字节范围都是0xA1-0xFE（范围0xA1A1 - 0xFEFE ）。
// 该范围的并不是每一个码位都编有字符，GB2312 只编码汉字6763个和非汉字图形字符682个。
// http://tools.jb51.net/table/gb2312
//
// GBK 编码同样使用单字节和双字节来进行编码，
// 单字节也同样采用ASCII的编码，
// 双字节GBK编码范围在0x8140 - 0xFEFE之间，覆盖了GB2312的编码范围。
// 收入 21886 个汉字和图形符号，其中汉字（包括部首和构件）21003 个，图形符号 883 个。
// http://tools.jb51.net/table/gbk_table
//
// GB18030 则在兼容 GBK 的基础上使用四个字节来达到扩展编码的目的。
//
// 等效于 simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(text))，
// 减少了 byte -> string、grow([]byte) 的损耗。
func DecodeGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, S2b(text), true)
	if err != nil {
		return text, err
	}

	return B2s(dst[:nDst]), nil
}

// 返回 UTF8 字符
func UTF8Text(text string) string {
	if !IsUTF8(text) {
		text, _ = DecodeGBK(text)
	}
	return text
}

// 返回 utf8 的 io.Reader
func UTF8Reader(r io.Reader) (io.Reader, error) {
	// 读取 Body，提前转换成 UTF-8
	body, err := io.ReadAll(r)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return strings.NewReader(UTF8Text(B2s(body))), nil
}
