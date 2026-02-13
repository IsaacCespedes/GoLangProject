package rollups

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/storage"
)

type Service struct {
	pool   *pgxpool.Pool
	rollups *storage.RollupsRepo
}

func NewService(pool *pgxpool.Pool, rollups *storage.RollupsRepo) *Service {
	return &Service{pool: pool, rollups: rollups}
}

func (s *Service) RecomputeForClass(ctx context.Context, classID string) error {
	var completionRate float64
	var avgScore *float64
	err := s.pool.QueryRow(ctx, `
		WITH assigned AS (
			SELECT COUNT(DISTINCT (payload->>'student_id') || '|' || (payload->>'assignment_id')) AS cnt
			FROM events
			WHERE type = 'ASSIGNMENT_ASSIGNED' AND payload->>'class_id' = $1
		),
		graded AS (
			SELECT COUNT(DISTINCT (payload->>'student_id') || '|' || (payload->>'assignment_id')) AS cnt
			FROM events
			WHERE type = 'SUBMISSION_GRADED' AND payload->>'class_id' = $1
		),
		scores AS (
			SELECT AVG((payload->>'score')::float) AS avg_s
			FROM events
			WHERE type = 'SUBMISSION_GRADED' AND payload->>'class_id' = $1
		)
		SELECT COALESCE(g.cnt::float / NULLIF(a.cnt, 0), 0), s.avg_s
		FROM assigned a, graded g, scores s
	`, classID).Scan(&completionRate, &avgScore)
	if err != nil {
		return err
	}
	// Clamp to [0, 1] to satisfy DB check constraint (avoids float rounding > 1)
	completionRate = math.Max(0, math.Min(1, completionRate))
	return s.rollups.UpsertClassRollup(ctx, classID, completionRate, avgScore)
}
