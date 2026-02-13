// Lesson 10: Standard library — JSON, HTTP, files, time
//
// Go's stdlib is batteries-included. No npm needed for basics.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// JSON: struct tags for marshaling
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Println("=== Lesson 10: Standard Library ===")

	// --- encoding/json ---
	p := Person{Name: "Alice", Age: 30}
	bytes, _ := json.Marshal(p)
	fmt.Println("JSON:", string(bytes))

	var p2 Person
	json.Unmarshal(bytes, &p2)
	fmt.Println("Unmarshaled:", p2)

	// --- strings ---
	s := "  hello, world  "
	fmt.Println("TrimSpace:", strings.TrimSpace(s))
	fmt.Println("Contains:", strings.Contains(s, "world"))
	fmt.Println("Split:", strings.Split("a,b,c", ","))
	fmt.Println("Join:", strings.Join([]string{"a", "b"}, "-"))

	// --- path/filepath ---
	fmt.Println("Base:", filepath.Base("/usr/local/bin"))
	fmt.Println("Join:", filepath.Join("dir", "sub", "file.txt"))

	// --- time ---
	now := time.Now()
	fmt.Println("Now:", now.Format(time.RFC3339))
	fmt.Println("Unix:", now.Unix())

	dur := 2 * time.Second
	time.Sleep(100 * time.Millisecond) // truncated for demo
	fmt.Println("Duration:", dur)

	// --- os: files ---
	// Write
	tmp := os.TempDir()
	fpath := filepath.Join(tmp, "go-tutorial-demo.txt")
	os.WriteFile(fpath, []byte("Hello from Go\n"), 0644)
	fmt.Println("Wrote:", fpath)

	// Read
	data, err := os.ReadFile(fpath)
	if err != nil {
		fmt.Println("Read error:", err)
	} else {
		fmt.Println("Read:", strings.TrimSpace(string(data)))
	}
	os.Remove(fpath)

	// --- net/http (client) ---
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("HTTP status:", resp.Status)
	// For body: io.ReadAll(resp.Body)
}
