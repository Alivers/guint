// bitwise.go implements bitwise operations for Uint512
package uint512

import "math/bits"

// And performs bitwise AND: result = a & b.
func (u *Uint512) And(other *Uint512) *Uint512 {
	result := &Uint512{}
	for i := range u.words {
		result.words[i] = u.words[i] & other.words[i]
	}
	return result
}

// AndInPlace performs bitwise AND in place: u = u & other.
func (u *Uint512) AndInPlace(other *Uint512) {
	for i := range u.words {
		u.words[i] &= other.words[i]
	}
}

// Or performs bitwise OR: result = a | b.
func (u *Uint512) Or(other *Uint512) *Uint512 {
	result := &Uint512{}
	for i := range u.words {
		result.words[i] = u.words[i] | other.words[i]
	}
	return result
}

// OrInPlace performs bitwise OR in place: u = u | other.
func (u *Uint512) OrInPlace(other *Uint512) {
	for i := range u.words {
		u.words[i] |= other.words[i]
	}
}

// Xor performs bitwise XOR: result = a ^ b.
func (u *Uint512) Xor(other *Uint512) *Uint512 {
	result := &Uint512{}
	for i := range u.words {
		result.words[i] = u.words[i] ^ other.words[i]
	}
	return result
}

// XorInPlace performs bitwise XOR in place: u = u ^ other.
func (u *Uint512) XorInPlace(other *Uint512) {
	for i := range u.words {
		u.words[i] ^= other.words[i]
	}
}

// Not performs bitwise NOT: result = ^a.
func (u *Uint512) Not() *Uint512 {
	result := &Uint512{}
	for i := range u.words {
		result.words[i] = ^u.words[i]
	}
	return result
}

// NotInPlace performs bitwise NOT in place: u = ^u.
func (u *Uint512) NotInPlace() {
	for i := range u.words {
		u.words[i] = ^u.words[i]
	}
}

// Shl performs left shift: result = a << n.
func (u *Uint512) Shl(n uint) *Uint512 {
	result := u.Clone()
	result.ShlInPlace(n)
	return result
}

// ShlInPlace performs left shift in place: u = u << n.
func (u *Uint512) ShlInPlace(n uint) {
	if n == 0 {
		return
	}

	if n >= 512 {
		// All bits are shifted out
		for i := range u.words {
			u.words[i] = 0
		}
		return
	}

	wordShift := n / 64
	bitShift := n % 64

	if wordShift > 0 {
		// Shift entire words
		for i := len(u.words) - 1; i >= int(wordShift); i-- {
			u.words[i] = u.words[i-int(wordShift)]
		}
		for i := 0; i < int(wordShift); i++ {
			u.words[i] = 0
		}
	}

	if bitShift > 0 {
		// Shift bits within words
		carry := uint64(0)
		for i := int(wordShift); i < len(u.words); i++ {
			newCarry := u.words[i] >> (64 - bitShift)
			u.words[i] = (u.words[i] << bitShift) | carry
			carry = newCarry
		}
	}
}

// Shr performs right shift: result = a >> n.
func (u *Uint512) Shr(n uint) *Uint512 {
	result := u.Clone()
	result.ShrInPlace(n)
	return result
}

// ShrInPlace performs right shift in place: u = u >> n.
func (u *Uint512) ShrInPlace(n uint) {
	if n == 0 {
		return
	}

	if n >= 512 {
		// All bits are shifted out
		for i := range u.words {
			u.words[i] = 0
		}
		return
	}

	wordShift := n / 64
	bitShift := n % 64

	if wordShift > 0 {
		// Shift entire words
		for i := 0; i < len(u.words)-int(wordShift); i++ {
			u.words[i] = u.words[i+int(wordShift)]
		}
		for i := len(u.words) - int(wordShift); i < len(u.words); i++ {
			u.words[i] = 0
		}
	}

	if bitShift > 0 {
		// Shift bits within words
		carry := uint64(0)
		for i := len(u.words) - int(wordShift) - 1; i >= 0; i-- {
			newCarry := u.words[i] << (64 - bitShift)
			u.words[i] = (u.words[i] >> bitShift) | carry
			carry = newCarry
		}
	}
}

// Bit returns the value of the bit at position i (0 is least significant).
func (u *Uint512) Bit(i int) bool {
	if i < 0 || i >= 512 {
		return false
	}
	wordIndex := i / 64
	bitIndex := i % 64
	return (u.words[wordIndex] & (1 << bitIndex)) != 0
}

// SetBit sets the bit at position i to 1.
func (u *Uint512) SetBit(i int) {
	if i < 0 || i >= 512 {
		return
	}
	wordIndex := i / 64
	bitIndex := i % 64
	u.words[wordIndex] |= (1 << bitIndex)
}

// ClearBit sets the bit at position i to 0.
func (u *Uint512) ClearBit(i int) {
	if i < 0 || i >= 512 {
		return
	}
	wordIndex := i / 64
	bitIndex := i % 64
	u.words[wordIndex] &^= (1 << bitIndex)
}

// FlipBit flips the bit at position i.
func (u *Uint512) FlipBit(i int) {
	if i < 0 || i >= 512 {
		return
	}
	wordIndex := i / 64
	bitIndex := i % 64
	u.words[wordIndex] ^= (1 << bitIndex)
}

// LeadingZeros returns the number of leading zero bits.
func (u *Uint512) LeadingZeros() int {
	for i := len(u.words) - 1; i >= 0; i-- {
		if u.words[i] != 0 {
			return (len(u.words)-1-i)*64 + bits.LeadingZeros64(u.words[i])
		}
	}
	return 512
}

// TrailingZeros returns the number of trailing zero bits.
func (u *Uint512) TrailingZeros() int {
	for i := 0; i < len(u.words); i++ {
		if u.words[i] != 0 {
			return i*64 + bits.TrailingZeros64(u.words[i])
		}
	}
	return 512
}

// OnesCount returns the number of one bits (population count).
func (u *Uint512) OnesCount() int {
	count := 0
	for _, word := range u.words {
		count += bits.OnesCount64(word)
	}
	return count
}
