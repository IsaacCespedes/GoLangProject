// Lesson 07: Error handling
//
// Go has no exceptions. Errors are values — return them explicitly.
// Check errors at call sites. "if err != nil" is idiomatic.

package main

import (
	"errors"
	"fmt"
)

// Custom error
var ErrNotFound = errors.New("not found")

// Sentinel errors for comparison
var ErrInvalidInput = errors.New("invalid input")

// Function that can fail
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Return sentinel for specific cases
func find(id int) (string, error) {
	if id < 0 {
		return "low", ErrInvalidInput
	}
	if id > 10 {
		return "high", ErrNotFound
	}
	return fmt.Sprintf("item-%d", id), nil
}

// Custom error type (wrap context)
type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %v", e.Field, e.Err)
}

func validate(name string) error {
	if name == "" {
		return &ValidationError{Field: "name", Err: ErrInvalidInput}
	}
	return nil
}

func main() {
	fmt.Println("=== Lesson 07: Error Handling ===")

	// Basic pattern
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("10/2 =", result)

	// Fail case
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Expected:", err)
	}

	// Sentinel check
	val, err := find(99)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Handle not found:", val)
	}

	_, err = find(-1)
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("Handle invalid input")
	}

	// Type assertion for custom errors
	err = validate("")
	if err != nil {
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Println("Validation error:", valErr.Field, valErr.Err)
		}
	}

	// errors.Join (Go 1.20+): combine errors
	err1 := errors.New("first")
	err2 := errors.New("second")
	combined := errors.Join(err1, err2)
	fmt.Println("Combined:", combined)

	// Panic/recover: avoid for normal flow; use for programmer errors
	// defer func() {
	//   if r := recover(); r != nil {
	//     fmt.Println("Recovered:", r)
	//   }
	// }()
	// panic("something went wrong")
}
