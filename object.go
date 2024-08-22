package libjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

func indexByKey(data any, key any) (any, error) {
	switch v := data.(type) {
	case nil:
		return nil, errors.New("Can not index into null")
	case string:
		return nil, errors.New("Can not index into string")
	case float64:
		return nil, errors.New("Can not index into number")
	case []any:
		if len(v) == 0 {
			return nil, nil
		}
		if k, ok := key.(int); !ok {
			return nil, fmt.Errorf("Can not use %T::%v to index into %T::%v", key, key, data, data)
		} else {
			return v[k], nil
		}
	case map[string]any:
		if len(v) == 0 {
			return nil, nil
		}
		if k, ok := key.(string); !ok {
			return nil, fmt.Errorf("Can not use %T::%v to index into %T::%v", key, key, data, data)
		} else {
			return v[k], nil
		}
	default:
		return nil, fmt.Errorf("Unsupported %T, can not index", data)
	}
}

func parsePath(path string) (func(any) (any, error), error) {
	if len(path) == 0 {
		return nil, errors.New("Unexpected index syntax, top level element is available via '.'")
	}

	// fast paths for '.' path / parent access
	if len(path) == 1 && path[0] == '.' {
		return func(a any) (any, error) {
			return a, nil
		}, nil
	}

	// skip first . because we handled that above
	path = path[1:]

	keys := make([]any, 0, len(path)/4)
	lastIndex := 0
	for i, b := range path {
		if b == '.' {
			keys = append(keys, path[lastIndex:i])
			lastIndex = i + 1
		} else if i+1 == len(path) {
			keys = append(keys, path[lastIndex:i+1])
		}
	}

	return func(a any) (any, error) {
		val := a
		for _, k := range keys {
			key := k.(string)
			if key[0] >= '0' && key[0] <= '9' {
				if k1, err := strconv.ParseInt(key, 10, 32); err == nil {
					k = int(k1)
				}
			}

			if v, err := indexByKey(val, k); err != nil {
				return nil, err
			} else {
				val = v
			}
		}
		return val, nil
	}, nil
}

type JSON struct {
	obj any
}

func (j *JSON) get(path string) (any, error) {
	f, err := parsePath(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", errors.ErrUnsupported, path)
	}
	return f(j.obj)
}

func (j *JSON) set(path string, value any) error {
	// TODO:
	return nil
}

func (j *JSON) compile() (func() (any, error), error) {
	f := func() (any, error) { return nil, nil }
	// TODO:
	return f, errors.ErrUnsupported
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.obj)
}
