package decode

import (
	"reflect"
	"testing"
)

func TestRLPDecoding(t *testing.T) {
	tests := map[string]struct {
		input    []byte
		expected any
	}{
		"string dog": {
			input:    []byte{0x83, 'd', 'o', 'g'},
			expected: []byte("dog"),
		},
		"list [cat, dog]": {
			input: []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'},
			expected: []any{
				[]byte("cat"),
				[]byte("dog"),
			},
		},
		"bytes": {
			input:    []byte{0x83, 'd', 'o', 'g'},
			expected: []byte("dog"),
		},
		"empty string": {
			input:    []byte{0x80},
			expected: []byte{},
		},
		"empty list": {
			input:    []byte{0xc0},
			expected: []any{},
		},
		"integer 0": {
			input:    []byte{0x80},
			expected: []byte{},
		},
		"integer 15": {
			input:    []byte{0x0f},
			expected: []byte{0x0f},
		},
		"integer 1024": {
			input:    []byte{0x82, 0x04, 0x00},
			expected: []byte{0x04, 0x00},
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
			input:    []byte{0x82, 0x04, 0x00},
			expected: []byte{0x04, 0x00},
		},
		"set theoretical representation [ [], [[]], [ [], [[]] ] ]": {
			input: []byte{0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0},
			expected: []any{
				[]any{},
				[]any{[]any{}},
				[]any{
					[]any{},
					[]any{[]any{}},
				},
			},
		},
		"long string Lorem ipsum...": {
			input:    append([]byte{0xb8, 0x38}, []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit")...),
			expected: []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := RlpDecode(tt.input)
			if err != nil {
				t.Fatalf("RlpDecode failed: %v", err)
			}

			switch expected := tt.expected.(type) {
			case []byte:
				actual, ok := val.([]byte)
				if !ok {
					t.Fatalf("Expected []byte, got %T", val)
				}
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Decoded []byte = %x, want %x", actual, expected)
				}

			case []any:
				actual, ok := val.([]any)
				if !ok {
					t.Fatalf("Expected []any, got %T", val)
				}
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Decoded list = %#v\nExpected list = %#v", actual, expected)
				}

			default:
				t.Fatalf("Unsupported expected type: %T", expected)
			}
		})
	}
}
