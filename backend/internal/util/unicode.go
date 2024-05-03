package util

import (
	"bytes"
	"unicode/utf8"
)

func ValidateUTF8(str string) bool {
	return utf8.ValidString(str)
}

func CleanString(in string) string {
	var buf bytes.Buffer
	for i, r := range in {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(in[i:])
			if size == 1 {
				continue // 忽略无效的 utf8 字符
			}
		}
		buf.WriteRune(r)
	}
	return buf.String()
}
