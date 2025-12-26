package libjson

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

// state encodes all possible states the deserializer can be in. A transation
// is always defined as
//
//	Next(state, character) -> state
type state uint8

// TODO: FUCK DFA's something is wrong and i have no clue what

const (
	Start state = iota
	InObject
	EndObject
	InArray
	EndArray
	String
	Number
	Atom
)

var transitions = map[state]map[byte]state{
	Start: {
		'{': InObject,
		'[': InArray,
		'"': String,
		'n': Atom,
		't': Atom,
		'f': Atom,
	},
	// no ':', since we always expect a colon between key and value
	// all rhs in InObject should be handled by Start transitions
	InObject: {
		// results in us popping a container and adding to its parent or
		// returning it if it is the only value
		'}': EndObject,
	},
	InArray: {
		// results in us popping a container and adding to its parent or
		// returning it if it is the only value
		']': EndArray,
	},
}

func init() {
	for b := byte('0'); b <= '9'; b++ {
		transitions[Start][b] = Number
	}
	transitions[Start]['-'] = Number
}

type container struct {
	isObj  bool   // true if InObject, for InArray false
	target any    // depending on container.t: map[string]any or []any
	key    string // current key if InObject
}

func insertValue(containerStack *[]container, v any) (any, bool) {
	if len(*containerStack) == 0 {
		return v, true
	}

	parent := &(*containerStack)[len(*containerStack)-1]
	if parent.isObj {
		parent.target.(map[string]any)[parent.key] = v
		parent.key = ""
	} else {
		parent.target = append(parent.target.([]any), v)
	}

	return nil, false
}

// converts any valid JSON to go values, result may be:
//
//	T = map[string]T, []T, string, float64, true, false, nil
//
// deserialize merges the concepts of lexical analysis with semantic and
// syntactical analysis, while producing direct go values out of the
// aforementioned. This results in a large performance improvements over the
// traditional approach, since the deserialisation process no longer requires
// multiple passes and intermediate values.
//
// At a high level this works by representing JSON as a table of states and
// possible input characters determining the follow state (table driven DFA).
// This makes deserialisation of large JSON inputs very fast.
func deserialize(src []byte) (any, error) {
	var pos int
	state := Start
	containerStack := make([]container, 0, 16)

	for pos < len(src) {
		for pos < len(src) && (src[pos] == ' ' || src[pos] == '\n' || src[pos] == '\t' || src[pos] == '\r') {
			pos++
		}
		if pos >= len(src) {
			break
		}

		b := src[pos]

		if b == ',' {
			pos++
			continue
		}

		next := transitions[state][b]

		fmt.Print("(", state.String(), ",", string(b), ")", "->", next.String(), "\n")

		switch next {
		case InArray:
			containerStack = append(containerStack, container{
				isObj:  false,
				target: []any{},
			})
			next = InArray
		case InObject:
			containerStack = append(containerStack, container{
				isObj:  true,
				target: map[string]any{},
			})
			next = InObject
		case EndArray, EndObject:
			last := containerStack[len(containerStack)-1]
			containerStack = containerStack[:len(containerStack)-1]

			if len(containerStack) == 0 {
				return last.target, nil
			}
			parent := &containerStack[len(containerStack)-1]
			if parent.isObj {
				parent.target.(map[string]any)[parent.key] = last.target
				parent.key = ""
			} else {
				parent.target = append(parent.target.([]any), last.target)
			}
			next = Start
		case String: // TODO: add support for escaping strings
			pos++ // skip "
			offset := bytes.IndexByte(src[pos:], '"')
			if offset < 0 {
				return nil, errors.New("Unterminated string")
			}
			end := pos + offset
			slice := src[pos:end]
			s := unsafe.String(unsafe.SliceData(slice), len(slice))
			pos = end

			if len(containerStack) == 0 {
				return s, nil
			}

			parent := &containerStack[len(containerStack)-1]
			if parent.isObj {
				if parent.key == "" {
					parent.key = s

					// since this is after the string at whose end we are
					if pos >= len(src) || src[pos+1] != ':' {
						return nil, fmt.Errorf("Expected ':' after object key")
					}
					pos++
				} else {
					parent.target.(map[string]any)[parent.key] = s
					parent.key = ""
				}
			} else {
				parent.target = append(parent.target.([]any), s)
			}

			next = Start
		case Number:
			start := pos
			for pos < len(src) && numChar[src[pos]] {
				pos++
			}

			if f, err := parseFloat(src[start:pos]); err != nil {
				return nil, err
			} else {
				if out, done := insertValue(&containerStack, f); done {
					return out, nil
				}
			}
			state = Start
			continue

		case Atom:
			var literal any
			switch src[pos] {
			case 't':
				if pos+3 > len(src) ||
					src[pos+1] != 'r' ||
					src[pos+2] != 'u' ||
					src[pos+3] != 'e' {
					return nil, errors.New("invalid true attempt")
				}
				pos += 4
				literal = true
			case 'f':
				if pos+4 > len(src) ||
					src[pos+1] != 'a' ||
					src[pos+2] != 'l' ||
					src[pos+3] != 's' ||
					src[pos+4] != 'e' {
					return nil, errors.New("invalid false attempt")
				}
				pos += 5
				literal = false
			case 'n':
				if pos+3 > len(src) ||
					src[pos+1] != 'u' ||
					src[pos+2] != 'l' ||
					src[pos+3] != 'l' {
					return nil, errors.New("invalid null attempt")
				}
				pos += 4
				literal = nil
			}

			if out, done := insertValue(&containerStack, literal); done {
				return out, nil
			}
			state = Start
			continue
		default:
			return nil, fmt.Errorf("No state transition found for (%s, %q)", state.String(), b)
		}
		state = next
		pos++
	}

	if len(containerStack) != 0 {
		return nil, fmt.Errorf("Unexpected end of input, unclosed container: %+v", containerStack)
	}

	return nil, nil
}
