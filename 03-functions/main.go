// Lesson 03: Functions — multiple returns, variadic, first-class
//
// Go functions: explicit returns, no default args, no overloading.
// First-class: can be assigned, passed, returned.

package main

import (
	"fmt"
)

// Basic function: func name(params) returnType
func add(a int, b int) int {
	return a + b
}

// Same types: can shorten param list
func add3(a, b, c int) int {
	return a + b + c
}

// Multiple return values (common in Go — used for errors)
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Named return values: variables declared in signature
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // "naked" return — returns x, y
}

// Variadic function: accepts zero or more args of same type
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Variadic can be mixed with regular params (variadic must be last)
func join(sep string, parts ...string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

func main() {
	fmt.Println("=== Lesson 03: Functions ===")

	fmt.Println("add(2,3):", add(2, 3))
	fmt.Println("add3(1,2,3):", add3(1, 2, 3))

	// Multiple returns
	q, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10/2 =", q)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Expected error:", err)
	}

	// Named returns
	a, b := split(17)
	fmt.Println("split(17):", a, b)

	// Variadic
	fmt.Println("sum(1,2,3,4,5):", sum(1, 2, 3, 4, 5))
	fmt.Println("sum():", sum())
	fmt.Println("join:", join("-", "a", "b", "c"))

	// Spread slice into variadic
	nums := []int{1, 2, 3}
	fmt.Println("sum(nums...):", sum(nums...))

	// First-class functions
	double := func(x int) int { return x * 2 }
	fmt.Println("double(5):", double(5))

	// Closures capture by reference (like JS)
	adder := func() func(int) int {
		total := 0
		return func(x int) int {
			total += x
			return total
		}
	}
	acc := adder()
	fmt.Println(acc(1), acc(2), acc(3)) // 1 3 6

	// Higher-order: function as argument
	apply := func(f func(int) int, x int) int {
		return f(x)
	}
	fmt.Println("apply(double, 7):", apply(double, 7))
}
