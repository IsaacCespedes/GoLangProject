package domain

import "time"

type StudentMastery struct {
	StudentID   string    `json:"student_id"`
	StandardID  string    `json:"standard_id"`
	MasteryScore float64   `json:"mastery_score"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClassRollup struct {
	ClassID        string    `json:"class_id"`
	CompletionRate float64   `json:"completion_rate"`
	AvgScore       *float64  `json:"avg_score,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RiskFlag struct {
	StudentID  string    `json:"student_id"`
	ClassID    string    `json:"class_id"`
	Reason     string    `json:"reason"`
	ComputedAt time.Time `json:"computed_at"`
}

const (
	RiskReasonMissingSubmissions = "missing_submissions"
	RiskReasonScoreTrendDown     = "score_trend_down"
	RiskReasonBelowMedian        = "completion_below_median"
)
