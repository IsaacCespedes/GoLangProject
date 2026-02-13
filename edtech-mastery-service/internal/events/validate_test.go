package events

import (
	"testing"
	"time"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

func TestValidateAndSetType(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name    string
		in      domain.IncomingEvent
		want    string
		wantErr bool
	}{
		{
			name: "inferred submission_created (assignment + standards, no score)",
			in: domain.IncomingEvent{
				EventID: "e1", Source: "s1", Timestamp: now,
				StudentID: "st1", ClassID: "c1", AssignmentID: "a1", StandardIDs: []string{"std1"},
			},
			want:    domain.EventTypeSubmissionCreated,
			wantErr: false,
		},
		{
			name: "explicit assignment_assigned via type field",
			in: domain.IncomingEvent{
				EventID: "e0", Source: "s1", Timestamp: now,
				StudentID: "st1", ClassID: "c1", AssignmentID: "a1", StandardIDs: []string{"std1"},
				Type: domain.EventTypeAssignmentAssigned,
			},
			want:    domain.EventTypeAssignmentAssigned,
			wantErr: false,
		},
		{
			name: "submission_graded",
			in: domain.IncomingEvent{
				EventID: "e2", Source: "s1", Timestamp: now,
				StudentID: "st1", ClassID: "c1", AssignmentID: "a1", StandardIDs: []string{"std1"},
				Score: ptrFloat64(85.5),
			},
			want:    domain.EventTypeSubmissionGraded,
			wantErr: false,
		},
		{
			name: "missing_event_id",
			in: domain.IncomingEvent{
				Source: "s1", Timestamp: now,
				StudentID: "st1", ClassID: "c1", AssignmentID: "a1", StandardIDs: []string{"std1"},
			},
			wantErr: true,
		},
		{
			name: "missing_student_id",
			in: domain.IncomingEvent{
				EventID: "e1", Source: "s1", Timestamp: now,
				ClassID: "c1", AssignmentID: "a1", StandardIDs: []string{"std1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndSetType(&tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndSetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ValidateAndSetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptrFloat64(f float64) *float64 { return &f }
