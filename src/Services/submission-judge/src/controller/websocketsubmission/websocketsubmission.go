package websocketsubmission

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission/interactor"
	"github.com/gin-gonic/gin"
)

func HandleSubmissionWSRequest(wsSubmissionInteractor *interactor.WSSubmissionInteractor) gin.HandlerFunc {
	return common.InvokeWSUseCase(
		toSubmissionWSRequest,
		wsSubmissionInteractor.SubmissionStatus,
	)
}

func toSubmissionWSRequest(c *gin.Context) (*usecase.WSSubmissionInput, error) {
	username := c.DefaultQuery("username", "*")
	problemId := c.DefaultQuery("problem_id", "*")
	submissionId := c.DefaultQuery("submission_id", "*")

	return &usecase.WSSubmissionInput{
		Username:     username,
		ProblemId:    problemId,
		SubmissionId: submissionId,
	}, nil
}
