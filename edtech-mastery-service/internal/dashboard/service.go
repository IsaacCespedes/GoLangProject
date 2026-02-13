package dashboard

import (
	"context"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
)

type Service struct {
	rollups   *storage.RollupsRepo
	risk      *storage.RiskRepo
	recent    *storage.RecentActivityRepo
	mastery   *storage.MasteryRepo
	timeline  *storage.TimelineRepo
}

func NewService(rollups *storage.RollupsRepo, risk *storage.RiskRepo, recent *storage.RecentActivityRepo, mastery *storage.MasteryRepo, timeline *storage.TimelineRepo) *Service {
	return &Service{
		rollups:  rollups,
		risk:     risk,
		recent:   recent,
		mastery:  mastery,
		timeline: timeline,
	}
}

func (s *Service) TeacherDashboard(ctx context.Context, teacherID, classID string) (*domain.TeacherDashboard, error) {
	rollup, err := s.rollups.GetClassRollup(ctx, classID)
	if err != nil {
		return nil, err
	}
	atRisk, err := s.risk.GetAtRiskByClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	recent, err := s.recent.GetRecentByClass(ctx, classID, 20)
	if err != nil {
		return nil, err
	}
	return &domain.TeacherDashboard{
		ClassID:        classID,
		CompletionRate: rollup.CompletionRate,
		AverageScore:   rollup.AvgScore,
		AtRiskStudents: atRisk,
		RecentActivity: recent,
	}, nil
}

func (s *Service) StudentMastery(ctx context.Context, studentID string) (*domain.StudentMasteryView, error) {
	mastery, err := s.mastery.GetMasteryByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return &domain.StudentMasteryView{StudentID: studentID, Mastery: mastery}, nil
}

func (s *Service) StudentTimeline(ctx context.Context, studentID, classID string, limit int) (*domain.StudentTimeline, error) {
	events, err := s.timeline.GetRecentEventsForStudentClass(ctx, studentID, classID, limit)
	if err != nil {
		return nil, err
	}
	return &domain.StudentTimeline{StudentID: studentID, ClassID: classID, Events: events}, nil
}
