// arithmetic.go implements arithmetic operations for Uint1024
package uint1024

import (
	"fmt"
	"math/bits"
)

// Add performs addition: result = a + b.
func (u *Uint1024) Add(other *Uint1024) *Uint1024 {
	result := &Uint1024{}
	var carry uint64

	for i := range u.words {
		sum, c1 := bits.Add64(u.words[i], other.words[i], carry)
		result.words[i] = sum
		carry = c1
	}

	return result
}

// AddInPlace performs addition in place: u = u + other.
func (u *Uint1024) AddInPlace(other *Uint1024) {
	var carry uint64

	for i := range u.words {
		sum, c1 := bits.Add64(u.words[i], other.words[i], carry)
		u.words[i] = sum
		carry = c1
	}
}

// Sub performs subtraction: result = a - b.
func (u *Uint1024) Sub(other *Uint1024) *Uint1024 {
	result := &Uint1024{}
	var borrow uint64

	for i := range u.words {
		diff, b1 := bits.Sub64(u.words[i], other.words[i], borrow)
		result.words[i] = diff
		borrow = b1
	}

	return result
}

// SubInPlace performs subtraction in place: u = u - other.
func (u *Uint1024) SubInPlace(other *Uint1024) {
	var borrow uint64

	for i := range u.words {
		diff, b1 := bits.Sub64(u.words[i], other.words[i], borrow)
		u.words[i] = diff
		borrow = b1
	}
}

// Mul performs multiplication: result = a * b.
// Note: This truncates the result to fit in Uint1024.
// In practice, you might want to return an error or handle overflow differently.
func (u *Uint1024) Mul(other *Uint1024) *Uint1024 {
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

// Div performs division: result = a / b.
// Returns quotient and error (if divisor is zero).
func (u *Uint1024) Div(other *Uint1024) (*Uint1024, error) {
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
	for i := 1023; i >= 0; i-- {
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
func (u *Uint1024) Mod(other *Uint1024) (*Uint1024, error) {
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
	for i := 1023; i >= 0; i-- {
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
