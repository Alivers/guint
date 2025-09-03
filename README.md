# guint - Large Unsigned Integer Library for Go

A high-performance Go library providing separate packages for 512-bit and 1024-bit unsigned integer types with comprehensive arithmetic, bitwise, and comparison operations.

## Features

- **uint512 Package**: 512-bit unsigned integer implementation
- **uint1024 Package**: 1024-bit unsigned integer implementation
- **Independent Packages**: Each package can be used independently
- **Global Constants**: ZERO, ONE, MAX as global variables in each package
- **Arithmetic Operations**: Addition, subtraction, multiplication, division, modulo
- **Bitwise Operations**: AND, OR, XOR, NOT, left/right shifts
- **Comparison Operations**: Equal, less than, greater than, etc.
- **String Conversion**: Decimal and hexadecimal representation
- **Bit Manipulation**: Individual bit operations, bit counting
- **High Performance**: Optimized algorithms with in-place operations
- **Comprehensive Testing**: Extensive test coverage with benchmarks

## Installation

```bash
go get github.com/Alivers/guint
```

## Quick Start

### Using uint512 package

```go
package main

import (
    "fmt"
    "github.com/Alivers/guint/uint512"
)

func main() {
    // Create new 512-bit integers
    a := uint512.New(123456789)
    b := uint512.New(987654321)

    // Use global constants
    zero := uint512.ZERO
    one := uint512.ONE
    max := uint512.MAX

    // Arithmetic operations
    sum := a.Add(b)
    fmt.Printf("Sum: %s\n", sum.String())

    // Bitwise operations
    result := a.And(b)
    fmt.Printf("Bitwise AND: %s\n", result.String())

    // Comparison with global constants
    if a.Greater(zero) {
        fmt.Println("a is greater than zero")
    }
}
```

### Using uint1024 package

```go
package main

import (
    "fmt"
    "github.com/Alivers/guint/uint1024"
)

func main() {
    // Create new 1024-bit integers
    a := uint1024.New(123456789)
    b := uint1024.New(987654321)

    // Use global constants
    zero := uint1024.ZERO
    one := uint1024.ONE
    max := uint1024.MAX

    // Arithmetic operations
    product := a.Mul(b)
    fmt.Printf("Product: %s\n", product.String())

    // Comparison with global constants
    if product.Less(max) {
        fmt.Println("Product is less than maximum value")
    }
}
```

## Package Structure

The library is split into two independent packages:

- `uint512/` - Contains all 512-bit integer functionality
- `uint1024/` - Contains all 1024-bit integer functionality

Each package provides the same API interface but operates on different bit widths.

## API Reference

### Global Constants

Each package provides global constants for common values:

```go
// Available in both uint512 and uint1024 packages
var (
    ZERO *Uint512  // Zero value
    ONE  *Uint512  // Value 1
    MAX  *Uint512  // Maximum value (all bits set)
)
```

### Constructors

```go
// Available in both packages
a := uint512.New(42)        // Create from uint64

// Use global constants
zero := uint512.ZERO        // Reference global zero
one := uint512.ONE          // Reference global one
max := uint512.MAX          // Reference global max

// Clone when you need to modify
a := uint512.ZERO.Clone()   // Clone from global constant

// Create from limbs (uint64 slice)
limbs := []uint64{1, 2, 3, 4, 5, 6, 7, 8}
num := uint512.FromLimbs(limbs)

// Create from byte slices
leData := []byte{1, 2, 3, 4, 5, 6, 7, 8}     // Little-endian
beData := []byte{1, 2, 3, 4, 5, 6, 7, 8}     // Big-endian
numLE := uint512.FromLeBytes(leData)
numBE := uint512.FromBeBytes(beData)
```

### Arithmetic Operations

```go
// Same API for both uint512 and uint1024
sum := a.Add(b)
diff := a.Sub(b)
product := a.Mul(b)         // Returns same type
quotient, err := a.Div(b)
mod, err := a.Mod(b)

// In-place operations
a.AddInPlace(b)
a.SubInPlace(b)
```

### Bitwise Operations

```go
// Logical operations (same API for both packages)
and := a.And(b)
or := a.Or(b)
xor := a.Xor(b)
not := a.Not()

// In-place operations
a.AndInPlace(b)
a.OrInPlace(b)
a.XorInPlace(b)
a.NotInPlace()

// Shift operations
left := a.Shl(10)
right := a.Shr(10)

// In-place shifts
a.ShlInPlace(10)
a.ShrInPlace(10)
```

### Bit Manipulation

```go
// Individual bit operations (same API for both packages)
bit := a.Bit(5)        // Get bit at position 5
a.SetBit(5)            // Set bit 5 to 1
a.ClearBit(5)          // Set bit 5 to 0
a.FlipBit(5)           // Flip bit 5

// Bit counting
leadingZeros := a.LeadingZeros()
trailingZeros := a.TrailingZeros()
onesCount := a.OnesCount()
```

### Comparison Operations

```go
// Same API for both packages
equal := a.Equal(b)
less := a.Less(b)
greater := a.Greater(b)
cmp := a.Compare(b)     // Returns -1, 0, or 1

// Utility comparisons
isOdd := a.IsOdd()
isEven := a.IsEven()
min := a.Min(b)
max := a.Max(b)

// Compare with global constants
if a.Equal(uint512.ZERO) {
    fmt.Println("a is zero")
}
```

### String Conversion and Data Export

```go
// Same API for both packages
decimal := a.String()           // Decimal string
hex := a.Hex()                  // Hexadecimal string
isZero := a.IsZero()           // Zero check

// Export to different formats
limbs := a.ToLimbs()           // Export as uint64 slice
leBytes := a.ToLeBytes()       // Export as little-endian bytes
beBytes := a.ToBeBytes()       // Export as big-endian bytes
```

## Examples

### Working with Global Constants

```go
package main

import (
    "fmt"
    "github.com/Alivers/guint/uint512"
)

func main() {
    // Use global constants directly
    fmt.Printf("Zero: %s\n", uint512.ZERO.String())
    fmt.Printf("One: %s\n", uint512.ONE.String())
    fmt.Printf("Max: %s\n", uint512.MAX.Hex())

    // Operations with global constants
    a := uint512.New(42)

    // Check if zero
    if a.Equal(uint512.ZERO) {
        fmt.Println("a is zero")
    } else {
        fmt.Println("a is not zero")
    }

    // Add one
    result, _ := a.Add(uint512.ONE)
    fmt.Printf("42 + 1 = %s\n", result.String())

    // Check if maximum
    if a.Equal(uint512.MAX) {
        fmt.Println("a is at maximum value")
    }
}
```

### Data Conversion Examples

```go
package main

import (
    "fmt"
    "github.com/Alivers/guint/uint512"
)

func main() {
    // Create from limbs
    limbs := []uint64{0x1234567890abcdef, 0xfedcba0987654321}
    num := uint512.FromLimbs(limbs)
    fmt.Printf("From limbs: %s\n", num.Hex())

    // Round-trip conversion
    exported := num.ToLimbs()
    restored := uint512.FromLimbs(exported)
    fmt.Printf("Round-trip equal: %t\n", num.Equal(restored))

    // Byte conversion
    leBytes := num.ToLeBytes()
    beBytes := num.ToBeBytes()

    fromLE := uint512.FromLeBytes(leBytes)
    fromBE := uint512.FromBeBytes(beBytes)

    fmt.Printf("LE conversion equal: %t\n", num.Equal(fromLE))
    fmt.Printf("BE conversion equal: %t\n", num.Equal(fromBE))

    // Working with partial data
    shortData := []byte{1, 2, 3, 4}
    partial := uint512.FromLeBytes(shortData)
    fmt.Printf("Partial data: %s\n", partial.Hex())
}
```

### Cross-Package Operations

Since the packages are independent, you need to convert between them when needed:

```go
package main

import (
    "github.com/Alivers/guint/uint512"
    "github.com/Alivers/guint/uint1024"
)

// Convert uint512 to uint1024
func convert512to1024(u512 *uint512.Uint512) *uint1024.Uint1024 {
    // This is a simplified conversion - you'd need to implement proper conversion
    // by copying the internal words representation
    result := uint1024.ZERO.Clone()
    // Copy implementation would go here
    return result
}

func main() {
    // Work with both types independently
    a512 := uint512.New(12345)
    a1024 := uint1024.New(12345)

    // Each package has its own operations
    sum512, _ := a512.Add(uint512.ONE)
    sum1024, _ := a1024.Add(uint1024.ONE)

    // Results are in their respective types
    fmt.Printf("512-bit result: %s\n", sum512.String())
    fmt.Printf("1024-bit result: %s\n", sum1024.String())
}
```

### Cryptographic Example

```go
package main

import (
    "fmt"
    "github.com/Alivers/guint/uint512"
)

// Simple modular exponentiation using uint512
func modExp(base, exp, mod *uint512.Uint512) *uint512.Uint512 {
    result := uint512.ONE.Clone()
    base = base.Clone()
    exp = exp.Clone()

    for !exp.Equal(uint512.ZERO) {
        if exp.IsOdd() {
            // result = (result * base) mod mod
            product := result.Mul(base)
            // Convert Uint1024 back to Uint512 with modulo
            // (This is simplified - proper implementation would handle the conversion)
            result, _ = result.Mod(mod)
        }

        // base = (base * base) mod mod
        product := base.Mul(base)
        // Convert and mod (simplified)
        base, _ = base.Mod(mod)

        exp.ShrInPlace(1)
    }

    return result
}

func main() {
    base := uint512.New(2)
    exp := uint512.New(10)
    mod := uint512.New(1000)

    result := modExp(base, exp, mod)
    fmt.Printf("2^10 mod 1000 = %s\n", result.String())
}
```

## Performance

Both packages are optimized for performance with:

- **Independent compilation**: Each package compiles separately
- **In-place operations** to reduce memory allocations
- **Word-level operations** using native uint64 arithmetic
- **Optimized algorithms** for multiplication and division
- **Minimal overhead** for common operations

Run tests for each package:

```bash
cd uint512 && go test -v
cd uint1024 && go test -v
```

## Implementation Details

### Package Independence

- Each package is completely independent
- No shared code between uint512 and uint1024 packages
- Global constants (ZERO, ONE, MAX) are defined separately in each package
- Identical API interface for both packages

### Memory Layout

Both packages use little-endian word order:

- `words[0]` contains the least significant 64 bits
- `words[n-1]` contains the most significant 64 bits

### Global Constants

Global constants are initialized at package load time:

- `ZERO`: All words set to 0
- `ONE`: First word set to 1, others to 0
- `MAX`: All words set to ^uint64(0) (all bits 1)

## Changelog

### v2.0.0

- **Breaking Change**: Split into separate uint512 and uint1024 packages
- Added global constants ZERO, ONE, MAX to each package
- Packages are now independent and can be used separately
- Same API interface maintained for both packages
- Improved performance through package separation

### v1.0.0

- Initial monolithic implementation
- Combined uint512 and uint1024 in single package

## Migration from v1.0.0

To migrate from the combined package:

```go
// Old (v1.0.0)
import "github.com/Alivers/guint"
a := guint.NewUint512(42)
b := guint.NewUint1024(42)

// New (v2.0.0)
import (
    "github.com/Alivers/guint/uint512"
    "github.com/Alivers/guint/uint1024"
)
a := uint512.New(42)
b := uint1024.New(42)

// Use global constants
zero512 := uint512.ZERO
zero1024 := uint1024.ZERO
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
