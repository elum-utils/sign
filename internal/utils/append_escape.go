package utils

// appendEscape performs percent-encoding (RFC 3986) on a string.
// Writes directly into dst without allocating intermediate strings.
func AppendEscape(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			dst = append(dst, c)
		} else {
			dst = append(dst, '%')
			dst = append(dst, "0123456789ABCDEF"[c>>4])
			dst = append(dst, "0123456789ABCDEF"[c&15])
		}
	}
	return dst
}
