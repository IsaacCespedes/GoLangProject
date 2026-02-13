package domain

import "time"

type TeacherDashboard struct {
	ClassID           string         `json:"class_id"`
	CompletionRate    float64        `json:"completion_rate"`
	AverageScore      *float64       `json:"average_score,omitempty"`
	AtRiskStudents    []AtRiskStudent `json:"at_risk_students"`
	RecentActivity    []RecentActivity `json:"recent_activity"`
}

type AtRiskStudent struct {
	StudentID string   `json:"student_id"`
	Reasons   []string `json:"reasons"`
}

type RecentActivity struct {
	EventType   string    `json:"event_type"`
	StudentID   string    `json:"student_id"`
	AssignmentID string   `json:"assignment_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type StudentMasteryView struct {
	StudentID string             `json:"student_id"`
	Mastery   []StandardMastery  `json:"mastery"`
}

type StandardMastery struct {
	StandardID   string  `json:"standard_id"`
	MasteryScore float64 `json:"mastery_score"`
}

type StudentTimeline struct {
	StudentID string          `json:"student_id"`
	ClassID   string          `json:"class_id"`
	Events    []TimelineEvent `json:"events"`
}

type TimelineEvent struct {
	EventType   string    `json:"event_type"`
	AssignmentID string   `json:"assignment_id"`
	Score       *float64  `json:"score,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
