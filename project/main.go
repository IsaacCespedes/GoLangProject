// Project: Simple CLI Tool — Task Manager
//
// Combines: structs, slices, JSON, file I/O, flags, error handling.
// Run: go run project/main.go [add|list|done] [args]

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
	NextID int   `json:"next_id"`
}

func (tl *TaskList) Add(title string) {
	tl.Tasks = append(tl.Tasks, Task{
		ID:    tl.NextID,
		Title: title,
		Done:  false,
	})
	tl.NextID++
}

func (tl *TaskList) MarkDone(id int) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].Done = true
			return true
		}
	}
	return false
}

func (tl *TaskList) Save(path string) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadTaskList(path string) (*TaskList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskList{NextID: 1}, nil
		}
		return nil, err
	}
	var tl TaskList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, err
	}
	if tl.NextID == 0 {
		tl.NextID = 1
	}
	return &tl, nil
}

func getDataPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".go-tutorial-tasks.json")
}

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	path := getDataPath()
	tl, err := LoadTaskList(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading tasks:", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		args := addCmd.Args()
		if len(args) == 0 {
			fmt.Println("Usage: go run main.go add <title>")
			os.Exit(1)
		}
		title := strings.Join(args, " ")
		tl.Add(title)
		if err := tl.Save(path); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving:", err)
			os.Exit(1)
		}
		fmt.Printf("Added: %s\n", title)

	case "list":
		listCmd.Parse(os.Args[2:])
		if len(tl.Tasks) == 0 {
			fmt.Println("No tasks.")
			return
		}
		for _, t := range tl.Tasks {
			status := " "
			if t.Done {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
		}

	case "done":
		doneCmd.Parse(os.Args[2:])
		args := doneCmd.Args()
		if len(args) == 0 {
			fmt.Println("Usage: go run main.go done <id>")
			os.Exit(1)
		}
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			fmt.Println("Invalid ID")
			os.Exit(1)
		}
		if tl.MarkDone(id) {
			tl.Save(path)
			fmt.Printf("Marked %d done.\n", id)
		} else {
			fmt.Println("Task not found")
			os.Exit(1)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Task Manager - Go Tutorial Project

Usage:
  go run main.go add <title>   Add a task
  go run main.go list          List all tasks
  go run main.go done <id>     Mark task as done`)
}
