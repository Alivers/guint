// comparison.go implements comparison operations for Uint1024
package uint1024

// Equal returns true if a == b.
func (u *Uint1024) Equal(other *Uint1024) bool {
	for i := range u.words {
		if u.words[i] != other.words[i] {
			return false
		}
	}
	return true
}

// Less returns true if a < b.
func (u *Uint1024) Less(other *Uint1024) bool {
	// Compare from most significant word to least significant
	for i := len(u.words) - 1; i >= 0; i-- {
		if u.words[i] < other.words[i] {
			return true
		}
		if u.words[i] > other.words[i] {
			return false
		}
	}
	return false // Equal
}

// LessOrEqual returns true if a <= b.
func (u *Uint1024) LessOrEqual(other *Uint1024) bool {
	return u.Less(other) || u.Equal(other)
}

// Greater returns true if a > b.
func (u *Uint1024) Greater(other *Uint1024) bool {
	return other.Less(u)
}

// GreaterOrEqual returns true if a >= b.
func (u *Uint1024) GreaterOrEqual(other *Uint1024) bool {
	return u.Greater(other) || u.Equal(other)
}

// NotEqual returns true if a != b.
func (u *Uint1024) NotEqual(other *Uint1024) bool {
	return !u.Equal(other)
}

// Compare returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func (u *Uint1024) Compare(other *Uint1024) int {
	// Compare from most significant word to least significant
	for i := len(u.words) - 1; i >= 0; i-- {
		if u.words[i] < other.words[i] {
			return -1
		}
		if u.words[i] > other.words[i] {
			return 1
		}
	}
	return 0 // Equal
}

// IsOdd returns true if the number is odd.
func (u *Uint1024) IsOdd() bool {
	return u.words[0]&1 == 1
}

// IsEven returns true if the number is even.
func (u *Uint1024) IsEven() bool {
	return u.words[0]&1 == 0
}

// Min returns the smaller of two numbers.
func (u *Uint1024) Min(other *Uint1024) *Uint1024 {
	if u.Less(other) {
		return u.Clone()
	}
	return other.Clone()
}

// Max returns the larger of two numbers.
func (u *Uint1024) Max(other *Uint1024) *Uint1024 {
	if u.Greater(other) {
		return u.Clone()
	}
	return other.Clone()
}
