// Lesson 05: Interfaces
//
// Go interfaces: implicit satisfaction (no "implements" keyword).
// "Accept interfaces, return structs." — design by capability, not hierarchy.

package main

import (
	"fmt"
)

// Interface: set of method signatures
type Speaker interface {
	Speak() string
}

type Writer interface {
	Write(data []byte) (int, error)
}

// Empty interface: any type satisfies it (like void* or any)
// Prefer generics (Go 1.18+) when possible
func describe(i interface{}) {
	fmt.Printf("Type: %T, value: %v\n", i, i)
}

// Dog implements Speaker (implicitly — no keyword needed)
type Dog struct{ Name string }

func (d Dog) Speak() string {
	return "Woof! I'm " + d.Name
}

// Cat implements Speaker
type Cat struct{ Name string }

func (c Cat) Speak() string {
	return "Meow! I'm " + c.Name
}

// Interface as parameter: accepts any type with Speak()
func makeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	fmt.Println("=== Lesson 05: Interfaces ===")

	dog := Dog{"Rex"}
	cat := Cat{"Whiskers"}

	// Both satisfy Speaker
	makeSpeak(dog)
	makeSpeak(cat)

	// Interface variable holds (value, type) pair
	var s Speaker
	s = dog
	fmt.Println("s (Dog):", s.Speak())

	s = cat
	fmt.Println("s (Cat):", s.Speak())

	// Type assertion: extract concrete type from interface
	var i interface{} = "hello"
	str, ok := i.(string)
	if ok {
		fmt.Println("String value:", str)
	}

	num, ok := i.(int)
	if !ok {
		fmt.Println("Not an int, ok =", ok, "num =", num)
	}

	// Type switch (seen in 02-control-flow)
	describe(42)
	describe("hello")
	describe(dog)

	// Nil interface vs nil concrete value
	var s2 Speaker
	fmt.Println("s2 == nil:", s2 == nil)

	var d *Dog = nil
	s2 = d
	fmt.Println("s2 (nil *Dog) == nil:", s2 == nil) // false! interface holds (nil, *Dog)
	// This is a common gotcha: check the concrete value if needed

	// Embedding interfaces
	type Animal interface {
		Speaker
	}

	// Any type that has Speak() satisfies Animal
	makeSpeak(dog) // works as Animal too
}
