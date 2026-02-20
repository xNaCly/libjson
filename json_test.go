package libjson

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const amount = 50_000
const naiveInput = `{"key1":"value","array":[],"obj":{},"atomArray":[11201,1e112,true,false,null,"str"]},`
const escapedInput = `{"text":"line1\nline2\nline3","quote":"\"hello\"","path":"C:\\\\Users\\\\name","unicode":"\u0041\u0042\u0043","mix":"abc\\ndef\"ghi\u263A"},`
const hardInput = `{
	"id":12345,
	"name":"very_long_string_with_no_escapes_but_large_payload_abcdefghijklmnopqrstuvwxyz_0123456789",
	"description":"This string contains\nmultiple\nlines\nand \"quotes\" and unicode \u2764\u2764\u2764",
	"nested":{
		"level1":{
			"level2":{
				"array":[
					"short",
					"string_with_escape\\n",
					"another\\tvalue",
					"unicode\u2603",
					1234567890,
					-1.2345e67,
					true,
					false,
					null
				]
			}
		}
	}
},`

func benchmarkWithInput(b *testing.B, input string) {
	data := strings.Repeat(input, amount)
	d := []byte("[" + data[:len(data)-1] + "]")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, len(d))
		copy(buf, d)
		b.StartTimer()
		_, err := New(buf)
		b.StopTimer()
		assert.NoError(b, err)
	}
	b.ReportAllocs()
}

func benchmarkEncodingJsonWithInput(b *testing.B, input string) {
	data := strings.Repeat(input, amount)
	d := []byte("[" + data[:len(data)-1] + "]")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v any
		err := json.Unmarshal(d, &v)
		assert.NoError(b, err)
	}
	b.ReportAllocs()
}

func BenchmarkLibJson_Naive(b *testing.B) {
	benchmarkWithInput(b, naiveInput)
}

func BenchmarkLibJson_Escaped(b *testing.B) {
	benchmarkWithInput(b, escapedInput)
}

func BenchmarkLibJson_Hard(b *testing.B) {
	benchmarkWithInput(b, hardInput)
}

func BenchmarkEncodingJson_Naive(b *testing.B) {
	benchmarkEncodingJsonWithInput(b, naiveInput)
}

func BenchmarkEncodingJson_Escaped(b *testing.B) {
	benchmarkEncodingJsonWithInput(b, escapedInput)
}

func BenchmarkEncodingJson_Hard(b *testing.B) {
	benchmarkEncodingJsonWithInput(b, hardInput)
}
