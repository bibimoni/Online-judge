package ratingusecase

import (
	domain "contest/src/domain/entity"
	"context"
)

type RatingInteractor interface {
	ComputeContestRatings(ctx context.Context, input *ComputeContestRatingsInput) (*ComputeContestRatingsOutput, error)
	GetContestRatingResults(ctx context.Context, input *GetContestRatingResultsInput) (*GetContestRatingResultsOutput, error)
	GetUserRatingHistory(ctx context.Context, input *GetUserRatingHistoryInput) (*GetUserRatingHistoryOutput, error)
	GetUserRating(ctx context.Context, input *GetUserRatingInput) (*GetUserRatingOutput, error)
}

type ComputeContestRatingsInput struct {
	ContestId string `json:"contest_id" binding:"required"`
	Username  string
	UserRole  string
}

type ComputeContestRatingsOutput struct {
	Results []domain.RatingResult `json:"results"`
}

type GetContestRatingResultsInput struct {
	ContestId string `json:"contest_id" binding:"required"`
}

type GetContestRatingResultsOutput struct {
	Results []domain.RatingResult `json:"results"`
}

type GetUserRatingHistoryInput struct {
	Username string `json:"username" binding:"required"`
}

type GetUserRatingHistoryOutput struct {
	History []domain.RatingResult `json:"history"`
}

type GetUserRatingInput struct {
	Username string `json:"username" binding:"required"`
}

type GetUserRatingOutput struct {
	State *domain.UserRatingState `json:"state"`
}
