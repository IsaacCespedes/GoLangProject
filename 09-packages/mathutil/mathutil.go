// Package mathutil provides simple math utilities.
// Package name matches directory; exported names start with uppercase.
package mathutil

// Add adds two integers (exported: uppercase)
func Add(a, b int) int {
	return a + b
}

// multiply is unexported (lowercase) — not visible outside package
func multiply(a, b int) int {
	return a * b
}
