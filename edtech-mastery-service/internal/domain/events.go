package domain

import "time"

const (
	EventTypeAssignmentAssigned = "ASSIGNMENT_ASSIGNED"
	EventTypeSubmissionCreated  = "SUBMISSION_CREATED"
	EventTypeSubmissionGraded   = "SUBMISSION_GRADED"
)

// IncomingEvent is the API payload for POST /events
type IncomingEvent struct {
	EventID      string    `json:"event_id"`
	Source       string    `json:"source"`
	Timestamp    time.Time `json:"timestamp"`
	StudentID    string    `json:"student_id"`
	ClassID      string    `json:"class_id"`
	AssignmentID string    `json:"assignment_id"`
	StandardIDs  []string  `json:"standard_ids"`
	Score        *float64  `json:"score,omitempty"`
	RubricTags   []string  `json:"rubric_tags,omitempty"`
	Type         string    `json:"type,omitempty"` // optional: ASSIGNMENT_ASSIGNED, SUBMISSION_CREATED, SUBMISSION_GRADED
}

// Event is the stored event row (append-only)
type Event struct {
	ID        int64     `json:"id"`
	EventID   string    `json:"event_id"`
	Source    string    `json:"source"`
	Type      string    `json:"type"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// OutboxItem represents a pending event for the worker
type OutboxItem struct {
	ID         int64
	EventDBID  int64
	Status     string
	Attempts   int
	LastError  *string
	CreatedAt  time.Time
	ProcessedAt *time.Time
}
