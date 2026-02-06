package scoreboardservice

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

type ScoreboardService interface {
	BuildScoreboardSnapshot(ctx context.Context, contestId string, includeVirtual bool, includeUnrated bool) (*domain.ScoreboardSnapshot, error)
}

var (
	ErrInvalidScoringType     = errors.New("invalid or unsupported scoring type")
	ErrNoCalculatorFound      = errors.New("no calculator found for scoring type")
	ErrScoreboardNotFinalized = errors.New("scoreboard has not been finalized")
)
