package encode

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
)

func RlpEncode(input any) []byte {
	switch v := input.(type) {
	case string:
		data := []byte(v)
		if len(data) == 1 && data[0] < 0x80 {
			return data
		}
		return append(encodeLength(len(data), 0x80), data...)

	case []byte:
		if len(v) == 1 && v[0] < 0x80 {
			return v
		}
		return append(encodeLength(len(v), 0x80), v...)

	default:
		// Handle slices of any type (e.g., []string, []int, []any)
		reflectedValue := reflect.ValueOf(input)
		kind := reflectedValue.Kind()

		if reflectedValue.Kind() == reflect.Slice {
			var output []byte
			for i := 0; i < reflectedValue.Len(); i++ {
				item := reflectedValue.Index(i).Interface()
				output = append(output, RlpEncode(item)...)
			}
			return append(encodeLength(len(output), 0xc0), output...)
		}

		// Handle all integer kinds (signed and unsigned)
		if isIntegerKind(kind) {
			n := toInt(reflectedValue)
			if n == 0 {
				return []byte{0x80}
			}
			return encodeInteger(n)
		}

		panic(fmt.Sprintf("unsupported type: %T", input))
	}
}

func encodeLength(length int, offset int) []byte {
	if length < 56 {
		return []byte{byte(length + offset)}
	}

	l := big.NewInt(int64(length))
	limit := new(big.Int).Lsh(big.NewInt(1), 64) // 2^64
	if l.Cmp(limit) >= 0 {
		panic("input too long")
	}

	bl := toBinary(length)
	return append([]byte{byte(len(bl) + offset + 55)}, bl...)
}

func toBinary(x int) []byte {
	if x == 0 {
		return []byte{}
	}
	var buf bytes.Buffer
	for x > 0 {
		buf.WriteByte(byte(x & 0xff))
		x >>= 8
	}
	// Reverse to make big-endian
	b := buf.Bytes()
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func encodeInteger(n int) []byte {
	if n < 0 {
		panic("RLP only supports unsigned integers")
	}
	buf := toBinary(n)
	if len(buf) == 1 && buf[0] < 0x80 {
		return buf
	}
	return append(encodeLength(len(buf), 0x80), buf...)
}

func isIntegerKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func toInt(v reflect.Value) int {
	// Convert to int (you can use int64 if you want bigger range)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(v.Uint())
	default:
		panic("not an integer kind")
	}
}
