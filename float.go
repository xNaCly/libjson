package libjson

import (
	"errors"
)

func pow10(exp int) float64 {
	res := 1.0
	if exp > 0 {
		for i := 0; i < exp; i++ {
			res *= 10
		}
	} else {
		for i := 0; i < -exp; i++ {
			res /= 10
		}
	}
	return res
}

// non allocating float parsing
func parseFloat(input []byte) (float64, error) {
	if len(input) == 0 {
		return 0, errors.New("empty input")
	}

	pos := 0
	neg := false
	if input[pos] == '-' {
		neg = true
		pos++
	}

	mantissa := uint64(0)
	exponent := 0
	seenDot := false

	for pos < len(input) {
		c := input[pos]
		if c >= '0' && c <= '9' {
			mantissa = mantissa*10 + uint64(c-'0')
			if seenDot {
				exponent--
			}
			pos++
		} else if c == '.' {
			if seenDot {
				return 0, errors.New("multiple dots in number")
			}
			seenDot = true
			pos++
		} else {
			break
		}
	}

	// weird eE+- handling
	if pos < len(input) && (input[pos] == 'e' || input[pos] == 'E') {
		pos++
		expNeg := false
		if pos < len(input) && input[pos] == '-' {
			expNeg = true
			pos++
		} else if pos < len(input) && input[pos] == '+' {
			pos++
		}

		if pos >= len(input) || input[pos] < '0' || input[pos] > '9' {
			return 0, errors.New("missing digits in exponent")
		}

		expVal := 0
		for pos < len(input) && input[pos] >= '0' && input[pos] <= '9' {
			expVal = expVal*10 + int(input[pos]-'0')
			pos++
		}
		if expNeg {
			expVal = -expVal
		}
		exponent += expVal
	}

	if mantissa == 0 {
		return 0, nil
	}

	result := float64(mantissa) * pow10(exponent)
	if neg {
		result = -result
	}
	return result, nil
}
