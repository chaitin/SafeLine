package fsl

import (
	"fmt"
)

// see : https://chaitin.cn/patronus/fvm/-/blob/master/src/util/StrUtil.cpp#L18

func appendEscapedByte(buf []byte, b byte) []byte {
	if b >= 0x20 && b <= 0x7e && b != '\'' && b != '"' && b != '\\' {
		return append(buf, b)
	}
	switch b {
	case '\'':
		return append(buf, `\'`...)
	case '"':
		return append(buf, `\"`...)
	case '\n':
		return append(buf, `\n`...)
	case '\t':
		return append(buf, `\t`...)
	case '\r':
		return append(buf, `\r`...)
	case '\b':
		return append(buf, `\b`...)
	case '\f':
		return append(buf, `\f`...)
	case '\\':
		return append(buf, `\\`...)
	default:
		return append(buf, fmt.Sprintf("\\x%x", b)...)
	}
}

func appendQuotedWith(buf []byte, s string, quote byte) []byte {
	buf = append(buf, quote)
	for i := 0; i < len(s); i++ {
		buf = appendEscapedByte(buf, s[i])
	}
	buf = append(buf, quote)
	return buf
}

func quoteWith(s string, quote byte) string {
	return string(appendQuotedWith(make([]byte, 0, 3*len(s)/2), s, quote))
}

func Quote(s string) string {
	return quoteWith(s, '\'')
}
