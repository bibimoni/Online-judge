package scoreboardservice

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

// ScoreboardCalculator is the interface for calculating scoreboards
type ScoreboardService interface {
	// BuildScoreboardSnapshot computes the scoreboard from contest and submissions
	BuildScoreboardSnapshot(ctx context.Context, contestId string, includeVirtual bool, includeUnrated bool) (*domain.ScoreboardSnapshot, error)

	// GetScoringType returns the scoring type this calculator handles
	GetScoringType() domain.ScoringType
}

// ScoreboardService orchestrates scoreboard operations
// type ScoreboardService interface {
// 	// BuildScoreboard calculates and returns a scoreboard snapshot
// 	BuildScoreboard(ctx context.Context, contest *domain.Contest, submissions []domain.ContestSubmission) (*domain.ScoreboardSnapshot, error)

// 	// GetLiveScoreboard retrieves the current live scoreboard
// 	GetLiveScoreboard(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error)

// 	// SaveLiveScoreboard saves the live scoreboard to Redis
// 	SaveLiveScoreboard(ctx context.Context, contestId string, snapshot *domain.ScoreboardSnapshot) error

// 	// FinalizeScoreboard creates a final snapshot and saves it to MongoDB
// 	FinalizeScoreboard(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error)

// 	// GetFinalScoreboard retrieves the final scoreboard
// 	GetFinalScoreboard(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error)

// 	// RefreshLiveScoreboard recalculates and updates live scoreboard
// 	RefreshLiveScoreboard(ctx context.Context, contestId string) error
// }

var (
	ErrInvalidScoringType    = errors.New("invalid or unsupported scoring type")
	ErrNoCalculatorFound     = errors.New("no calculator found for scoring type")
	ErrScoreboardNotFinalized = errors.New("scoreboard has not been finalized")
)
