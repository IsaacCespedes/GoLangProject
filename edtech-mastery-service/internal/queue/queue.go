package queue

import (
	"context"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
)

// Queue is the postgres-backed outbox for async event processing.
type Queue struct {
	outbox *storage.OutboxRepo
}

func NewQueue(outbox *storage.OutboxRepo) *Queue {
	return &Queue{outbox: outbox}
}

func (q *Queue) Claim(ctx context.Context, limit int) ([]domain.OutboxItem, error) {
	return q.outbox.ClaimNext(ctx, limit)
}

func (q *Queue) MarkProcessed(ctx context.Context, outboxID int64) error {
	return q.outbox.MarkProcessed(ctx, outboxID)
}

func (q *Queue) MarkFailed(ctx context.Context, outboxID int64, errMsg string) error {
	return q.outbox.MarkFailed(ctx, outboxID, errMsg)
}
