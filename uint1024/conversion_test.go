package uint1024

import (
	"bytes"
	"reflect"
	"testing"
)

// TestFromLimbs tests the FromLimbs constructor
func TestFromLimbs(t *testing.T) {
	tests := []struct {
		name     string
		limbs    []uint64
		expected *Uint1024
	}{
		{
			name:     "Empty slice",
			limbs:    []uint64{},
			expected: ZERO.Clone(),
		},
		{
			name:     "Single limb",
			limbs:    []uint64{42},
			expected: New(42),
		},
		{
			name:     "Multiple limbs",
			limbs:    []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected: &Uint1024{words: [16]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
		},
		{
			name:     "Too many limbs (truncated)",
			limbs:    []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			expected: &Uint1024{words: [16]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
		},
		{
			name:     "Partial limbs",
			limbs:    []uint64{100, 200, 300},
			expected: &Uint1024{words: [16]uint64{100, 200, 300, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromLimbs(tt.limbs)
			if !result.Equal(tt.expected) {
				t.Errorf("FromLimbs(%v) = %v, want %v", tt.limbs, result.ToLimbs(), tt.expected.ToLimbs())
			}
		})
	}
}

// TestToLimbs tests the ToLimbs method
func TestToLimbs(t *testing.T) {
	tests := []struct {
		name     string
		input    *Uint1024
		expected []uint64
	}{
		{
			name:     "Zero",
			input:    ZERO.Clone(),
			expected: make([]uint64, 16),
		},
		{
			name:     "One",
			input:    ONE.Clone(),
			expected: append([]uint64{1}, make([]uint64, 15)...),
		},
		{
			name:     "Custom value",
			input:    &Uint1024{words: [16]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
			expected: []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToLimbs()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToLimbs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFromLeBytes tests little-endian byte conversion
func TestFromLeBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected *Uint1024
	}{
		{
			name:     "Empty bytes",
			bytes:    []byte{},
			expected: ZERO.Clone(),
		},
		{
			name:     "Single byte",
			bytes:    []byte{42},
			expected: New(42),
		},
		{
			name:     "8 bytes (one word)",
			bytes:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
			expected: &Uint1024{words: [16]uint64{0x0807060504030201, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
		{
			name:     "16 bytes (two words)",
			bytes:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected: &Uint1024{words: [16]uint64{0x0807060504030201, 0x100f0e0d0c0b0a09, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromLeBytes(tt.bytes)
			if !result.Equal(tt.expected) {
				t.Errorf("FromLeBytes(%v) = %v, want %v", tt.bytes, result.ToLimbs(), tt.expected.ToLimbs())
			}
		})
	}
}

// TestToLeBytes tests little-endian byte output
func TestToLeBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    *Uint1024
		expected []byte
	}{
		{
			name:     "Zero",
			input:    ZERO.Clone(),
			expected: make([]byte, 128),
		},
		{
			name:     "Value 42",
			input:    New(42),
			expected: append([]byte{42}, make([]byte, 127)...),
		},
		{
			name:  "Custom value",
			input: &Uint1024{words: [16]uint64{0x0807060504030201, 0x100f0e0d0c0b0a09, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			expected: append(
				[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				make([]byte, 112)...,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToLeBytes()
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ToLeBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFromBeBytes tests big-endian byte conversion
func TestFromBeBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected *Uint1024
	}{
		{
			name:     "Empty bytes",
			bytes:    []byte{},
			expected: ZERO.Clone(),
		},
		{
			name:     "Single byte",
			bytes:    []byte{42},
			expected: &Uint1024{words: [16]uint64{42, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
		{
			name:     "8 bytes (one word)",
			bytes:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
			expected: &Uint1024{words: [16]uint64{0x0102030405060708, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromBeBytes(tt.bytes)
			if !result.Equal(tt.expected) {
				t.Errorf("FromBeBytes(%v) = %v, want %v", tt.bytes, result.ToLimbs(), tt.expected.ToLimbs())
			}
		})
	}
}

// TestToBeBytes tests big-endian byte output
func TestToBeBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    *Uint1024
		expected []byte
	}{
		{
			name:     "Zero",
			input:    ZERO.Clone(),
			expected: make([]byte, 128),
		},
		{
			name:  "Value 42",
			input: New(42),
			expected: append(
				make([]byte, 127),
				42,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToBeBytes()
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ToBeBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestRoundTripConversions tests that conversions are reversible
func TestRoundTripConversions(t *testing.T) {
	original := &Uint1024{words: [16]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}}

	// Test limbs round trip
	limbs := original.ToLimbs()
	fromLimbs := FromLimbs(limbs)
	if !original.Equal(fromLimbs) {
		t.Error("Limbs round trip failed")
	}

	// Test little-endian bytes round trip
	leBytes := original.ToLeBytes()
	fromLeBytes := FromLeBytes(leBytes)
	if !original.Equal(fromLeBytes) {
		t.Error("Little-endian bytes round trip failed")
	}

	// Test big-endian bytes round trip
	beBytes := original.ToBeBytes()
	fromBeBytes := FromBeBytes(beBytes)
	if !original.Equal(fromBeBytes) {
		t.Error("Big-endian bytes round trip failed")
	}
}
