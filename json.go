package libjson

import (
	"io"
	"os"
)

func NewReader(r io.Reader) (JSON, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return JSON{}, err
	}
	return newJSON(data, nil)
}

// New consumes data and deserializes it into a JSON object
//
// data is consumed and possibly mutated, DO NOT REUSE
func New(data []byte) (JSON, error) {
	return newJSON(data, nil)
}

func newJSON(data []byte, cleanup func() error) (JSON, error) {
	p := parser{l: lexer{data: data, len: len(data)}}
	obj, err := p.parse(data)
	if err != nil {
		return JSON{}, err
	}
	return JSON{obj: obj, arena: p.arena, cleanup: cleanup}, nil
}

// FromFile is the same as New but zero copy via mmap
//
// The returned JSON value owns the memory mapping; call Close once the parsed
// data and any zero-copy strings derived from it are no longer needed.
func FromFile(f *os.File) (JSON, error) {
	data, cleanup, err := mmapFile(f)
	if err != nil {
		return JSON{}, err
	}
	obj, err := newJSON(data, cleanup)
	if err != nil {
		if cleanup != nil {
			_ = cleanup()
		}
		return JSON{}, err
	}
	return obj, nil
}
