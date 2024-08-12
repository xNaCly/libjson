package libjson

// json type
type t_json uint8

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
	t_unkown                     // for any non structure characters outside of strings and numbers
)

type token struct {
	Type t_json
	Val  any // only populated for number and string (can be number or string, depends on .Type)
}

type Node interface {
	Type() t_json
	Compute() any
}
