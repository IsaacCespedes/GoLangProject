// Lesson 08: Concurrency — goroutines and channels
//
// Goroutines: lightweight threads (not OS threads).
// Channels: typed conduit for sending/receiving (CSP model).

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Lesson 08: Concurrency ===")

	// --- Goroutines ---
	// Launch with "go" — runs concurrently
	go func() {
		fmt.Println("Hello from goroutine")
	}()

	time.Sleep(100 * time.Millisecond) // crude sync; channels are better

	// --- Channels ---
	ch := make(chan int)

	go func() {
		ch <- 42 // send
	}()

	v := <-ch // receive (blocks until value available)
	fmt.Println("Received:", v)

	// Buffered channel (non-blocking send until full)
	buf := make(chan int, 2)
	buf <- 1
	buf <- 2
	// buf <- 3 would block
	fmt.Println(<-buf, <-buf)

	// Close channel (sender closes; receiver can detect)
	ch2 := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for v := range ch2 {
		fmt.Print(v, " ")
	}
	fmt.Println("(channel closed)")

	// --- Select: multi-channel ---
	chA := make(chan string)
	chB := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		chA <- "from A"
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		chB <- "from B"
	}()

	select {
	case msg := <-chA:
		fmt.Println("Select got:", msg)
	case msg := <-chB:
		fmt.Println("Select got:", msg)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("Timeout")
	}

	// --- sync.WaitGroup ---
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("All workers done")

	// --- sync.Mutex (when channels aren't the right fit) ---
	var mu sync.Mutex
	counter := 0
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("Counter:", counter)
}
