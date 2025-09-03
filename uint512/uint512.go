// Package uint512 provides implementation of 512-bit unsigned integer
// with comprehensive arithmetic, bitwise, and comparison operations.
package uint512

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Uint512 represents a 512-bit unsigned integer.
// It's implemented as an array of 8 uint64 values, stored in little-endian order.
type Uint512 struct {
	// words stores the 512-bit value as 8 64-bit words in little-endian order
	// words[0] contains the least significant 64 bits
	// words[7] contains the most significant 64 bits
	words [8]uint64
}

// Global constants
var (
	// ZERO represents the zero value for Uint512
	ZERO = &Uint512{}

	// ONE represents the value 1 for Uint512
	ONE = &Uint512{words: [8]uint64{1, 0, 0, 0, 0, 0, 0, 0}}

	// MAX represents the maximum value for Uint512 (all bits set to 1)
	MAX = &Uint512{words: [8]uint64{^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0)}}
)

// New creates a new Uint512 from a uint64 value.
func New(val uint64) *Uint512 {
	u := &Uint512{}
	u.words[0] = val
	return u
}

// FromLimbs creates a new Uint512 from a slice of uint64 limbs in little-endian order.
// If the slice is longer than 8 elements, only the first 8 are used.
// If shorter, the remaining words are set to zero.
func FromLimbs(limbs []uint64) *Uint512 {
	u := &Uint512{}
	n := len(limbs)
	if n > 8 {
		n = 8
	}
	copy(u.words[:n], limbs[:n])
	return u
}

// FromLeBytes creates a new Uint512 from a byte slice in little-endian order.
// The byte slice should be exactly 64 bytes (512 bits).
// If shorter, it's padded with zeros. If longer, only the first 64 bytes are used.
func FromLeBytes(data []byte) *Uint512 {
	u := &Uint512{}

	// Ensure we don't read beyond the slice
	dataLen := len(data)
	if dataLen > 64 {
		dataLen = 64
	}

	// Convert bytes to words in little-endian order
	for i := 0; i < 8; i++ {
		start := i * 8
		end := start + 8

		if start < dataLen {
			// Determine how many bytes we can read for this word
			bytesToRead := 8
			if end > dataLen {
				bytesToRead = dataLen - start
			}

			// Create a temp slice with padding if needed
			wordBytes := make([]byte, 8)
			copy(wordBytes, data[start:start+bytesToRead])

			u.words[i] = binary.LittleEndian.Uint64(wordBytes)
		}
	}

	return u
}

// FromBeBytes creates a new Uint512 from a byte slice in big-endian order.
// The byte slice should be exactly 64 bytes (512 bits).
// If shorter, it's padded with zeros. If longer, only the first 64 bytes are used.
func FromBeBytes(data []byte) *Uint512 {
	u := &Uint512{}

	// Ensure we don't read beyond the slice
	dataLen := len(data)
	if dataLen > 64 {
		dataLen = 64
	}

	if dataLen == 0 {
		return u
	}

	// For big-endian input, we need to place the data at the high-order end
	// Pad the data to 64 bytes with leading zeros
	padded := make([]byte, 64)
	copy(padded[64-dataLen:], data[:dataLen])

	// Convert bytes to words in big-endian order
	for i := 0; i < 8; i++ {
		start := i * 8
		u.words[7-i] = binary.BigEndian.Uint64(padded[start : start+8])
	}

	return u
}

// Clone creates a copy of the Uint512.
func (u *Uint512) Clone() *Uint512 {
	result := &Uint512{}
	copy(result.words[:], u.words[:])
	return result
}

// IsZero returns true if the value is zero.
func (u *Uint512) IsZero() bool {
	return u.Equal(ZERO)
}

// ToLimbs returns the Uint512 as a slice of uint64 limbs in little-endian order.
// Returns a copy of the internal words slice.
func (u *Uint512) ToLimbs() []uint64 {
	limbs := make([]uint64, 8)
	copy(limbs, u.words[:])
	return limbs
}

// ToLeBytes returns the Uint512 as a 64-byte slice in little-endian order.
func (u *Uint512) ToLeBytes() []byte {
	bytes := make([]byte, 64)

	for i := range u.words {
		start := i * 8
		binary.LittleEndian.PutUint64(bytes[start:start+8], u.words[i])
	}

	return bytes
}

// ToBeBytes returns the Uint512 as a 64-byte slice in big-endian order.
func (u *Uint512) ToBeBytes() []byte {
	bytes := make([]byte, 64)

	// For big-endian, we reverse the word order and use big-endian encoding
	for i := range u.words {
		wordIndex := 7 - i // Reverse word order for big-endian
		start := i * 8
		binary.BigEndian.PutUint64(bytes[start:start+8], u.words[wordIndex])
	}

	return bytes
}

// String returns the decimal string representation of the number.
func (u *Uint512) String() string {
	if u.IsZero() {
		return "0"
	}

	// Convert to decimal using repeated division by 10
	temp := u.Clone()
	var digits []byte

	for !temp.IsZero() {
		remainder := temp.divBySmall(10)
		digits = append(digits, byte('0'+remainder))
	}

	// Reverse the digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return string(digits)
}

// Hex returns the hexadecimal string representation of the number.
func (u *Uint512) Hex() string {
	if u.IsZero() {
		return "0x0"
	}

	var result strings.Builder
	result.WriteString("0x")

	// Find the most significant non-zero word
	msw := -1
	for i := len(u.words) - 1; i >= 0; i-- {
		if u.words[i] != 0 {
			msw = i
			break
		}
	}

	if msw == -1 {
		return "0x0"
	}

	// Write the most significant word without leading zeros
	result.WriteString(fmt.Sprintf("%x", u.words[msw]))

	// Write remaining words with leading zeros
	for i := msw - 1; i >= 0; i-- {
		result.WriteString(fmt.Sprintf("%016x", u.words[i]))
	}

	return result.String()
}

// divBySmall divides the number by a small divisor (< 2^64) and returns the remainder.
// This modifies the receiver in place.
func (u *Uint512) divBySmall(divisor uint64) uint64 {
	var remainder uint64
	for i := len(u.words) - 1; i >= 0; i-- {
		dividend := remainder<<32 | u.words[i]>>32
		u.words[i] = (u.words[i] & 0xFFFFFFFF) | (dividend/divisor)<<32
		remainder = dividend % divisor

		dividend = remainder<<32 | (u.words[i] & 0xFFFFFFFF)
		u.words[i] = (u.words[i] & 0xFFFFFFFF00000000) | (dividend / divisor)
		remainder = dividend % divisor
	}
	return remainder
}
