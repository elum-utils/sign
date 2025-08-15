package utils

import (
	"strings"
	"unsafe"
)

// QueryUnescape is a lightweight replacement for url.QueryUnescape.
//
// It decodes %XX sequences and replaces '+' with space.
// Uses an external buffer to avoid extra allocations.
// Returns false if the input contains invalid percent-encoding.
func QueryUnescape(s string, dstBuf *[]byte) (string, bool) {
	// Быстрая проверка — ничего не кодируем
	needsDecode := false
	for i := 0; i < len(s); i++ {
		if s[i] == '%' || s[i] == '+' {
			needsDecode = true
			break
		}
	}
	if !needsDecode {
		return s, true
	}

	b := (*dstBuf)[:0]
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			if i+2 >= len(s) {
				return "", false
			}
			hi := FromHex(s[i+1])
			lo := FromHex(s[i+2])
			if hi == 255 || lo == 255 {
				return "", false
			}
			b = append(b, hi<<4|lo)
			i += 3
		case '+':
			b = append(b, ' ')
			i++
		default:
			b = append(b, s[i])
			i++
		}
	}
	*dstBuf = b
	return unsafe.String(&b[0], len(b)), true
}

func QueryUnescapeBytes(s string, buf *[]byte) ([]byte, bool) {
	if strings.IndexByte(s, '%') == -1 && strings.IndexByte(s, '+') == -1 {
		return unsafeStringToBytes(s), true
	}

	b := (*buf)[:0]
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			if i+2 >= len(s) {
				return nil, false
			}
			hi := FromHex(s[i+1])
			lo := FromHex(s[i+2])
			if hi == 255 || lo == 255 {
				return nil, false
			}
			b = append(b, hi<<4|lo)
			i += 3
		case '+':
			b = append(b, ' ')
			i++
		default:
			b = append(b, s[i])
			i++
		}
	}
	*buf = b
	return b, true
}
