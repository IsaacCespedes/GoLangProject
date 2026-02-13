package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

type RecentActivityRepo struct {
	pool *pgxpool.Pool
}

func NewRecentActivityRepo(pool *pgxpool.Pool) *RecentActivityRepo {
	return &RecentActivityRepo{pool: pool}
}

func (r *RecentActivityRepo) GetRecentByClass(ctx context.Context, classID string, limit int) ([]domain.RecentActivity, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx,
		`SELECT e.type, e.payload->>'student_id', e.payload->>'assignment_id', e.created_at
		 FROM events e
		 WHERE e.payload->>'class_id' = $1
		 ORDER BY e.created_at DESC LIMIT $2`,
		classID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.RecentActivity
	for rows.Next() {
		var a domain.RecentActivity
		var studentID, assignmentID *string
		if err := rows.Scan(&a.EventType, &studentID, &assignmentID, &a.CreatedAt); err != nil {
			return nil, err
		}
		if studentID != nil {
			a.StudentID = *studentID
		}
		if assignmentID != nil {
			a.AssignmentID = *assignmentID
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
