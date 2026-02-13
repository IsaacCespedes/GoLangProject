package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

type OutboxRepo struct {
	pool *pgxpool.Pool
}

func NewOutboxRepo(pool *pgxpool.Pool) *OutboxRepo {
	return &OutboxRepo{pool: pool}
}

func (r *OutboxRepo) ClaimNext(ctx context.Context, limit int) ([]domain.OutboxItem, error) {
	rows, err := r.pool.Query(ctx,
		`UPDATE event_outbox SET status = 'processing', attempts = attempts + 1
		 WHERE id IN (
		   SELECT id FROM event_outbox WHERE status = 'pending' ORDER BY id LIMIT $1 FOR UPDATE SKIP LOCKED
		 )
		 RETURNING id, event_db_id, status, attempts, last_error, created_at, processed_at`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OutboxItem
	for rows.Next() {
		var o domain.OutboxItem
		err := rows.Scan(&o.ID, &o.EventDBID, &o.Status, &o.Attempts, &o.LastError, &o.CreatedAt, &o.ProcessedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	return items, rows.Err()
}

func (r *OutboxRepo) MarkProcessed(ctx context.Context, outboxID int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE event_outbox SET status = 'processed', processed_at = NOW() WHERE id = $1`,
		outboxID,
	)
	return err
}

func (r *OutboxRepo) MarkFailed(ctx context.Context, outboxID int64, errMsg string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE event_outbox SET status = 'pending', last_error = $2 WHERE id = $1`,
		outboxID, errMsg,
	)
	return err
}
