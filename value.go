package libjson

import (
	"encoding/json"
	"unsafe"
)

type JSONKind uint8

const (
	JSONNull JSONKind = iota
	JSONBool
	JSONNumber
	JSONString
	JSONArray
	JSONObject
)

type JSONVal struct {
	Num float64
	Ptr uintptr
}

const (
	tagMask   = uintptr(0x7)
	tagNull   = uintptr(0)
	tagFalse  = uintptr(1)
	tagTrue   = uintptr(2)
	tagNumber = uintptr(3)
	tagString = uintptr(4)
	tagArray  = uintptr(5)
	tagObject = uintptr(6)
)

const arenaChunkSize = 256

type valueArena struct {
	stringChunks [][]stringNode
	arrayChunks  [][]arrayNode
	objectChunks [][]objectNode
}

type stringNode struct {
	value string
}

type arrayNode struct {
	value []JSONVal
}

type objectNode struct {
	value map[string]JSONVal
}

func (a *valueArena) allocStringNode(s string) uintptr {
	if len(a.stringChunks) == 0 || len(a.stringChunks[len(a.stringChunks)-1]) == cap(a.stringChunks[len(a.stringChunks)-1]) {
		a.stringChunks = append(a.stringChunks, make([]stringNode, 0, arenaChunkSize))
	}
	last := len(a.stringChunks) - 1
	idx := len(a.stringChunks[last])
	a.stringChunks[last] = append(a.stringChunks[last], stringNode{value: s})
	return uintptr(unsafe.Pointer(&a.stringChunks[last][idx]))
}

func (a *valueArena) allocArrayNode(v []JSONVal) uintptr {
	if len(a.arrayChunks) == 0 || len(a.arrayChunks[len(a.arrayChunks)-1]) == cap(a.arrayChunks[len(a.arrayChunks)-1]) {
		a.arrayChunks = append(a.arrayChunks, make([]arrayNode, 0, arenaChunkSize))
	}
	last := len(a.arrayChunks) - 1
	idx := len(a.arrayChunks[last])
	a.arrayChunks[last] = append(a.arrayChunks[last], arrayNode{value: v})
	return uintptr(unsafe.Pointer(&a.arrayChunks[last][idx]))
}

func (a *valueArena) allocObjectNode(v map[string]JSONVal) uintptr {
	if len(a.objectChunks) == 0 || len(a.objectChunks[len(a.objectChunks)-1]) == cap(a.objectChunks[len(a.objectChunks)-1]) {
		a.objectChunks = append(a.objectChunks, make([]objectNode, 0, arenaChunkSize))
	}
	last := len(a.objectChunks) - 1
	idx := len(a.objectChunks[last])
	a.objectChunks[last] = append(a.objectChunks[last], objectNode{value: v})
	return uintptr(unsafe.Pointer(&a.objectChunks[last][idx]))
}

func (a *valueArena) NewStringVal(s string) JSONVal {
	return JSONVal{Ptr: a.allocStringNode(s) | tagString}
}

func (a *valueArena) NewArrayVal(v []JSONVal) JSONVal {
	return JSONVal{Ptr: a.allocArrayNode(v) | tagArray}
}

func (a *valueArena) NewObjectVal(v map[string]JSONVal) JSONVal {
	return JSONVal{Ptr: a.allocObjectNode(v) | tagObject}
}

func NewNullVal() JSONVal {
	return JSONVal{Ptr: tagNull}
}

func NewBoolVal(v bool) JSONVal {
	if v {
		return JSONVal{Ptr: tagTrue}
	}
	return JSONVal{Ptr: tagFalse}
}

func NewNumberVal(v float64) JSONVal {
	return JSONVal{Num: v, Ptr: tagNumber}
}

func (v JSONVal) Kind() JSONKind {
	switch v.Ptr & tagMask {
	case tagNull:
		return JSONNull
	case tagFalse, tagTrue:
		return JSONBool
	case tagNumber:
		return JSONNumber
	case tagString:
		return JSONString
	case tagArray:
		return JSONArray
	case tagObject:
		return JSONObject
	default:
		panic("unreachable JSON tag")
	}
}

func (v JSONVal) Bool() bool {
	return (v.Ptr & tagMask) == tagTrue
}

func (v JSONVal) ptr() unsafe.Pointer {
	return unsafe.Pointer(v.Ptr &^ tagMask)
}

func (v JSONVal) String() string {
	return (*stringNode)(v.ptr()).value
}

func (v JSONVal) Array() []JSONVal {
	return (*arrayNode)(v.ptr()).value
}

func (v JSONVal) Object() map[string]JSONVal {
	return (*objectNode)(v.ptr()).value
}

func (v JSONVal) Interface() any {
	switch v.Kind() {
	case JSONNull:
		return nil
	case JSONBool:
		return v.Bool()
	case JSONNumber:
		return v.Num
	case JSONString:
		return v.String()
	case JSONArray:
		arr := v.Array()
		out := make([]any, len(arr))
		for i := range arr {
			out[i] = arr[i].Interface()
		}
		return out
	case JSONObject:
		obj := v.Object()
		out := make(map[string]any, len(obj))
		for k, child := range obj {
			out[k] = child.Interface()
		}
		return out
	default:
		panic("unreachable JSON kind")
	}
}

func (v JSONVal) MarshalJSON() ([]byte, error) {
	switch v.Kind() {
	case JSONNull:
		return []byte("null"), nil
	case JSONBool:
		if v.Bool() {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case JSONNumber:
		return json.Marshal(v.Num)
	case JSONString:
		return json.Marshal(v.String())
	case JSONArray:
		return json.Marshal(v.Array())
	case JSONObject:
		return json.Marshal(v.Object())
	default:
		panic("unreachable JSON kind")
	}
}
