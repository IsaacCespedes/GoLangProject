// Lesson 04: Pointers and structs
//
// Pointers exist but are simpler than C: no pointer arithmetic.
// Structs are like C structs but can have methods (receiver functions).

package main

import (
	"fmt"
)

// Struct: aggregate type (like C struct, or a class without inheritance)
type Person struct {
	Name string
	Age  int
}

// Method: function with receiver (receiver goes before func name)
// Value receiver — operates on copy (read-only or small structs)
func (p Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s, %d years old.", p.Name, p.Age)
}

// Value receiver that modifies — only changes the copy, not the original
func (p Person) HaveBirthdayCopy() {
	p.Age++
}

// Pointer receiver — modifies the original; Go passes &p even when called on a value
func (p *Person) HaveBirthday() {
	p.Age++
}

// Constructor idiom (Go doesn't have constructors; use a function)
func NewPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

func main() {
	fmt.Println("=== Lesson 04: Pointers & Structs ===")

	// --- Pointers ---
	x := 42
	ptr := &x
	fmt.Println("x:", x, "ptr:", ptr, "*ptr:", *ptr)

	*ptr = 100
	fmt.Println("After *ptr=100, x:", x)

	// nil is the zero value for pointers (like nullptr in C++)
	var p *int
	fmt.Println("nil pointer:", p)
	// *p would panic

	// --- Struct literals ---
	p1 := Person{"Alice", 30}
	p2 := Person{Name: "Bob", Age: 25}
	p3 := Person{Name: "Charlie"} // Age gets zero value
	fmt.Println(p1, p2, p3)

	// Pointer to struct
	p4 := &Person{"Diana", 28}
	// Go automatically dereferences: p4.Name same as (*p4).Name
	fmt.Println(p4.Name, p4.Age)

	// --- Methods: value vs pointer receiver ---
	fmt.Println(p1.Greet())

	// Value receiver: receives a copy; modification doesn't affect original
	p1.HaveBirthdayCopy()
	fmt.Println("After HaveBirthdayCopy (value receiver):", p1.Age) // still 30

	// Pointer receiver: receives address; modification affects original
	// (Go passes &p1 automatically when calling on a value)
	p1.HaveBirthday()
	fmt.Println("After HaveBirthday (pointer receiver):", p1.Age) // now 31

	p4.HaveBirthday() // same: pointer receiver modifies p4
	fmt.Println("After HaveBirthday on p4 (pointer):", p4.Age)

	// --- Struct embedding (composition, not inheritance) ---
	type Employee struct {
		Person
		Title string
	}

	emp := Employee{
		Person: Person{"Eve", 35},
		Title:  "Engineer",
	}
	// Embedded fields promoted: emp.Name, emp.Age, emp.Greet() all work
	fmt.Println(emp.Name, emp.Title, emp.Greet())

	// --- new vs &Type{} ---
	// new(T) returns *T, zero-initialized
	np := new(Person)
	fmt.Println("new(Person):", np)

	// &Person{} is idiomatic when you need initial values
	np2 := &Person{Name: "Frank", Age: 40}
	fmt.Println(np2)
}
