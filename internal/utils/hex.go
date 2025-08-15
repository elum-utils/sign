package utils

import (
	"encoding/hex"
	"errors"
	"unsafe"
)

// DecodeHexStringInto декодирует строку hex в готовый буфер dst.
// Возвращает количество записанных байт и ошибку (если есть).
func DecodeHexStringInto(src string, dst []byte) (int, error) {
	if len(src)%2 != 0 {
		return 0, errors.New("invalid hex string length")
	}
	if len(dst) < len(src)/2 {
		return 0, errors.New("destination buffer too small")
	}

	// hex.Decode работает с []byte, но []byte(src) создаст копию.
	// Чтобы избежать аллокации, декодируем напрямую через hex.Decode:
	return hex.Decode(dst, unsafeStringToBytes(src))
}

// unsafeStringToBytes — zero-copy преобразование string → []byte.
// Использовать только для read-only операций!
func unsafeStringToBytes(s string) []byte {
	sh := (*[2]uintptr)(unsafe.Pointer(&s))
	bh := [3]uintptr{sh[0], sh[1], sh[1]}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
