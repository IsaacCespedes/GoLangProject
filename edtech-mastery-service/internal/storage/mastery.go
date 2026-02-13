package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

type MasteryRepo struct {
	pool *pgxpool.Pool
}

func NewMasteryRepo(pool *pgxpool.Pool) *MasteryRepo {
	return &MasteryRepo{pool: pool}
}

func (r *MasteryRepo) UpsertMastery(ctx context.Context, studentID, standardID string, score float64) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO student_mastery (student_id, standard_id, mastery_score, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (student_id, standard_id) DO UPDATE SET mastery_score = $3, updated_at = NOW()`,
		studentID, standardID, score,
	)
	return err
}

func (r *MasteryRepo) GetMasteryByStudent(ctx context.Context, studentID string) ([]domain.StandardMastery, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT standard_id, mastery_score FROM student_mastery WHERE student_id = $1 ORDER BY standard_id`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.StandardMastery
	for rows.Next() {
		var s domain.StandardMastery
		if err := rows.Scan(&s.StandardID, &s.MasteryScore); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
