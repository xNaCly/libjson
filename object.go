package libjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type JSON struct {
	inner any
}

// takes a JSON.inner value and converts it to Go, for instance merges the obj
// fields into a map
func toGo(json any) any {
	switch v := json.(type) {
	case obj:
		m := make(map[string]any, len(v.Fields))
		for _, f := range v.Fields {
			m[f.Key] = toGo(f.Value)
		}
		return m
	case []any:
		arr := make([]any, len(v))
		for i, el := range v {
			arr[i] = toGo(el)
		}
		return arr
	default:
		return v
	}
}

func Get[T any](obj *JSON, path string) (T, error) {
	val, err := obj.get(path)
	if err != nil {
		var e T
		return e, err
	}

	// normalise inner json representation into something Go can deal with
	val = toGo(val)

	if val == nil {
		var e T
		return e, nil
	}

	if castVal, ok := val.(T); !ok {
		var e T
		return e, fmt.Errorf("Expected value of type %T, got type %T", e, val)
	} else {
		return castVal, nil
	}
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
	case obj:
		if len(v.Fields) == 0 {
			return nil, nil
		}

		if k, ok := key.(string); !ok {
			return nil, fmt.Errorf("Can not use %T::%v to index into %T::%v", key, key, data, data)
		} else {
			i := 0
			for ; i < len(v.Fields); i++ {
				cur := v.Fields[i]
				if cur.Key == k {
					return cur.Value, nil
				}
			}
			return nil, nil
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

func (j *JSON) get(path string) (any, error) {
	f, err := parsePath(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", errors.ErrUnsupported, path)
	}
	return f(j.inner)
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(toGo(j.inner))
}
