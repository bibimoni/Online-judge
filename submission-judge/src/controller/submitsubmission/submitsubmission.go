package transportsubmitsubmission

import (
	"slices"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	helper "github.com/bibimoni/Online-judge/submission-judge/src/controller"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fmt"

	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission/interactor"
)

func HandleInternalContestSubmitSubmissionRequest(submissioninteractor *interactor.SubmissionInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toInternalContestSubmitSubmissionType,
		submissioninteractor.InternalContestSubmitSubmission,
		helper.WriteCreatedOutput,
	)
}

func HandleSubmitSubmissionRequest(submissioninteractor *interactor.SubmissionInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toSubmitSubmissionType,
		submissioninteractor.SubmitSubmission,
		helper.WriteCreatedOutput,
	)
}

func HandleRejudgeSubmissionRequest(submissioninteractor *interactor.SubmissionInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toRejudgeSubmissionType,
		submissioninteractor.RejudgeSubmission,
		helper.WriteCreatedOutput,
	)
}

func toInternalContestSubmitSubmissionType(c *gin.Context) (*usecase.InternalContestSubmitSubmissionInput, error) {
	if !checkInternal(c) {
		log.Error().Msgf("invalid internal secret")
		return nil, fmt.Errorf("forbidden")
	}

	var input usecase.InternalContestSubmitSubmissionInput
	if err := c.BindJSON(&input); err != nil {
		log.Error().Msgf("%s", err.Error())
		return nil, fmt.Errorf("invalid Request Body")
	}

	return &input, nil
}

func toRejudgeSubmissionType(c *gin.Context) (*usecase.RejudgeSubmissionInput, error) {
	if !checkInternal(c) {
		log.Error().Msgf("invalid internal secret")
		return nil, fmt.Errorf("forbidden")
	}

	log := config.GetLogger()
	var input usecase.RejudgeSubmissionInput

	if err := c.BindJSON(&input); err != nil {
		log.Error().Msgf("%s", err.Error())
		return nil, fmt.Errorf("invalid Request Body")
	}

	// Make sure submission ids are unique
	slices.Sort(input.SubmissionIds)
	input.SubmissionIds = slices.Compact(input.SubmissionIds)
	return &input, nil
}

func toSubmitSubmissionType(c *gin.Context) (*usecase.SubmitSubmissionInput, error) {
	log := config.GetLogger()
	var input usecase.SubmitSubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Msgf("%s", err.Error())
		return nil, fmt.Errorf("invalid Request Body")
	}
	input.Username = c.GetHeader("X-Username")
	input.Role = common.RoleName(c.GetHeader("X-User-Role"))

	// Guard submission type, i think all the validation should happen here
	// as long as it doesn't require any service / repository
	if input.SubmissionType != domain.SubmissionType(domain.ICPC) {
		return nil, fmt.Errorf("sorry, we currently don't support thi type of submission: %s", input.SubmissionType)
	}

	return &input, nil
}

func checkInternal(c *gin.Context) bool {
	secretHeader := c.GetHeader("X-Internal-Secret")
	cfg, err := config.Load()
	if err != nil {
		log.Panic().Err(err).Msg("failed to get config")
		return false
	}

	return secretHeader == cfg.InternalSecret
}
