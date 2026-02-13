// Simulate generates students, assignments, and random learning events for demo and load testing.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const defaultBaseURL = "http://localhost:8080"

func main() {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	nStudents := 20
	nAssignments := 10
	if s := os.Getenv("STUDENTS"); s != "" {
		if n, _ := strconv.Atoi(s); n > 0 {
			nStudents = n
		}
	}
	if s := os.Getenv("ASSIGNMENTS"); s != "" {
		if n, _ := strconv.Atoi(s); n > 0 {
			nAssignments = n
		}
	}
	classID := "class-1"
	teacherID := "teacher-1"
	standards := []string{"std-math-1", "std-math-2", "std-ela-1"}

	client := &http.Client{Timeout: 10 * time.Second}
	var assigned [][]string
	for a := 0; a < nAssignments; a++ {
		assignmentID := fmt.Sprintf("assign-%d", a+1)
		for s := 0; s < nStudents; s++ {
			studentID := fmt.Sprintf("student-%d", s+1)
			ev := map[string]interface{}{
				"event_id":      uuid.New().String(),
				"source":        "simulate",
				"timestamp":     time.Now().UTC().Format(time.RFC3339),
				"student_id":    studentID,
				"class_id":      classID,
				"assignment_id": assignmentID,
				"standard_ids":  standards,
				"type":          "ASSIGNMENT_ASSIGNED",
			}
			if err := postEvent(client, baseURL, ev); err != nil {
				fmt.Fprintf(os.Stderr, "assign err: %v\n", err)
				continue
			}
			assigned = append(assigned, []string{studentID, assignmentID})
		}
	}
	fmt.Printf("Posted %d ASSIGNMENT_ASSIGNED events\n", len(assigned))

	for _, pair := range assigned {
		studentID, assignmentID := pair[0], pair[1]
		ev := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"source":        "simulate",
			"timestamp":     time.Now().UTC().Format(time.RFC3339),
			"student_id":    studentID,
			"class_id":      classID,
			"assignment_id": assignmentID,
			"standard_ids":  standards,
			"type":          "SUBMISSION_CREATED",
		}
		if err := postEvent(client, baseURL, ev); err != nil {
			fmt.Fprintf(os.Stderr, "submission err: %v\n", err)
			continue
		}
		score := 60.0 + float64(time.Now().UnixNano()%4000)/100.0
		ev["score"] = score
		ev["event_id"] = uuid.New().String()
		ev["type"] = "SUBMISSION_GRADED"
		if err := postEvent(client, baseURL, ev); err != nil {
			fmt.Fprintf(os.Stderr, "graded err: %v\n", err)
			continue
		}
	}
	fmt.Printf("Posted SUBMISSION_CREATED + SUBMISSION_GRADED for %d student-assignment pairs\n", len(assigned))

	// Quick dashboard check
	resp, err := client.Get(baseURL + "/teachers/" + teacherID + "/classes/" + classID + "/dashboard")
	if err != nil {
		fmt.Fprintf(os.Stderr, "dashboard get: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Dashboard response status: %d\n", resp.StatusCode)
	fmt.Println("Simulate done. Run worker to process events, then GET dashboard again.")
}

func postEvent(client *http.Client, baseURL string, ev map[string]interface{}) error {
	body, _ := json.Marshal(ev)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/events", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}
