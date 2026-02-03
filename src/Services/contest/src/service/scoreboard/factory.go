package scoreboardservice

import (
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	contestsubmissionrepo "contest/src/domain/repository/contest-submission"
	scoreboardrepo "contest/src/domain/repository/scoreboard"
	scoreboardserviceimpl "contest/src/service/scoreboard/impl"
	"context"
)

// ScoreboardServiceFactory creates and manages scoreboard services
type ScoreboardServiceFactory struct {
	calculators           map[domain.ScoringType]ScoreboardService
	contestRepo           contestrepo.ContestRepository
	contestSubmissionRepo contestsubmissionrepo.ContestSubmissionRepository
	scoreboardRepo        scoreboardrepo.ScoreboardRepository
}

// NewScoreboardServiceFactory creates a factory with all calculators registered
func NewScoreboardServiceFactory(
	contestSubmissionRepo contestsubmissionrepo.ContestSubmissionRepository,
	contestRepo contestrepo.ContestRepository,
	scoreboardRepo scoreboardrepo.ScoreboardRepository,
) *ScoreboardServiceFactory {
	factory := &ScoreboardServiceFactory{
		calculators:           make(map[domain.ScoringType]ScoreboardService),
		contestRepo:           contestRepo,
		contestSubmissionRepo: contestSubmissionRepo,
		scoreboardRepo:        scoreboardRepo,
	}

	// Register all available calculators
	icpcCalculator := scoreboardserviceimpl.NewICPCCalculator(
		contestSubmissionRepo,
		contestRepo,
		scoreboardRepo,
	)
	factory.calculators[domain.ICPC] = icpcCalculator

	// Register IOI calculator
	ioiCalculator := scoreboardserviceimpl.NewIOICalculator(
		contestSubmissionRepo,
		contestRepo,
		scoreboardRepo,
	)
	factory.calculators[domain.IOI] = ioiCalculator

	return factory
}

// GetService returns a scoreboard service that automatically routes to the correct calculator
func (f *ScoreboardServiceFactory) NewScoreboardFactory() ScoreboardService {
	return &factoryScoreboardService{
		factory: f,
	}
}

// GetCalculator returns a specific calculator by scoring type
func (f *ScoreboardServiceFactory) GetCalculator(scoringType domain.ScoringType) (ScoreboardService, error) {
	calculator, exists := f.calculators[scoringType]
	if !exists {
		return nil, ErrNoCalculatorFound
	}
	return calculator, nil
}

// factoryScoreboardService is a wrapper that routes to the appropriate calculator
type factoryScoreboardService struct {
	factory *ScoreboardServiceFactory
}

// BuildScoreboardSnapshot automatically selects and uses the correct calculator
func (s *factoryScoreboardService) BuildScoreboardSnapshot(
	ctx context.Context,
	contestId string,
	includeVirtual bool,
	includeUnrated bool,
) (*domain.ScoreboardSnapshot, error) {
	// Get contest to determine scoring type
	contest, err := s.factory.contestRepo.GetById(ctx, contestId)
	if err != nil {
		return nil, err
	}

	// Get the appropriate calculator
	calculator, exists := s.factory.calculators[contest.ContestRule.ScoringType]
	if !exists {
		return nil, ErrNoCalculatorFound
	}

	// Delegate to the specific calculator
	return calculator.BuildScoreboardSnapshot(ctx, contestId, includeVirtual, includeUnrated)
}

// GetScoringType is not applicable for the factory service wrapper
func (s *factoryScoreboardService) GetScoringType() domain.ScoringType {
	return "" // This service routes to multiple calculators
}
