package encode

import (
	"reflect"
	"testing"
)

func TestRLPEncoding(t *testing.T) {
	tests := map[string]struct {
		input    any
		expected []byte
	}{
		"string dog": {
			input:    "dog",
			expected: []byte{0x83, 'd', 'o', 'g'},
		},
		"list [cat, dog]": {
			input:    []string{"cat", "dog"},
			expected: []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'},
		},
		"bytes": {
			input:    []byte("dog"),
			expected: []byte{0x83, 'd', 'o', 'g'},
		},
		"empty string": {
			input:    "",
			expected: []byte{0x80},
		},
		"empty list": {
			input:    []any{},
			expected: []byte{0xc0},
		},
		"integer 0": {
			input:    0,
			expected: []byte{0x80},
		},
		"integer 15": {
			input:    15,
			expected: []byte{0x0f},
		},
		"integer 1024": {
			input:    1024,
			expected: []byte{0x82, 0x04, 0x00},
		},
		"byte 0x00": {
			input:    []byte{0x00},
			expected: []byte{0x00},
		},
		"byte 0x0f": {
			input:    []byte{0x0f},
			expected: []byte{0x0f},
		},
		"bytes 0x04 0x00": {
			input:    []byte{0x04, 0x00},
			expected: []byte{0x82, 0x04, 0x00},
		},
		"set theoretical representation [ [], [[]], [ [], [[]] ] ]": {
			input: []any{
				[]any{},
				[]any{[]any{}},
				[]any{
					[]any{},
					[]any{[]any{}},
				},
			},
			expected: []byte{0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0},
		},
		"long string Lorem ipsum...": {
			input:    "Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			expected: append([]byte{0xb8, 0x38}, []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit")...),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := RlpEncode(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Encode(%v) = %v, want %v", tt.input, result, tt.expected)
				t.Errorf("Encode(%v) = %x, want %x", tt.input, result, tt.expected)
			}
		})
	}
}
