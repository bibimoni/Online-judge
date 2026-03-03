package ratingservice

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

type RatingService interface {
	ComputeRatings(
		ctx context.Context,
		contest *domain.Contest,
		standings []domain.StandingRow,
		oldStates map[string]domain.UserRatingState,
	) ([]domain.RatingResult, map[string]domain.UserRatingState, error)
}

var (
	ErrContestNotEnded     = errors.New("contest has not ended yet")
	ErrNoStandings         = errors.New("no standings to compute ratings for")
	ErrAlreadyRated        = errors.New("ratings have already been computed for this contest")
)
