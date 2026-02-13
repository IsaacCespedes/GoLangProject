package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
)

var ErrDuplicateEvent = errors.New("duplicate event (source, event_id)")

type EventRepo struct {
	pool *pgxpool.Pool
}

func NewEventRepo(pool *pgxpool.Pool) *EventRepo {
	return &EventRepo{pool: pool}
}

// InsertEvent inserts into events and outbox in one transaction. Idempotent: duplicate (source, event_id) does not create a new outbox row.
func (r *EventRepo) InsertEvent(ctx context.Context, eventID, source, eventType string, payload []byte) (eventDBID int64, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx,
		`INSERT INTO events (event_id, source, type, payload, created_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (source, event_id) DO NOTHING
		 RETURNING id`,
		eventID, source, eventType, payload,
	).Scan(&id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	if err == nil {
		// New row inserted; enqueue for worker
		_, err = tx.Exec(ctx,
			`INSERT INTO event_outbox (event_db_id, status) VALUES ($1, 'pending')`,
			id,
		)
		if err != nil {
			return 0, err
		}
		if err = tx.Commit(ctx); err != nil {
			return 0, err
		}
		return id, nil
	}

	// Duplicate: get existing id
	err = tx.QueryRow(ctx,
		`SELECT id FROM events WHERE source = $1 AND event_id = $2`,
		source, eventID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *EventRepo) GetEventByID(ctx context.Context, id int64) (*domain.Event, error) {
	var e domain.Event
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, source, type, payload, created_at FROM events WHERE id = $1`,
		id,
	).Scan(&e.ID, &e.EventID, &e.Source, &e.Type, &e.Payload, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}
