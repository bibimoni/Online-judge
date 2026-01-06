package contestcontroller

import (
	"contest/src/common"
	"contest/src/infrastructure/config"
	"contest/src/usecase/contest"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Create(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toCreateContestInput,
		contestInteractor.CreateContest,
		common.WriteCreatedOutput[contestusecase.CreateContestOutput],
	)
}

func Edit(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toEditContestInput,
		contestInteractor.EditContest,
		common.WriteSuccessOutput[contestusecase.EditContestOutput],
	)
}

func Patch(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toPatchContestInput,
		contestInteractor.PatchContest,
		common.WriteSuccessOutput[contestusecase.PatchContestOutput],
	)
}

func toPatchContestInput(c *gin.Context) (*contestusecase.PatchContestInput, error) {
	contestId := c.Param("contest_id")
	var req contestusecase.PatchContestInput

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	req.ContestId = contestId
	req.UserRole = c.GetHeader("X-User-Role")

	config.GetLogger().Debug().Msgf("PatchContestInput: %+v", req)
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}

func toCreateContestInput(c *gin.Context) (*contestusecase.CreateContestInput, error) {
	// creator := c.DefaultQuery("creator", "")
	contestName := c.DefaultQuery("name", "")
	if contestName == "" {
		return nil, fmt.Errorf("contest creator and contest name is required")
	}

	var req struct {
		Username string `json:"username" validate:"min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.NewForbiddenError("No user found in request")
	}

	return &contestusecase.CreateContestInput{
		Username: req.Username,
		Name:     contestName,
		UserRole: c.GetHeader("X-User-Role"),
	}, nil
}

func toEditContestInput(c *gin.Context) (*contestusecase.EditContestInput, error) {
	editType := c.DefaultQuery("edit_type", "")
	if editType == "" {
		return nil, fmt.Errorf("edit_type is required")
	}

	var req contestusecase.EditContestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	req.EditType = contestusecase.EditType(editType)
	req.UserRole = c.GetHeader("X-User-Role")
	return &req, nil
}
