package contestsubmissionusecase

import (
	domain "contest/src/domain/entity"
	"context"
)

type ContestSubmissionInteractor interface {
	InternalIngestContestSubmission(
		ctx context.Context,
		input *InternalIngestContestSubmissionInput,
	) (*InternalIngestContestSubmissionOutput, error)
}

type (
	InternalIngestContestSubmissionInput struct {
		SubmissionId string         `json:"submission_id,omitempty"`
		Verdict      domain.Verdict `json:"verdict,omitempty"`
		Points       float64        `json:"points,omitempty"`
	}

	InternalIngestContestSubmissionOutput struct {
		Message string `json:"message,omitempty"`
	}
)
