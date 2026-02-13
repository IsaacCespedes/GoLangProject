package events

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

var (
	ErrInvalidEventType = errors.New("invalid event type")
	ErrMissingFields    = errors.New("missing required fields")
)

var validTypes = map[string]bool{
	domain.EventTypeAssignmentAssigned: true,
	domain.EventTypeSubmissionCreated:  true,
	domain.EventTypeSubmissionGraded:   true,
}

func ValidateAndSetType(in *domain.IncomingEvent) (eventType string, err error) {
	if in.EventID == "" || in.Source == "" || in.StudentID == "" || in.ClassID == "" {
		return "", fmt.Errorf("%w: event_id, source, student_id, class_id required", ErrMissingFields)
	}
	if in.Type != "" && validTypes[in.Type] {
		// Client-provided type must be consistent (e.g. GRADED must have score)
		if in.Type == domain.EventTypeSubmissionGraded && in.Score == nil {
			return "", fmt.Errorf("%w: SUBMISSION_GRADED requires score", ErrInvalidEventType)
		}
		return in.Type, nil
	}
	switch {
	case len(in.StandardIDs) > 0 && in.Score != nil:
		eventType = domain.EventTypeSubmissionGraded
	case in.AssignmentID != "" && in.Score == nil && len(in.StandardIDs) > 0:
		eventType = domain.EventTypeSubmissionCreated
	case in.AssignmentID != "" && len(in.StandardIDs) > 0:
		eventType = domain.EventTypeAssignmentAssigned
	default:
		return "", fmt.Errorf("%w: cannot infer type from payload", ErrInvalidEventType)
	}
	return eventType, nil
}

func PayloadFromIncoming(in *domain.IncomingEvent) ([]byte, error) {
	return json.Marshal(in)
}
