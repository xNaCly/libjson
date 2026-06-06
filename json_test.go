package libjson

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const amount = 50_000
const naiveInput = `{"key1":"value","array":[],"obj":{},"atomArray":[11201,1e112,true,false,null,"str"]},`
const escapedInput = `{"text":"line1\nline2\nline3","quote":"\"hello\"","path":"C:\\\\Users\\\\name","unicode":"\u0041\u0042\u0043","mix":"abc\\ndef\"ghi\u263A"},`
const hardInput = `{
  "id": 12345,
  "name": "very_long_string_with_escapes_and_unicode_abcdefghijklmnopqrstuvwxyz_0123456789",
  "description": "This string contains\nmultiple\nlines\nand \"quotes\" and unicode \u2764\u2764\u2764",
  "nested": {
    "level1": {
      "level2": {
        "level3": {
          "level4": {
            "array": [
              "short",
              "string_with_escape\\n",
              "another\\tvalue",
              "unicode\u2603",
              "escaped_quote_\"_and_backslash_\\",
              1234567890,
              -1.2345e67,
              3.141592653589793,
              true,
              false,
              null,
              "ABC\u00a9\u20ac",
              "mix\\n\\t\\r\\\\\\\"end"
            ]
          }
        }
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

func benchmarkNewReaderWithInput(b *testing.B, input string) {
	data := strings.Repeat(input, amount)
	d := []byte("[" + data[:len(data)-1] + "]")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewReader(bytes.NewReader(d))
		assert.NoError(b, err)
	}
	b.ReportAllocs()
}

func benchmarkFromFileWithInput(b *testing.B, input string) {
	data := strings.Repeat(input, amount)
	d := []byte("[" + data[:len(data)-1] + "]")

	f, err := os.CreateTemp(b.TempDir(), "libjson-bench-*.json")
	assert.NoError(b, err)
	_, err = f.Write(d)
	assert.NoError(b, err)
	_, err = f.Seek(0, 0)
	assert.NoError(b, err)
	b.Cleanup(func() {
		assert.NoError(b, f.Close())
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Seek(0, 0)
		assert.NoError(b, err)
		doc, err := FromFile(f)
		assert.NoError(b, err)
		assert.NoError(b, doc.Close())
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

func BenchmarkLibJson_NewReader_Hard(b *testing.B) {
	benchmarkNewReaderWithInput(b, hardInput)
}

func BenchmarkLibJson_FromFile_Hard(b *testing.B) {
	benchmarkFromFileWithInput(b, hardInput)
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
