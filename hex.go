package libjson

import "errors"

var invalid_hex_err = errors.New("invalid hex")

var hexTable [256]byte

func init() {
	for i := 0; i < 256; i++ {
		hexTable[i] = 0xFF
	}
	for i := byte('0'); i <= '9'; i++ {
		hexTable[i] = i - '0'
	}
	for i := byte('a'); i <= 'f'; i++ {
		hexTable[i] = i - 'a' + 10
	}
	for i := byte('A'); i <= 'F'; i++ {
		hexTable[i] = i - 'A' + 10
	}
}

// hex4 converts 4 ASCII hex bytes to a rune.
// Returns an error if any byte is invalid.
func hex4(b []byte) (r rune, err error) {
	var v byte

	v = hexTable[b[0]]
	if v == 0xFF {
		return 0, invalid_hex_err
	}
	r = rune(v) << 12

	v = hexTable[b[1]]
	if v == 0xFF {
		return 0, invalid_hex_err
	}
	r |= rune(v) << 8

	v = hexTable[b[2]]
	if v == 0xFF {
		return 0, invalid_hex_err
	}
	r |= rune(v) << 4

	v = hexTable[b[3]]
	if v == 0xFF {
		return 0, invalid_hex_err
	}
	r |= rune(v)

	return r, nil
}
