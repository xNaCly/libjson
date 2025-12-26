package libjson

import (
	"reflect"
	"testing"
)

func TestDeserialize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "empty object",
			input:    `{}`,
			expected: map[string]any{},
		},
		{
			name:     "empty array",
			input:    `[]`,
			expected: []any{},
		},
		{
			name:     "simple string",
			input:    `"abc"`,
			expected: "abc",
		},
		{
			name:     "int",
			input:    `1234`,
			expected: 1234.0,
		},
		{
			name:     "int single char",
			input:    `1`,
			expected: 1.0,
		},
		{
			name:     "double",
			input:    `3.1415`,
			expected: 3.1415,
		},
		{
			name:     "true",
			input:    `true`,
			expected: true,
		},
		{
			name:     "false",
			input:    `false`,
			expected: false,
		},
		{
			name:     "null",
			input:    `null`,
			expected: nil,
		},
		{
			name:  "array of literals",
			input: `["str", 3, 1.5, true, false, null]`,
			expected: []any{
				"str",
				3.0,
				1.5,
				true,
				false,
				nil,
			},
		},
		{
			name:  "simple object",
			input: `{"foo": "bar"}`,
			expected: map[string]any{
				"foo": "bar",
			},
		},
		{
			name:  "nested object",
			input: `{"a": {"b": [1,2,3]}}`,
			expected: map[string]any{
				"a": map[string]any{
					"b": []any{float64(1), float64(2), float64(3)},
				},
			},
		},

		// {
		// 	name:  "mixed array",
		// 	input: `[{"x":1}, "str", 3.14]`,
		// 	expected: []any{
		// 		map[string]any{"x": float64(1)},
		// 		"str",
		// 		3.14,
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deserialize([]byte(tt.input))
			if err != nil {
				t.Fatalf("deserialize() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("deserialize() = %v, want %v", got, tt.expected)
			}
		})
	}
}
