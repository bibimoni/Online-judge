package controller

import (
	"contest/src/common"
	"contest/src/common/helper"
	"contest/src/components"
	createcontest "contest/src/usecase/create_contest"
	editcontest "contest/src/usecase/edit_contest"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ContestController struct {
	appCtx components.AppContext
}

func NewContestController(appCtx components.AppContext) *ContestController {
	return &ContestController{
		appCtx: appCtx,
	}
}

func (cc *ContestController) Create() gin.HandlerFunc {
	return common.InvokeUseCase(
		cc.toCreateContestInput,
		cc.appCtx.GetCreateContestInteractor().CreateContest,
		helper.WriteCreatedOutput[createcontest.CreateContestOutput],
	)
}

func (cc *ContestController) Edit() gin.HandlerFunc {
	return common.InvokeUseCase(
		cc.toEditContestInput,
		cc.appCtx.GetEditContestInteractor().EditContest,
		helper.WriteSuccessOutput[editcontest.EditContestOutput],
	)
}

func (cc *ContestController) toCreateContestInput(c *gin.Context) (*createcontest.CreateContestInput, error) {
	authorIdStr := c.DefaultQuery("author-id", "0")
	authorId, err := strconv.ParseUint(authorIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	if authorId == 0 {
		return nil, fmt.Errorf("author-id is required")
	}
	return &createcontest.CreateContestInput{
		AuthorId: authorId,
	}, nil
}

func (cc *ContestController) toEditContestInput(c *gin.Context) (*editcontest.EditContestInput, error) {
	editType := c.DefaultQuery("edit-type", "")
	if editType == "" {
		return nil, fmt.Errorf("edit-type is required")
	}

	var req editcontest.EditContestInput
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
