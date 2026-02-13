// Lesson 11 (Bonus): Generics — Go 1.18+
//
// Type parameters for reusable, type-safe code.
// Familiar if you know C++ templates or TypeScript generics.

package main

import (
	"fmt"
)

// Generic function: [T typeParam] before params
func Identity[T any](x T) T {
	return x
}

// Constrained type parameter (comparable = supports == and !=)
func Contains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// Constraint: ordered types (support < > etc)
func Min[T ~int | ~int64 | ~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Generic struct
type Box[T any] struct {
	Value T
}

func (b Box[T]) Get() T {
	return b.Value
}

func main() {
	fmt.Println("=== Lesson 11: Generics ===")

	fmt.Println("Identity(42):", Identity(42))
	fmt.Println("Identity(\"hello\"):", Identity("hello"))

	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Contains(nums, 3):", Contains(nums, 3))
	fmt.Println("Contains(nums, 99):", Contains(nums, 99))

	fmt.Println("Min(3, 7):", Min(3, 7))
	fmt.Println("Min(3.14, 2.71):", Min(3.14, 2.71))

	b := Box[int]{Value: 42}
	fmt.Println("Box.Get():", b.Get())
}
