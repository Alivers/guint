package uint512

import (
	"testing"
)

// TestUint512Creation tests basic creation and initialization of Uint512
func TestUint512Creation(t *testing.T) {
	// Test New
	u1 := New(42)
	if u1.words[0] != 42 {
		t.Errorf("New(42) failed: got %d, want 42", u1.words[0])
	}
	for i := 1; i < len(u1.words); i++ {
		if u1.words[i] != 0 {
			t.Errorf("New(42) word[%d] should be 0, got %d", i, u1.words[i])
		}
	}

	// Test global constants directly
	if !ZERO.IsZero() {
		t.Error("Global ZERO should be zero value")
	}

	if ONE.words[0] != 1 {
		t.Error("Global ONE should have first word as 1")
	}
	for i := 1; i < len(ONE.words); i++ {
		if ONE.words[i] != 0 {
			t.Errorf("Global ONE word[%d] should be 0, got %d", i, ONE.words[i])
		}
	}

	for i := 0; i < len(MAX.words); i++ {
		if MAX.words[i] != ^uint64(0) {
			t.Errorf("Global MAX word[%d] should be all 1s", i)
		}
	}
}

// TestGlobalConstants tests the global ZERO, ONE, MAX constants
func TestGlobalConstants(t *testing.T) {
	// Test ZERO
	if !ZERO.IsZero() {
		t.Error("Global ZERO should be zero")
	}

	// Test ONE
	if ONE.words[0] != 1 {
		t.Error("Global ONE should have first word as 1")
	}
	for i := 1; i < len(ONE.words); i++ {
		if ONE.words[i] != 0 {
			t.Errorf("Global ONE word[%d] should be 0", i)
		}
	}

	// Test MAX
	for i := 0; i < len(MAX.words); i++ {
		if MAX.words[i] != ^uint64(0) {
			t.Errorf("Global MAX word[%d] should be all 1s", i)
		}
	}
}

// TestArithmetic tests basic arithmetic operations
func TestArithmetic(t *testing.T) {
	a := New(100)
	b := New(200)

	// Test addition
	sum := a.Add(b)
	expected := New(300)
	if !sum.Equal(expected) {
		t.Errorf("100 + 200: got %s, want %s", sum.String(), expected.String())
	}

	// Test subtraction
	diff := b.Sub(a)
	expected = New(100)
	if !diff.Equal(expected) {
		t.Errorf("200 - 100: got %s, want %s", diff.String(), expected.String())
	}

	// Test multiplication
	product := a.Mul(b)
	if product.String() != "20000" {
		t.Errorf("100 * 200: got %s, want 20000", product.String())
	}
}

// TestBitwise tests bitwise operations
func TestBitwise(t *testing.T) {
	a := New(0b1100) // 12
	b := New(0b1010) // 10

	// Test AND
	result := a.And(b)
	expected := New(0b1000) // 8
	if !result.Equal(expected) {
		t.Errorf("12 & 10: got %s, want %s", result.String(), expected.String())
	}

	// Test OR
	result = a.Or(b)
	expected = New(0b1110) // 14
	if !result.Equal(expected) {
		t.Errorf("12 | 10: got %s, want %s", result.String(), expected.String())
	}

	// Test XOR
	result = a.Xor(b)
	expected = New(0b0110) // 6
	if !result.Equal(expected) {
		t.Errorf("12 ^ 10: got %s, want %s", result.String(), expected.String())
	}

	// Test NOT
	result = a.Not()
	// NOT should flip all bits
	if result.And(a).IsZero() && result.Or(a).Equal(MAX) {
		// This is expected behavior
	} else {
		t.Error("NOT operation failed")
	}
}

// TestComparison tests comparison operations
func TestComparison(t *testing.T) {
	a := New(100)
	b := New(200)
	c := New(100)

	// Test Equal
	if !a.Equal(c) {
		t.Error("Equal numbers should be equal")
	}
	if a.Equal(b) {
		t.Error("Different numbers should not be equal")
	}

	// Test Less
	if !a.Less(b) {
		t.Error("100 should be less than 200")
	}
	if a.Less(c) {
		t.Error("100 should not be less than 100")
	}

	// Test Greater
	if !b.Greater(a) {
		t.Error("200 should be greater than 100")
	}
	if a.Greater(c) {
		t.Error("100 should not be greater than 100")
	}

	// Test Compare
	if a.Compare(b) != -1 {
		t.Error("100.Compare(200) should return -1")
	}
	if a.Compare(c) != 0 {
		t.Error("100.Compare(100) should return 0")
	}
	if b.Compare(a) != 1 {
		t.Error("200.Compare(100) should return 1")
	}
}

// TestStringConversion tests string representation
func TestStringConversion(t *testing.T) {
	tests := []struct {
		value    uint64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{123, "123"},
		{^uint64(0), "18446744073709551615"},
	}

	for _, test := range tests {
		u := New(test.value)
		result := u.String()
		if result != test.expected {
			t.Errorf("String() for %d: got %s, want %s", test.value, result, test.expected)
		}
	}
}

// TestHexConversion tests hexadecimal representation
func TestHexConversion(t *testing.T) {
	tests := []struct {
		value    uint64
		expected string
	}{
		{0, "0x0"},
		{15, "0xf"},
		{255, "0xff"},
		{256, "0x100"},
	}

	for _, test := range tests {
		u := New(test.value)
		result := u.Hex()
		if result != test.expected {
			t.Errorf("Hex() for %d: got %s, want %s", test.value, result, test.expected)
		}
	}
}

// TestBitOperations tests individual bit operations
func TestBitOperations(t *testing.T) {
	u := ZERO.Clone()

	// Test SetBit
	u.SetBit(5)
	if !u.Bit(5) {
		t.Error("SetBit(5) should set bit 5")
	}

	// Test ClearBit
	u.ClearBit(5)
	if u.Bit(5) {
		t.Error("ClearBit(5) should clear bit 5")
	}

	// Test FlipBit
	u.FlipBit(3)
	if !u.Bit(3) {
		t.Error("FlipBit(3) should set bit 3")
	}
	u.FlipBit(3)
	if u.Bit(3) {
		t.Error("FlipBit(3) again should clear bit 3")
	}
}

// TestShiftOperations tests shift operations
func TestShiftOperations(t *testing.T) {
	// Test left shift
	u := New(1)
	result := u.Shl(4)
	expected := New(16)
	if !result.Equal(expected) {
		t.Errorf("1 << 4: got %s, want %s", result.String(), expected.String())
	}

	// Test right shift
	u = New(16)
	result = u.Shr(4)
	expected = New(1)
	if !result.Equal(expected) {
		t.Errorf("16 >> 4: got %s, want %s", result.String(), expected.String())
	}
}
