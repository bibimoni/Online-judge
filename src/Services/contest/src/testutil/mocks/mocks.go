package mocks

import (
	"context"
	"errors"

	"contest/src/domain/entity"
)

// MockContestRepository is a mock implementation of ContestRepository for testing
type MockContestRepository struct {
	CreateFunc              func(ctx context.Context, contest entity.Contest) (string, error)
	GetByIDFunc             func(ctx context.Context, id string) (*entity.Contest, error)
	UpdateMetaFunc          func(ctx context.Context, id string, patch map[string]interface{}) error
	UpdateProblemsFunc      func(ctx context.Context, id string, problems []entity.ContestProblem) error
	RegisterContestantFunc  func(ctx context.Context, contestID, username string) error
	SetRatedFunc            func(ctx context.Context, contestID string, rated bool, policy string) error
	SetScoringTypeFunc      func(ctx context.Context, contestID string, scoringType string, rules interface{}) error
	ListContestsFunc        func(ctx context.Context, filter map[string]interface{}) ([]entity.Contest, error)
	GetContestantsFunc      func(ctx context.Context, contestID string) ([]string, error)
	
	// For tracking calls
	CreateCalls             int
	GetByIDCalls            int
	UpdateMetaCalls         int
	UpdateProblemsCalls     int
	RegisterContestantCalls int
}

func NewMockContestRepository() *MockContestRepository {
	return &MockContestRepository{}
}

func (m *MockContestRepository) Create(ctx context.Context, contest entity.Contest) (string, error) {
	m.CreateCalls++
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, contest)
	}
	return contest.ID, nil
}

func (m *MockContestRepository) GetByID(ctx context.Context, id string) (*entity.Contest, error) {
	m.GetByIDCalls++
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *MockContestRepository) UpdateMeta(ctx context.Context, id string, patch map[string]interface{}) error {
	m.UpdateMetaCalls++
	if m.UpdateMetaFunc != nil {
		return m.UpdateMetaFunc(ctx, id, patch)
	}
	return nil
}

func (m *MockContestRepository) UpdateProblems(ctx context.Context, id string, problems []entity.ContestProblem) error {
	m.UpdateProblemsCalls++
	if m.UpdateProblemsFunc != nil {
		return m.UpdateProblemsFunc(ctx, id, problems)
	}
	return nil
}

func (m *MockContestRepository) RegisterContestant(ctx context.Context, contestID, username string) error {
	m.RegisterContestantCalls++
	if m.RegisterContestantFunc != nil {
		return m.RegisterContestantFunc(ctx, contestID, username)
	}
	return nil
}

func (m *MockContestRepository) SetRated(ctx context.Context, contestID string, rated bool, policy string) error {
	if m.SetRatedFunc != nil {
		return m.SetRatedFunc(ctx, contestID, rated, policy)
	}
	return nil
}

func (m *MockContestRepository) SetScoringType(ctx context.Context, contestID string, scoringType string, rules interface{}) error {
	if m.SetScoringTypeFunc != nil {
		return m.SetScoringTypeFunc(ctx, contestID, scoringType, rules)
	}
	return nil
}

func (m *MockContestRepository) ListContests(ctx context.Context, filter map[string]interface{}) ([]entity.Contest, error) {
	if m.ListContestsFunc != nil {
		return m.ListContestsFunc(ctx, filter)
	}
	return []entity.Contest{}, nil
}

func (m *MockContestRepository) GetContestants(ctx context.Context, contestID string) ([]string, error) {
	if m.GetContestantsFunc != nil {
		return m.GetContestantsFunc(ctx, contestID)
	}
	return []string{}, nil
}

// MockContestSubmissionRepository is a mock implementation for testing
type MockContestSubmissionRepository struct {
	UpsertFromJudgeEventFunc           func(ctx context.Context, submission entity.ContestSubmission) error
	ListByContestFunc                  func(ctx context.Context, contestID string, filter map[string]interface{}) ([]entity.ContestSubmission, error)
	SetIgnoredFunc                     func(ctx context.Context, contestID, submissionID string, ignored bool) error
	SetIgnoredByUserFunc               func(ctx context.Context, contestID, username string, ignored bool) error
	UpdateVerdictPointsForRejudgeFunc  func(ctx context.Context, contestID, submissionID string, verdict entity.Verdict, points float64) error
	
	UpsertCalls       int
	ListByContestCalls int
	SetIgnoredCalls   int
}

func NewMockContestSubmissionRepository() *MockContestSubmissionRepository {
	return &MockContestSubmissionRepository{}
}

func (m *MockContestSubmissionRepository) UpsertFromJudgeEvent(ctx context.Context, submission entity.ContestSubmission) error {
	m.UpsertCalls++
	if m.UpsertFromJudgeEventFunc != nil {
		return m.UpsertFromJudgeEventFunc(ctx, submission)
	}
	return nil
}

func (m *MockContestSubmissionRepository) ListByContest(ctx context.Context, contestID string, filter map[string]interface{}) ([]entity.ContestSubmission, error) {
	m.ListByContestCalls++
	if m.ListByContestFunc != nil {
		return m.ListByContestFunc(ctx, contestID, filter)
	}
	return []entity.ContestSubmission{}, nil
}

func (m *MockContestSubmissionRepository) SetIgnored(ctx context.Context, contestID, submissionID string, ignored bool) error {
	m.SetIgnoredCalls++
	if m.SetIgnoredFunc != nil {
		return m.SetIgnoredFunc(ctx, contestID, submissionID, ignored)
	}
	return nil
}

func (m *MockContestSubmissionRepository) SetIgnoredByUser(ctx context.Context, contestID, username string, ignored bool) error {
	if m.SetIgnoredByUserFunc != nil {
		return m.SetIgnoredByUserFunc(ctx, contestID, username, ignored)
	}
	return nil
}

func (m *MockContestSubmissionRepository) UpdateVerdictPointsForRejudge(ctx context.Context, contestID, submissionID string, verdict entity.Verdict, points float64) error {
	if m.UpdateVerdictPointsForRejudgeFunc != nil {
		return m.UpdateVerdictPointsForRejudgeFunc(ctx, contestID, submissionID, verdict, points)
	}
	return nil
}

// MockScoreboardCalculator is a mock implementation of the scoreboard calculator
type MockScoreboardCalculator struct {
	BuildScoreboardFunc func(contest entity.Contest, submissions []entity.ContestSubmission) (entity.ScoreboardSnapshot, error)
	BuildCalls          int
}

func NewMockScoreboardCalculator() *MockScoreboardCalculator {
	return &MockScoreboardCalculator{}
}

func (m *MockScoreboardCalculator) BuildScoreboard(contest entity.Contest, submissions []entity.ContestSubmission) (entity.ScoreboardSnapshot, error) {
	m.BuildCalls++
	if m.BuildScoreboardFunc != nil {
		return m.BuildScoreboardFunc(contest, submissions)
	}
	return entity.ScoreboardSnapshot{}, nil
}

// MockRatingPolicy is a mock implementation of rating policy
type MockRatingPolicy struct {
	ComputeFunc func(ctx context.Context, contest entity.Contest, standings entity.ScoreboardSnapshot, oldRatings map[string]int) ([]entity.RatingResult, error)
	ComputeCalls int
}

func NewMockRatingPolicy() *MockRatingPolicy {
	return &MockRatingPolicy{}
}

func (m *MockRatingPolicy) Compute(ctx context.Context, contest entity.Contest, standings entity.ScoreboardSnapshot, oldRatings map[string]int) ([]entity.RatingResult, error) {
	m.ComputeCalls++
	if m.ComputeFunc != nil {
		return m.ComputeFunc(ctx, contest, standings, oldRatings)
	}
	return []entity.RatingResult{}, nil
}
