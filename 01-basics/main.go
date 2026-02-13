// Lesson 01: Basics — Hello World, variables, types
//
// Go is compiled, statically typed, with garbage collection.
// Unlike C/C++, no manual memory management. Unlike JS, no runtime type changes.

package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Lesson 01: Basics ===")

	// --- Variable declaration ---
	// Three styles; all equivalent for local variables:

	var x int = 42
	var y = 42   // type inferred
	z := 42      // short declaration (most common; only inside functions)

	fmt.Println("Variables:", x, y, z)

	// --- Zero values ---
	// Uninitialized variables get zero values (unlike C's undefined, or JS's undefined):
	var a int     // 0
	var b float64 // 0.0
	var c bool    // false
	var d string  // ""

	fmt.Printf("Zero values: int=%d float=%f bool=%v string=%q\n", a, b, c, d)

	// --- Basic types ---
	var i8 int8 = 127
	var i16 int16 = 32000
	var i32 int32 = 2147483647
	var i64 int64 = 9223372036854775807

	// int/uint: platform-dependent (32 or 64 bit)
	// In practice, use int for sizes and indices
	var n int = 100

	// Floating point
	var f32 float32 = 3.14
	var f64 float64 = 3.14159265358979

	// Rune = int32, represents a Unicode code point (like wchar_t in C++)
	var r rune = 'A'

	// Byte = uint8, for raw bytes
	var b1 byte = 0xFF

	fmt.Println("Types:", i8, i16, i32, i64, n, f32, f64, r, b1)

	// --- Type conversion ---
	// Explicit conversion required (no implicit narrowing like C)
	f := 3.14
	conv := int(f) // 3, truncates
	fmt.Println("Conversion float->int:", conv)

	// --- Constants ---
	const Pi = 3.14159
	const (
		Zero = 0
		One  = 1
	)
	// iota: auto-incrementing constant generator (like enum)
	const (
		Sunday = iota // 0
		Monday        // 1
		Tuesday       // 2
	)
	fmt.Println("Constants:", Pi, Sunday, Monday, Tuesday)

	// --- Strings ---
	s1 := "Hello"
	s2 := `Raw string
can span
multiple lines`
	s3 := s1 + " " + "World"
	fmt.Println(s3)
	fmt.Println("Raw:", s2)

	// Strings are immutable (like in most languages)
	// s1[0] = 'h' // compile error

	// --- Printf ---
	name := "Go"
	version := 1.22
	fmt.Printf("Language: %s, version: %.2f\n", name, version)
	fmt.Printf("Binary: %b, Hex: %x\n", 255, 255)

	// --- Multiple assignment ---
	p, q := 1, 2
	p, q = q, p // swap without temp (like Python)
	fmt.Println("Swap:", p, q)
}
