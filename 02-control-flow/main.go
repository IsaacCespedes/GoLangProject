// Lesson 02: Control flow — if, for, switch
//
// Go simplifies control flow: no parentheses, braces always required.
// No while loop — use for. No ternary operator.

package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Lesson 02: Control Flow ===")

	// --- if ---
	// No parentheses; condition must be bool (no truthy/falsy like JS)
	x := 42
	if x > 0 {
		fmt.Println("x is positive")
	}

	// if with short statement (variable scoped to if block)
	if y := 10; y < 20 {
		fmt.Println("y is less than 20:", y)
	}
	// y not visible here

	// --- for: the only loop construct ---
	// C-style for
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// while-style (omit init and post)
	sum := 0
	for sum < 10 {
		sum += 2
	}
	fmt.Println("Sum:", sum)

	// infinite loop
	count := 0
	for {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Println("Broke at count:", count)

	// range over slice (like for-of in JS, or range-based for in C++11)
	nums := []int{1, 2, 3, 4, 5}
	for i, v := range nums {
		fmt.Printf("nums[%d] = %d\n", i, v)
	}

	// Skip index with _
	for _, v := range nums {
		fmt.Print(v, " ")
	}
	fmt.Println()

	// Go 1.22+: range over integer
	for i := range 5 {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// --- switch ---
	// No fall-through by default (unlike C)
	day := 2
	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3, 4, 5:
		fmt.Println("Midweek")
	default:
		fmt.Println("Weekend or invalid")
	}

	// switch with short statement
	switch n := 42; {
	case n < 0:
		fmt.Println("negative")
	case n == 0:
		fmt.Println("zero")
	default:
		fmt.Println("positive")
	}

	// type switch (we'll see more with interfaces)
	var val interface{} = "hello"
	switch v := val.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string:", v)
	default:
		fmt.Printf("other: %T\n", v)
	}

	// --- defer ---
	// Defers execution until surrounding function returns (LIFO order)
	// Like "finally" or cleanup; often for closing resources
	fmt.Print("Defer order: ")
	for i := 0; i < 3; i++ {
		defer fmt.Print(i, " ")
	}
	fmt.Println("(deferred)") // prints first; then 2 1 0 on return
}
