package contestsubmissioncontroller

import (
	"contest/src/common"
	contestsubmissionusecase "contest/src/usecase/contest-submission"

	"github.com/gin-gonic/gin"
)

func InternalIngestContestSubmission(contestsubmissionInteractor contestsubmissionusecase.ContestSubmissionInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toInternalIngestContestSubmissionInput,
		contestsubmissionInteractor.InternalIngestContestSubmission,
		common.WriteSuccessOutput[contestsubmissionusecase.InternalIngestContestSubmissionOutput],
	)
}

func toInternalIngestContestSubmissionInput(c *gin.Context) (*contestsubmissionusecase.InternalIngestContestSubmissionInput, error) {
	var req contestsubmissionusecase.InternalIngestContestSubmissionInput
	// here we bind because the param must be percise
	if err := c.BindJSON(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
