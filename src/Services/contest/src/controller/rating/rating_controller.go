package ratingcontroller

import (
	"contest/src/common"
	"contest/src/controller"
	ratingusecase "contest/src/usecase/rating"

	"github.com/gin-gonic/gin"
)

// ComputeContestRatings triggers Elo-MMR computation for a finished contest.
// POST /contest/:contest_id/rating/compute
func ComputeContestRatings(interactor ratingusecase.RatingInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toComputeContestRatingsInput,
		interactor.ComputeContestRatings,
		common.WriteSuccessOutput[ratingusecase.ComputeContestRatingsOutput],
	)
}

func toComputeContestRatingsInput(c *gin.Context) (*ratingusecase.ComputeContestRatingsInput, error) {
	contestId := c.Param("contest_id")
	if contestId == "" {
		return nil, common.NewBadRequestError("contest_id is required")
	}

	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("authentication required")
	}

	return &ratingusecase.ComputeContestRatingsInput{
		ContestId: contestId,
		Username:  rc.Username,
		UserRole:  rc.Role,
	}, nil
}

// GetContestRatingResults returns all rating changes for a contest.
// GET /contest/:contest_id/rating
func GetContestRatingResults(interactor ratingusecase.RatingInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetContestRatingResultsInput,
		interactor.GetContestRatingResults,
		common.WriteSuccessOutput[ratingusecase.GetContestRatingResultsOutput],
	)
}

func toGetContestRatingResultsInput(c *gin.Context) (*ratingusecase.GetContestRatingResultsInput, error) {
	contestId := c.Param("contest_id")
	if contestId == "" {
		return nil, common.NewBadRequestError("contest_id is required")
	}
	return &ratingusecase.GetContestRatingResultsInput{
		ContestId: contestId,
	}, nil
}

// GetUserRatingHistory returns the full rating history for a user.
// GET /rating/user/:username/history
func GetUserRatingHistory(interactor ratingusecase.RatingInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetUserRatingHistoryInput,
		interactor.GetUserRatingHistory,
		common.WriteSuccessOutput[ratingusecase.GetUserRatingHistoryOutput],
	)
}

func toGetUserRatingHistoryInput(c *gin.Context) (*ratingusecase.GetUserRatingHistoryInput, error) {
	username := c.Param("username")
	if username == "" {
		return nil, common.NewBadRequestError("username is required")
	}
	return &ratingusecase.GetUserRatingHistoryInput{
		Username: username,
	}, nil
}

// GetUserRating returns the current rating state for a user.
// GET /rating/user/:username
func GetUserRating(interactor ratingusecase.RatingInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetUserRatingInput,
		interactor.GetUserRating,
		common.WriteSuccessOutput[ratingusecase.GetUserRatingOutput],
	)
}

func toGetUserRatingInput(c *gin.Context) (*ratingusecase.GetUserRatingInput, error) {
	username := c.Param("username")
	if username == "" {
		return nil, common.NewBadRequestError("username is required")
	}
	return &ratingusecase.GetUserRatingInput{
		Username: username,
	}, nil
}
