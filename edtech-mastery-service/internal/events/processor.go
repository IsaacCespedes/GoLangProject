package events

import (
	"context"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/mastery"
	"github.com/edtech-mastery/student-progress-service/internal/risk"
	"github.com/edtech-mastery/student-progress-service/internal/rollups"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
)

type Processor struct {
	eventRepo *storage.EventRepo
	mastery   *mastery.Service
	rollups   *rollups.Service
	risk      *risk.Service
}

func NewProcessor(eventRepo *storage.EventRepo, mastery *mastery.Service, rollups *rollups.Service, risk *risk.Service) *Processor {
	return &Processor{
		eventRepo: eventRepo,
		mastery:   mastery,
		rollups:   rollups,
		risk:      risk,
	}
}

func (p *Processor) Process(ctx context.Context, eventDBID int64) error {
	event, err := p.eventRepo.GetEventByID(ctx, eventDBID)
	if err != nil || event == nil {
		return err
	}
	in, err := PayloadToIncoming(event.Payload)
	if err != nil {
		return err
	}
	classID := in.ClassID
	switch event.Type {
	case domain.EventTypeSubmissionGraded:
		if err := p.mastery.UpdateFromGradedEvent(ctx, in); err != nil {
			return err
		}
	}
	if classID != "" {
		if err := p.rollups.RecomputeForClass(ctx, classID); err != nil {
			return err
		}
		if err := p.risk.RecomputeForClass(ctx, classID); err != nil {
			return err
		}
	}
	return nil
}
