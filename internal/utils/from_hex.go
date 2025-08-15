package utils

// fromHex converts an ASCII hex digit into its value (0–15).
// Returns 255 if the character is not a valid hex digit.
func FromHex(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	default:
		return 255
	}
}
