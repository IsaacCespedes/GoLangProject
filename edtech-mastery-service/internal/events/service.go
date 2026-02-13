package events

import (
	"context"
	"encoding/json"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
	"github.com/edtech-mastery/student-progress-service/pkg/metrics"
)

type Service struct {
	eventRepo *storage.EventRepo
}

func NewService(eventRepo *storage.EventRepo) *Service {
	return &Service{eventRepo: eventRepo}
}

func (s *Service) Ingest(ctx context.Context, in *domain.IncomingEvent) (eventType string, eventDBID int64, err error) {
	eventType, err = ValidateAndSetType(in)
	if err != nil {
		metrics.EventsIngested.WithLabelValues("unknown", "validation_error").Inc()
		return "", 0, err
	}
	payload, err := PayloadFromIncoming(in)
	if err != nil {
		metrics.EventsIngested.WithLabelValues(eventType, "error").Inc()
		return "", 0, err
	}
	eventDBID, err = s.eventRepo.InsertEvent(ctx, in.EventID, in.Source, eventType, payload)
	if err != nil {
		metrics.EventsIngested.WithLabelValues(eventType, "error").Inc()
		return "", 0, err
	}
	metrics.EventsIngested.WithLabelValues(eventType, "ok").Inc()
	return eventType, eventDBID, nil
}

func PayloadToIncoming(payload []byte) (*domain.IncomingEvent, error) {
	var in domain.IncomingEvent
	if err := json.Unmarshal(payload, &in); err != nil {
		return nil, err
	}
	return &in, nil
}
