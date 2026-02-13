package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

type RollupsRepo struct {
	pool *pgxpool.Pool
}

func NewRollupsRepo(pool *pgxpool.Pool) *RollupsRepo {
	return &RollupsRepo{pool: pool}
}

func (r *RollupsRepo) UpsertClassRollup(ctx context.Context, classID string, completionRate float64, avgScore *float64) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO class_rollups (class_id, completion_rate, avg_score, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (class_id) DO UPDATE SET completion_rate = $2, avg_score = $3, updated_at = NOW()`,
		classID, completionRate, avgScore,
	)
	return err
}

func (r *RollupsRepo) GetClassRollup(ctx context.Context, classID string) (*domain.ClassRollup, error) {
	var c domain.ClassRollup
	err := r.pool.QueryRow(ctx,
		`SELECT class_id, completion_rate, avg_score, updated_at FROM class_rollups WHERE class_id = $1`,
		classID,
	).Scan(&c.ClassID, &c.CompletionRate, &c.AvgScore, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &domain.ClassRollup{ClassID: classID, CompletionRate: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
