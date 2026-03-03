package ratingrepo

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

// RatingRepository persists rating results and per-user Elo-MMR state.
type RatingRepository interface {
	// SaveRatingResults batch-inserts the rating results for a single contest.
	SaveRatingResults(ctx context.Context, results []domain.RatingResult) error

	// GetRatingResultsByContest returns all rating changes for a contest.
	GetRatingResultsByContest(ctx context.Context, contestId string) ([]domain.RatingResult, error)

	// GetRatingHistory returns the rating change history for a single user,
	// ordered newest-first.
	GetRatingHistory(ctx context.Context, username string) ([]domain.RatingResult, error)

	// SaveUserRatingStates upserts (by username) the per-user Elo-MMR state.
	SaveUserRatingStates(ctx context.Context, states []domain.UserRatingState) error

	// GetUserRatingState returns the current state for a single user.
	// Returns ErrUserRatingNotFound if the user has never been rated.
	GetUserRatingState(ctx context.Context, username string) (*domain.UserRatingState, error)

	// BatchGetUserRatingStates returns states for multiple users at once.
	// Users without a stored state are omitted from the returned map.
	BatchGetUserRatingStates(ctx context.Context, usernames []string) (map[string]domain.UserRatingState, error)
}

var (
	ErrUserRatingNotFound = errors.New("user rating state not found")
	ErrInvalidContestId   = errors.New("invalid contest id for rating lookup")
)
