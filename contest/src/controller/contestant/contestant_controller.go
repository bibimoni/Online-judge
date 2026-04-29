package contestantcontroller

import (
	"contest/src/common"
	"contest/src/controller"
	contestantusecase "contest/src/usecase/contestant"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Submit(contestantInteractor contestantusecase.ContestantInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toSubmitInput,
		contestantInteractor.Submit,
		common.WriteSuccessOutput[contestantusecase.SubmitOutput],
	)
}

func Register(contestantInteractor contestantusecase.ContestantInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toRegisterInput,
		contestantInteractor.Register,
		common.WriteSuccessOutput[contestantusecase.RegisterOutput],
	)
}

func Unregister(contestantInteractor contestantusecase.ContestantInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toUnregisterInput,
		contestantInteractor.Unregister,
		common.WriteSuccessOutput[contestantusecase.UnregisterOutput],
	)
}

func toRegisterInput(c *gin.Context) (*contestantusecase.RegisterInput, error) {
	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}

	var req contestantusecase.RegisterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	req.Username = rc.Username
	req.Role = rc.Role
	validator := validator.New()
	if err := validator.Struct(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func toUnregisterInput(c *gin.Context) (*contestantusecase.UnregisterInput, error) {
	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}

	var req contestantusecase.UnregisterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	req.Username = rc.Username
	req.Role = rc.Role
	validator := validator.New()
	if err := validator.Struct(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func toSubmitInput(c *gin.Context) (*contestantusecase.SubmitInput, error) {
	rc, ok := controller.GetContextRequest(c)
	if !ok {
		return nil, common.NewForbiddenError("unauthorized")
	}

	var req contestantusecase.SubmitInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	req.Username = rc.Username
	req.Role = rc.Role
	validator := validator.New()
	if err := validator.Struct(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
