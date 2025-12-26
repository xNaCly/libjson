package libjson

// deserializer merges the concepts of lexical analysis with semantic and
// syntactical analysis, while producing direct go values out of the
// aforementioned. This results in a large performance improvement, since the
// deserialisation process no longer requires multiple passes and intermediate
// values.
type deserializer struct {
	src []byte
	pos uint
}

// State encodes all possible states the deserializer can be in. A transation
// is always defined as
//
//	Next(State, character) -> State
type state uint8

const (
	Start state = iota
	InObject
	InArray
	InString
	InNumber
	InLiteral
	Done
)

// converts any json from string to go value, result may be:
//
//	T = map[string]T, []T, string, float64, true, false, nil
func (d *deserializer) deserialize(src []byte) (any, error) {
	return nil, nil
}
