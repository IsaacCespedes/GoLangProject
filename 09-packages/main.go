// Lesson 09: Packages and modules
//
// Package = directory; module = go.mod root.
// Exported names: uppercase. Unexported: lowercase.

package main

import (
	"fmt"

	"go-tutorial/09-packages/mathutil"
)

func main() {
	fmt.Println("=== Lesson 09: Packages ===")

	sum := mathutil.Add(3, 4)
	fmt.Println("mathutil.Add(3,4):", sum)

	// mathutil.multiply(2,3) — compile error: unexported
}
