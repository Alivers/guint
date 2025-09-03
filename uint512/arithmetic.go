// arithmetic.go implements arithmetic operations for Uint512
package uint512

import (
	"fmt"
	"math/bits"
)

// Add performs addition: result = a + b.
func (u *Uint512) Add(other *Uint512) *Uint512 {
	result := &Uint512{}
	var carry uint64

	for i := range u.words {
		sum, c1 := bits.Add64(u.words[i], other.words[i], carry)
		result.words[i] = sum
		carry = c1
	}

	return result
}

// AddInPlace performs addition in place: u = u + other.
func (u *Uint512) AddInPlace(other *Uint512) {
	var carry uint64

	for i := range u.words {
		sum, c1 := bits.Add64(u.words[i], other.words[i], carry)
		u.words[i] = sum
		carry = c1
	}
}

// Sub performs subtraction: result = a - b.
func (u *Uint512) Sub(other *Uint512) *Uint512 {
	result := &Uint512{}
	var borrow uint64

	for i := range u.words {
		diff, b1 := bits.Sub64(u.words[i], other.words[i], borrow)
		result.words[i] = diff
		borrow = b1
	}

	return result
}

// SubInPlace performs subtraction in place: u = u - other.
func (u *Uint512) SubInPlace(other *Uint512) {
	var borrow uint64

	for i := range u.words {
		diff, b1 := bits.Sub64(u.words[i], other.words[i], borrow)
		u.words[i] = diff
		borrow = b1
	}
}

// Uint1024 represents a 1024-bit result for multiplication
type Uint1024 struct {
	words [16]uint64
}

// Mul performs multiplication: result = a * b.
// Uses the schoolbook multiplication algorithm.
// Returns a Uint1024 to hold the full result.
func (u *Uint512) Mul(other *Uint512) *Uint1024 {
	result := &Uint1024{}

	for i := range u.words {
		if u.words[i] == 0 {
			continue
		}

		var carry uint64
		for j := 0; j < len(other.words) && i+j < len(result.words); j++ {
			if other.words[j] == 0 {
				continue
			}

			hi, lo := bits.Mul64(u.words[i], other.words[j])

			// Add lo to result[i+j]
			sum, c1 := bits.Add64(result.words[i+j], lo, carry)
			result.words[i+j] = sum
			carry = c1

			// Add hi to result[i+j+1] if it exists
			if i+j+1 < len(result.words) {
				sum, c2 := bits.Add64(result.words[i+j+1], hi, carry)
				result.words[i+j+1] = sum
				carry = c2
			}
		}

		// Propagate remaining carry
		k := i + len(other.words)
		for carry != 0 && k < len(result.words) {
			sum, c := bits.Add64(result.words[k], carry, 0)
			result.words[k] = sum
			carry = c
			k++
		}
	}

	return result
}

// String returns the decimal string representation of Uint1024.
func (u1024 *Uint1024) String() string {
	// Check if zero
	isZero := true
	for _, word := range u1024.words {
		if word != 0 {
			isZero = false
			break
		}
	}
	if isZero {
		return "0"
	}

	// Convert to decimal using repeated division by 10
	temp := &Uint1024{}
	copy(temp.words[:], u1024.words[:])
	var digits []byte

	for !temp.isZero() {
		remainder := temp.divBySmall(10)
		digits = append(digits, byte('0'+remainder))
	}

	// Reverse the digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return string(digits)
}

// isZero returns true if the Uint1024 is zero.
func (u1024 *Uint1024) isZero() bool {
	for _, word := range u1024.words {
		if word != 0 {
			return false
		}
	}
	return true
}

// divBySmall divides the Uint1024 by a small divisor and returns the remainder.
func (u1024 *Uint1024) divBySmall(divisor uint64) uint64 {
	var remainder uint64
	for i := len(u1024.words) - 1; i >= 0; i-- {
		dividend := remainder<<32 | u1024.words[i]>>32
		u1024.words[i] = (u1024.words[i] & 0xFFFFFFFF) | (dividend/divisor)<<32
		remainder = dividend % divisor

		dividend = remainder<<32 | (u1024.words[i] & 0xFFFFFFFF)
		u1024.words[i] = (u1024.words[i] & 0xFFFFFFFF00000000) | (dividend / divisor)
		remainder = dividend % divisor
	}
	return remainder
}

// Div performs division: result = a / b.
// Returns quotient and error (if divisor is zero).
func (u *Uint512) Div(other *Uint512) (*Uint512, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}

	if u.Less(other) {
		return ZERO.Clone(), nil
	}

	if u.Equal(other) {
		return ONE.Clone(), nil
	}

	// Use binary long division
	quotient := ZERO.Clone()
	remainder := ZERO.Clone()

	// Process bits from most significant to least significant
	for i := 511; i >= 0; i-- {
		// Shift remainder left by 1
		remainder.ShlInPlace(1)

		// Set the least significant bit of remainder to the i-th bit of dividend
		if u.Bit(i) {
			remainder.words[0] |= 1
		}

		// If remainder >= divisor, subtract divisor and set quotient bit
		if !remainder.Less(other) {
			remainder.SubInPlace(other)
			quotient.SetBit(i)
		}
	}

	return quotient, nil
}

// Mod performs modulo operation: result = a % b.
func (u *Uint512) Mod(other *Uint512) (*Uint512, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}

	if u.Less(other) {
		return u.Clone(), nil
	}

	if u.Equal(other) {
		return ZERO.Clone(), nil
	}

	// Use binary long division to compute remainder
	remainder := ZERO.Clone()

	// Process bits from most significant to least significant
	for i := 511; i >= 0; i-- {
		// Shift remainder left by 1
		remainder.ShlInPlace(1)

		// Set the least significant bit of remainder to the i-th bit of dividend
		if u.Bit(i) {
			remainder.words[0] |= 1
		}

		// If remainder >= divisor, subtract divisor
		if !remainder.Less(other) {
			remainder.SubInPlace(other)
		}
	}

	return remainder, nil
}
