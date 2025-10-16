package transportsubmitsubmission

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	appctx "github.com/bibimoni/Online-judge/submission-judge/src/components"
	"github.com/bibimoni/Online-judge/submission-judge/src/controller"
	controller_utils "github.com/bibimoni/Online-judge/submission-judge/src/controller/utils"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/gin-gonic/gin"

	"fmt"

	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission"
)

func HandleSubmitSubmissionRequest(appContext appctx.AppContext) gin.HandlerFunc {
	submissionInteractor, err := controller_utils.InitInteractor(appContext)
	if err != nil {
		return nil
	}

	return common.InvokeUseCase(
		toSubmitSubmissionType,
		submissionInteractor.SubmitSubmission,
		helper.WriteCreatedOutput,
	)
}

func toSubmitSubmissionType(c *gin.Context) (*usecase.SubmitSubmissionInput, error) {
	log := config.GetLogger()
	var input usecase.SubmitSubmissionInput
	if err := c.BindJSON(&input); err != nil {
		log.Error().Msgf("%s", err.Error())
		return nil, fmt.Errorf("Invalid Request Body")
	}

	// Guard submission type, i think all the validation should happen here
	// as long as it doesn't require any service / repository
	if input.SubmissionType != domain.SubmissionType(domain.ICPC) {
		return nil, fmt.Errorf("Sorry, we currently don't support this type of submission: %s", input.SubmissionType)
	}

	return &input, nil
}
