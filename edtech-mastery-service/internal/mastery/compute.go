package mastery

import (
	"context"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
)

// Default: score 0..100 maps to mastery 0..1
const maxScore = 100.0

type Service struct {
	mastery *storage.MasteryRepo
}

func NewService(mastery *storage.MasteryRepo) *Service {
	return &Service{mastery: mastery}
}

func (s *Service) UpdateFromGradedEvent(ctx context.Context, in *domain.IncomingEvent) error {
	if in.Score == nil || len(in.StandardIDs) == 0 {
		return nil
	}
	score := *in.Score
	if score < 0 {
		score = 0
	}
	if score > maxScore {
		score = maxScore
	}
	masteryScore := score / maxScore
	for _, std := range in.StandardIDs {
		if err := s.mastery.UpsertMastery(ctx, in.StudentID, std, masteryScore); err != nil {
			return err
		}
	}
	return nil
}
