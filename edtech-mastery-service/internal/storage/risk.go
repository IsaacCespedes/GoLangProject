package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

type RiskRepo struct {
	pool *pgxpool.Pool
}

func NewRiskRepo(pool *pgxpool.Pool) *RiskRepo {
	return &RiskRepo{pool: pool}
}

func (r *RiskRepo) UpsertRiskFlag(ctx context.Context, studentID, classID, reason string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO risk_flags (student_id, class_id, reason, computed_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (student_id, class_id, reason) DO UPDATE SET computed_at = NOW()`,
		studentID, classID, reason,
	)
	return err
}

func (r *RiskRepo) DeleteRiskFlagsForClass(ctx context.Context, classID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM risk_flags WHERE class_id = $1`, classID)
	return err
}

func (r *RiskRepo) GetAtRiskByClass(ctx context.Context, classID string) ([]domain.AtRiskStudent, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT student_id, reason FROM risk_flags WHERE class_id = $1 ORDER BY student_id`,
		classID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byStudent := make(map[string][]string)
	for rows.Next() {
		var studentID, reason string
		if err := rows.Scan(&studentID, &reason); err != nil {
			return nil, err
		}
		byStudent[studentID] = append(byStudent[studentID], reason)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var out []domain.AtRiskStudent
	for id, reasons := range byStudent {
		out = append(out, domain.AtRiskStudent{StudentID: id, Reasons: reasons})
	}
	return out, nil
}
