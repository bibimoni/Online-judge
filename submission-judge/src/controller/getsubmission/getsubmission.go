package transportgetsubmission

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	helper "github.com/bibimoni/Online-judge/submission-judge/src/controller"
	"github.com/gin-gonic/gin"

	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission/interactor"
)

func HandleGetSubmissionRequest(submissioninteractor *interactor.SubmissionInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetSubmissionType,
		submissioninteractor.GetSubmission,
		helper.WriteSuccessOutput,
	)
}

func HandleGetProblemSubmissionRequest(submissioninteractor *interactor.SubmissionInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetProblemSubmissionType,
		submissioninteractor.GetProblemSubmission,
		helper.WriteSuccessOutput,
	)
}

func toGetProblemSubmissionType(c *gin.Context) (*usecase.GetProblemSubmissionInput, error) {
	pid := c.Param("problem_id")
	log := config.GetLogger()
	log.Debug().Msgf("get problem id from request: %s", pid)
	return &usecase.GetProblemSubmissionInput{
		ProblemId: pid,
	}, nil
}

func toGetSubmissionType(c *gin.Context) (*usecase.GetSubmissionInput, error) {
	sid := c.Param("submission_id")
	log := config.GetLogger()
	log.Debug().Msgf("get submission id from request: %s", sid)
	return &usecase.GetSubmissionInput{SubmissionId: sid}, nil
}
