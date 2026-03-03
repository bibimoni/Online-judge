package scoreboardcontroller

import (
	"contest/src/common"
	"contest/src/controller"
	scoreboardusecase "contest/src/usecase/scoreboard"

	"github.com/gin-gonic/gin"
)

func BuildScoreboardSnapshot(ScoreboardInteractor scoreboardusecase.ScoreboardInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toBuildScoreboardSnapshotInput,
		ScoreboardInteractor.BuildScoreboardSnapshot,
		common.WriteSuccessOutput[scoreboardusecase.BuildScoreboardSnapshotOutput],
	)
}

func toBuildScoreboardSnapshotInput(c *gin.Context) (*scoreboardusecase.BuildScoreboardSnapshotInput, error) {
	var input scoreboardusecase.BuildScoreboardSnapshotInput
	if err := c.BindJSON(&input); err != nil {
		return nil, err
	}
	rc, ok := controller.GetContextRequest(c)
	if ok {
		input.Username = rc.Username
		input.UserRole = rc.Role
		input.Authenticated = true
	} else {
		input.Authenticated = false
	}
	return &input, nil
}

func GetScoreboardAt(ScoreboardInteractor scoreboardusecase.ScoreboardInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetScoreboardAtInput,
		ScoreboardInteractor.GetScoreboardAt,
		common.WriteSuccessOutput[scoreboardusecase.BuildScoreboardSnapshotOutput],
	)
}

func toGetScoreboardAtInput(c *gin.Context) (*scoreboardusecase.GetScoreboardAtInput, error) {
	var input scoreboardusecase.GetScoreboardAtInput
	if err := c.BindJSON(&input); err != nil {
		return nil, err
	}
	rc, ok := controller.GetContextRequest(c)
	if ok {
		input.Username = rc.Username
		input.UserRole = rc.Role
		input.Authenticated = true
	} else {
		input.Authenticated = false
	}
	return &input, nil
}