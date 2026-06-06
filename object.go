package libjson

import (
	"errors"
	"fmt"
	"strconv"
)

type JSON struct {
	obj     JSONVal
	arena   valueArena
	cleanup func() error
}

func Get[T any](obj *JSON, path string) (T, error) {
	val, err := obj.get(path)
	if err != nil {
		var e T
		return e, err
	}
	if castVal, ok := val.Interface().(T); !ok {
		var e T
		return e, fmt.Errorf("Expected value of type %T, got json kind %v", e, val.Kind())
	} else {
		return castVal, nil
	}
}

func indexByKey(data JSONVal, key any) (JSONVal, error) {
	switch data.Kind() {
	case JSONNull:
		return JSONVal{}, errors.New("Can not index into null")
	case JSONString:
		return JSONVal{}, errors.New("Can not index into string")
	case JSONNumber:
		return JSONVal{}, errors.New("Can not index into number")
	case JSONBool:
		return JSONVal{}, errors.New("Can not index into bool")
	case JSONArray:
		arr := data.Array()
		if len(arr) == 0 {
			return JSONVal{}, nil
		}
		if k, ok := key.(int); !ok {
			return JSONVal{}, fmt.Errorf("Can not use %T::%v to index into json array", key, key)
		} else {
			return arr[k], nil
		}
	case JSONObject:
		obj := data.Object()
		if len(obj) == 0 {
			return JSONVal{}, nil
		}
		if k, ok := key.(string); !ok {
			return JSONVal{}, fmt.Errorf("Can not use %T::%v to index into json object", key, key)
		} else {
			return obj[k], nil
		}
	default:
		return JSONVal{}, fmt.Errorf("Unsupported json kind %v, can not index", data.Kind())
	}
}

func parsePath(path string) (func(JSONVal) (JSONVal, error), error) {
	if len(path) == 0 {
		return nil, errors.New("Unexpected index syntax, top level element is available via '.'")
	}

	// fast paths for '.' path / parent access
	if len(path) == 1 && path[0] == '.' {
		return func(a JSONVal) (JSONVal, error) {
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

	return func(a JSONVal) (JSONVal, error) {
		val := a
		for _, k := range keys {
			key := k.(string)
			if key[0] >= '0' && key[0] <= '9' {
				if k1, err := strconv.ParseInt(key, 10, 32); err == nil {
					k = int(k1)
				}
			}

			if v, err := indexByKey(val, k); err != nil {
				return JSONVal{}, err
			} else {
				val = v
			}
		}
		return val, nil
	}, nil
}

func (j *JSON) get(path string) (JSONVal, error) {
	f, err := parsePath(path)
	if err != nil {
		return JSONVal{}, fmt.Errorf("%w: %q", errors.ErrUnsupported, path)
	}
	return f(j.obj)
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	return j.obj.MarshalJSON()
}

func (j *JSON) Close() error {
	if j.cleanup == nil {
		return nil
	}
	cleanup := j.cleanup
	j.cleanup = nil
	return cleanup()
}
