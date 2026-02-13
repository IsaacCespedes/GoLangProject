package risk

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
)

const (
	MissingSubmissionsThreshold = 3
	MaxRetries                 = 3
)

type Service struct {
	pool   *pgxpool.Pool
	riskRepo *storage.RiskRepo
}

func NewService(pool *pgxpool.Pool, riskRepo *storage.RiskRepo) *Service {
	return &Service{pool: pool, riskRepo: riskRepo}
}

// RecomputeForClass clears existing flags and recomputes at-risk students for the class.
func (s *Service) RecomputeForClass(ctx context.Context, classID string) error {
	if err := s.riskRepo.DeleteRiskFlagsForClass(ctx, classID); err != nil {
		return err
	}

	// Missing submissions: students with assigned but no graded submission count above threshold
	rows, err := s.pool.Query(ctx, `
		WITH assigned AS (
			SELECT DISTINCT (payload->>'student_id') AS student_id
			FROM events
			WHERE type = 'ASSIGNMENT_ASSIGNED' AND payload->>'class_id' = $1
		),
		graded AS (
			SELECT DISTINCT (payload->>'student_id') AS student_id
			FROM events
			WHERE type = 'SUBMISSION_GRADED' AND payload->>'class_id' = $1
		),
		missing AS (
			SELECT a.student_id FROM assigned a
			LEFT JOIN graded g ON a.student_id = g.student_id
			WHERE g.student_id IS NULL
		)
		SELECT student_id FROM missing
	`, classID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var studentID string
		if err := rows.Scan(&studentID); err != nil {
			return err
		}
		_ = s.riskRepo.UpsertRiskFlag(ctx, studentID, classID, domain.RiskReasonMissingSubmissions)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Completion below median: get class completion rate per student, then median, then flag below
	var medianRate sql.NullFloat64
	err = s.pool.QueryRow(ctx, `
		WITH student_rates AS (
			SELECT
				payload->>'student_id' AS student_id,
				COUNT(DISTINCT CASE WHEN type = 'SUBMISSION_GRADED' THEN payload->>'assignment_id' END)::float / NULLIF(COUNT(DISTINCT CASE WHEN type = 'ASSIGNMENT_ASSIGNED' THEN payload->>'assignment_id' END), 0) AS rate
			FROM events
			WHERE payload->>'class_id' = $1
			GROUP BY payload->>'student_id'
		),
		ordered AS (
			SELECT rate, row_number() OVER (ORDER BY rate) AS rn, count(*) OVER () AS cnt
			FROM student_rates WHERE rate IS NOT NULL
		)
		SELECT AVG(rate) FROM ordered WHERE rn IN (cnt/2, cnt/2 + 1)
	`, classID).Scan(&medianRate)
	if err == nil && medianRate.Valid {
		rows2, err := s.pool.Query(ctx, `
			WITH student_rates AS (
				SELECT
					payload->>'student_id' AS student_id,
					COUNT(DISTINCT CASE WHEN type = 'SUBMISSION_GRADED' THEN payload->>'assignment_id' END)::float / NULLIF(COUNT(DISTINCT CASE WHEN type = 'ASSIGNMENT_ASSIGNED' THEN payload->>'assignment_id' END), 0) AS rate
				FROM events WHERE payload->>'class_id' = $1
				GROUP BY payload->>'student_id'
			)
			SELECT student_id FROM student_rates WHERE rate IS NOT NULL AND rate < $2
		`, classID, medianRate.Float64)
		if err == nil {
			for rows2.Next() {
				var studentID string
				if err := rows2.Scan(&studentID); err != nil {
					rows2.Close()
					break
				}
				_ = s.riskRepo.UpsertRiskFlag(ctx, studentID, classID, domain.RiskReasonBelowMedian)
			}
			rows2.Close()
		}
	}

	// Score trend downward: last 2 graded scores - if latest < previous, flag
	rows3, err := s.pool.Query(ctx, `
		WITH graded AS (
			SELECT payload->>'student_id' AS student_id, (payload->>'score')::float AS score,
			       row_number() OVER (PARTITION BY payload->>'student_id' ORDER BY created_at DESC) AS rn
			FROM events
			WHERE type = 'SUBMISSION_GRADED' AND payload->>'class_id' = $1
		),
		pair AS (
			SELECT student_id, MAX(CASE WHEN rn = 1 THEN score END) AS last_score,
			       MAX(CASE WHEN rn = 2 THEN score END) AS prev_score
			FROM graded GROUP BY student_id
		)
		SELECT student_id FROM pair WHERE prev_score IS NOT NULL AND last_score < prev_score
	`, classID)
	if err != nil {
		return nil
	}
	defer rows3.Close()
	for rows3.Next() {
		var studentID string
		if err := rows3.Scan(&studentID); err != nil {
			break
		}
		_ = s.riskRepo.UpsertRiskFlag(ctx, studentID, classID, domain.RiskReasonScoreTrendDown)
	}

	return nil
}
