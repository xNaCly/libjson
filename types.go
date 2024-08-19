package libjson

// json type
type t_json uint32

var empty = token{Type: t_eof, Val: nil}

const (
	t_string       t_json = iota // anything between ""
	t_number                     // floating point, hex, etc
	t_true                       // true
	t_false                      // false
	t_null                       // null
	t_left_curly                 // {
	t_right_curly                // }
	t_left_braket                // [
	t_right_braket               // ]
	t_comma                      // ,
	t_colon                      // :
	t_eof                        // for any non structure characters outside of strings and numbers
)

var tokennames = map[t_json]string{
	t_string:       "string",
	t_number:       "number",
	t_true:         "true",
	t_false:        "false",
	t_null:         "null",
	t_left_curly:   "{",
	t_right_curly:  "}",
	t_left_braket:  "[",
	t_right_braket: "]",
	t_comma:        ",",
	t_colon:        ":",
	t_eof:          "EOF",
}

type token struct {
	Type t_json
	Val  []byte // only populated for number and string
}
