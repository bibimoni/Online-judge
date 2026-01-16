package contestcontroller

import (
	"contest/src/common"
	"contest/src/controller"
	"contest/src/infrastructure/config"
	"contest/src/usecase/contest"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetAllContest(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetAllContestInput,
		contestInteractor.GetAllContests,
		common.WriteSuccessOutput[contestusecase.GetContestsOutput],
	)
}

func GetContest(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetContestInput,
		contestInteractor.GetContestById,
		common.WriteSuccessOutput[contestusecase.GetContestByIdOutput],
	)
}

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

func ManageProblems(contestInteractor contestusecase.ContestInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toManageProblemsInput,
		contestInteractor.ManageContestProblems,
		common.WriteSuccessOutput[contestusecase.ManageContestProblemsOutput],
	)
}

func toGetAllContestInput(c *gin.Context) (*contestusecase.GetContestsInput, error) {
	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return &contestusecase.GetContestsInput{
			Authenticated: false,
		}, nil
	}
	return &contestusecase.GetContestsInput{
		Username:      rc.Username,
		Role:          rc.Role,
		Authenticated: true,
	}, nil
}

func toGetContestInput(c *gin.Context) (*contestusecase.GetContestByIdInput, error) {
	contestId := c.Param("contest_id")
	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return &contestusecase.GetContestByIdInput{
			ContestId: contestId,
			GetContestsInput: contestusecase.GetContestsInput{
				Authenticated: false,
			},
		}, nil
	}
	return &contestusecase.GetContestByIdInput{
		ContestId: contestId,
		GetContestsInput: contestusecase.GetContestsInput{
			Username:      rc.Username,
			Role:          rc.Role,
			Authenticated: true,
		},
	}, nil
}

func toManageProblemsInput(c *gin.Context) (*contestusecase.ManageContestProblemsInput, error) {
	contestId := c.Param("contest_id")
	var req contestusecase.ManageContestProblemsInput

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}
	req.ContestId = contestId
	req.Username = rc.Username

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}

func toPatchContestInput(c *gin.Context) (*contestusecase.PatchContestInput, error) {
	contestId := c.Param("contest_id")
	var req contestusecase.PatchContestInput

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}
	req.ContestId = contestId
	req.UserRole = rc.Role
	req.Username = rc.Username

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

	// var req struct {
	// 	Username string `json:"username" validate:"min=6"`
	// }
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	return nil, common.NewForbiddenError("No user found in request")
	// }

	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}
	return &contestusecase.CreateContestInput{
		Username: rc.Username,
		Name:     contestName,
		UserRole: rc.Role,
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

	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}
	req.EditType = contestusecase.EditType(editType)
	req.UserRole = rc.Role
	req.Username = rc.Username
	return &req, nil
}
