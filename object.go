package libjson

import (
	"encoding/json"
	"fmt"
)

func Set[T any](obj *JSON, path string, value T) error {
	return obj.set(path, value)
}

func Get[T any](obj *JSON, path string) (T, error) {
	val, err := obj.get(path)
	if err != nil {
		var e T
		return e, err
	}
	if castVal, ok := val.(T); !ok {
		var e T
		return e, fmt.Errorf("Expected value of type %T, got type %T", e, val)
	} else {
		return castVal, nil
	}
}

func Compile[T any](obj *JSON, path string) (func() (T, error), error) {
	f, err := obj.compile()
	if err != nil {
		return nil, err
	}
	return func() (T, error) {
		val, err := f()
		if err != nil {
			var e T
			return e, err
		}
		if castVal, ok := val.(T); !ok {
			var e T
			return e, fmt.Errorf("Expected value of type %T, got type %T", e, val)
		} else {
			return castVal, nil
		}
	}, nil
}

type JSON struct {
	obj map[string]any
}

func (j *JSON) get(path string) (any, error) {
	return nil, nil
}

func (j *JSON) set(path string, value any) error {
	return nil
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.obj)
}

func (j *JSON) compile() (func() (any, error), error) {
	f := func() (any, error) { return nil, nil }
	return f, nil
}
