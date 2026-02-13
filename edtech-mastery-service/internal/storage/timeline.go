package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

type TimelineRepo struct {
	pool *pgxpool.Pool
}

func NewTimelineRepo(pool *pgxpool.Pool) *TimelineRepo {
	return &TimelineRepo{pool: pool}
}

func (r *TimelineRepo) GetRecentEventsForStudentClass(ctx context.Context, studentID, classID string, limit int) ([]domain.TimelineEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx,
		`SELECT e.type, (e.payload->>'assignment_id')::text, (e.payload->>'score')::float, e.created_at
		 FROM events e
		 WHERE e.payload->>'student_id' = $1 AND e.payload->>'class_id' = $2
		 ORDER BY e.created_at DESC LIMIT $3`,
		studentID, classID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.TimelineEvent
	for rows.Next() {
		var t domain.TimelineEvent
		var score *float64
		var assignmentID *string
		if err := rows.Scan(&t.EventType, &assignmentID, &score, &t.CreatedAt); err != nil {
			return nil, err
		}
		if assignmentID != nil {
			t.AssignmentID = *assignmentID
		}
		t.Score = score
		out = append(out, t)
	}
	return out, rows.Err()
}
