// Lesson 06: Collections — slices and maps
//
// Go has arrays (fixed size) but slices (dynamic) are used almost always.
// Maps are hash maps. No generic "list" — slices cover it.

package main

import (
	"fmt"
	"slices"
)

func main() {
	fmt.Println("=== Lesson 06: Collections ===")

	// --- Arrays (fixed size, rarely used directly) ---
	var arr [5]int
	arr[0] = 1
	arr2 := [3]string{"a", "b", "c"}
	fmt.Println("Array:", arr, arr2)

	// --- Slices (dynamic, backed by array) ---
	// Literal
	sl := []int{1, 2, 3, 4, 5}
	fmt.Println("Slice:", sl, "len:", len(sl), "cap:", cap(sl))

	// make: type, length, capacity (capacity optional)
	s := make([]int, 5)      // len=5, cap=5, zero-initialized
	s2 := make([]int, 0, 10) // len=0, cap=10 (preallocate)
	fmt.Println("make:", s, s2)

	// Append (like push)
	s = append(s, 6)
	s = append(s, 7, 8, 9)
	fmt.Println("After append:", s)

	// Append slice (must use ...)
	other := []int{10, 11}
	s = append(s, other...)
	fmt.Println("Append slice:", s)

	// Slicing: s[low:high] (high exclusive)
	sub := s[2:5]
	fmt.Println("s[2:5]:", sub)

	// Slices share underlying array (reference semantics)
	sub[0] = 999
	fmt.Println("After sub[0]=999, s:", s)

	// Copy
	dest := make([]int, 3)
	copy(dest, s)
	fmt.Println("Copy:", dest)

	// Delete element (slices.Delete returns modified slice)
	s = slices.Delete(s, 2, 3)
	fmt.Println("After delete index 2:", s)

	// --- Maps ---
	m := make(map[string]int)
	m["apple"] = 1
	m["banana"] = 2
	fmt.Println("Map:", m)

	// Literal
	ages := map[string]int{
		"alice": 30,
		"bob":   25,
	}
	fmt.Println("Ages:", ages)

	// Access: returns (value, ok)
	v, ok := ages["alice"]
	fmt.Println("alice:", v, ok)

	v, ok = ages["charlie"]
	fmt.Println("charlie:", v, ok) // 0, false

	// Delete
	delete(ages, "bob")
	fmt.Println("After delete bob:", ages)

	// Iteration order is random (intentionally)
	for k, v := range ages {
		fmt.Printf("%s: %d\n", k, v)
	}

	// --- Slice of structs ---
	type Item struct {
		Name  string
		Count int
	}
	items := []Item{
		{"A", 1},
		{"B", 2},
		{"C", 3},
	}
	for _, item := range items {
		fmt.Printf("%s: %d\n", item.Name, item.Count)
	}
}
