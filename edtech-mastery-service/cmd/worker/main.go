package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/events"
	"github.com/edtech-mastery/student-progress-service/internal/mastery"
	"github.com/edtech-mastery/student-progress-service/internal/queue"
	"github.com/edtech-mastery/student-progress-service/internal/risk"
	"github.com/edtech-mastery/student-progress-service/internal/rollups"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
	"github.com/edtech-mastery/student-progress-service/pkg/logging"
	"github.com/edtech-mastery/student-progress-service/pkg/metrics"
)

const (
	workerConcurrency = 4
	claimSize         = 10
	maxRetries        = 3
	pollInterval      = 2 * time.Second
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	log := logging.Init(env)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, getDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("db connect")
	}
	defer pool.Close()

	eventRepo := storage.NewEventRepo(pool)
	outboxRepo := storage.NewOutboxRepo(pool)
	masteryRepo := storage.NewMasteryRepo(pool)
	rollupsRepo := storage.NewRollupsRepo(pool)
	riskRepo := storage.NewRiskRepo(pool)

	masterySvc := mastery.NewService(masteryRepo)
	rollupsSvc := rollups.NewService(pool, rollupsRepo)
	riskSvc := risk.NewService(pool, riskRepo)
	processor := events.NewProcessor(eventRepo, masterySvc, rollupsSvc, riskSvc)

	q := queue.NewQueue(outboxRepo)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < workerConcurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runWorker(ctx, log, workerID, q, processor)
		}(i)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down workers")
	cancel()
	wg.Wait()
	log.Info().Msg("workers stopped")
}

func getDSN() string {
	if s := os.Getenv("DATABASE_URL"); s != "" {
		return s
	}
	return "postgres://postgres:postgres@localhost:5432/edtech?sslmode=disable"
}

func runWorker(ctx context.Context, log zerolog.Logger, workerID int, q *queue.Queue, processor *events.Processor) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		items, err := q.Claim(ctx, claimSize)
		if err != nil {
			log.Warn().Err(err).Int("worker", workerID).Msg("claim")
			time.Sleep(pollInterval)
			continue
		}
		if len(items) == 0 {
			time.Sleep(pollInterval)
			continue
		}
		for _, it := range items {
			processOne(ctx, log, workerID, q, processor, it)
		}
	}
}

func processOne(ctx context.Context, log zerolog.Logger, workerID int, q *queue.Queue, processor *events.Processor, item domain.OutboxItem) {
	start := time.Now()
	err := processor.Process(ctx, item.EventDBID)
	metrics.WorkerProcessingLatency.WithLabelValues("event").Observe(time.Since(start).Seconds())

	if err != nil {
		metrics.WorkerFailures.WithLabelValues("event").Inc()
		log.Warn().Err(err).Int64("outbox_id", item.ID).Int64("event_id", item.EventDBID).Msg("process failed")
		_ = q.MarkFailed(ctx, item.ID, err.Error())
		return
	}
	if err := q.MarkProcessed(ctx, item.ID); err != nil {
		log.Warn().Err(err).Int64("outbox_id", item.ID).Msg("mark processed failed")
	}
}
