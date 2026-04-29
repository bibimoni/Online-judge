package transporthealth

import (
	"context"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	helper "github.com/bibimoni/Online-judge/submission-judge/src/controller"
	"github.com/gin-gonic/gin"
)

func HandleHealth() gin.HandlerFunc {
	return common.InvokeUseCase(
		func(c *gin.Context) (input *struct{}, err error) {
			return
		},
		func(context context.Context, input *struct{}) (output *HealthReport, err error) {
			return &HealthReport{Status: "ok"}, nil
		},
		helper.WriteSuccessOutput,
	)
}

type HealthReport struct {
	Status string `json:"status"`
}
