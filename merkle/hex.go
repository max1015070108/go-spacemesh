package merkle

// fromHexChar converts a hex character into its value and a success flag.
// Adapted from https://golang.org/src/encoding/hex/hex.go - too bad it is private
func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

// Returns the common prefix of 2 hex encoded strings
// Empty string is return if there's no common suffix of len >= 1
func commonPrefix(s string, s1 string) string {
	l := lenPrefix(s, s1)
	return s[:l]
}

// Returns the length of the common prefix of 2 hex encoded strings
func lenPrefix(a, b string) int {

	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}
