package contestcontroller

import (
	"contest/src/common"
	"contest/src/common/helper"
	"contest/src/usecase/contest"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Create(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toCreateContestInput,
		contestInteractor.CreateContest,
		helper.WriteCreatedOutput[contestusecase.CreateContestOutput],
	)
}

func Edit(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toEditContestInput,
		contestInteractor.EditContest,
		helper.WriteSuccessOutput[contestusecase.EditContestOutput],
	)
}

func toCreateContestInput(c *gin.Context) (*contestusecase.CreateContestInput, error) {
	creator := c.DefaultQuery("creator", "")
	if creator == "" {
		return nil, fmt.Errorf("creator is required")
	}
	return &contestusecase.CreateContestInput{
		Creator: creator,
	}, nil
}

func toEditContestInput(c *gin.Context) (*contestusecase.EditContestInput, error) {
	editType := c.DefaultQuery("edit-type", "")
	if editType == "" {
		return nil, fmt.Errorf("edit-type is required")
	}

	var req contestusecase.EditContestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	req.EditType = editType
	return &req, nil
}
