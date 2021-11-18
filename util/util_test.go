package util

import (
	"testing"
)

func TestUint16ToBytes(t *testing.T) {
	var value uint16 = 69
	bytes := Uint16ToBytes(value)
	decodedValue := BytesToUint16(bytes)

	if value != decodedValue {
		t.Fatalf(`Values did not match after decoding from bytes. Original: %v, Decoded: %v`, value, decodedValue)
	}
}

func TestFloat64ToBytes(t *testing.T) {
	var value float64 = 69.420
	bytes := Float64ToBytes(value)
	decodedValue := BytesToFloat64(bytes)

	if value != decodedValue {
		t.Fatalf(`Values did not match after decoding from bytes. Original: %v, Decoded: %v`, value, decodedValue)
	}
}
