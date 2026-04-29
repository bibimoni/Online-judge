package contestsubmissioninteractor

import (
	contestsubmissionservice "contest/src/service/contest-submission"
	contestsubmissionusecase "contest/src/usecase/contest-submission"
	"context"
)

type ContestSubmissionInteractor struct {
	contestsubmissionService contestsubmissionservice.ContestSubmissionService
}

func NewContestSubmissionInteractor(
	contestsubmissionService contestsubmissionservice.ContestSubmissionService,
) *ContestSubmissionInteractor {
	return &ContestSubmissionInteractor{
		contestsubmissionService: contestsubmissionService,
	}
}

func (csi *ContestSubmissionInteractor) InternalIngestContestSubmission(
	ctx context.Context,
	input *contestsubmissionusecase.InternalIngestContestSubmissionInput,
) (*contestsubmissionusecase.InternalIngestContestSubmissionOutput, error) {
	err := csi.contestsubmissionService.UpsertContestSubmissionFromJudgeEvent(
		ctx,
		input.SubmissionId,
		input.Verdict,
		input.Points,
	)
	if err != nil {
		return nil, err
	}

	return &contestsubmissionusecase.InternalIngestContestSubmissionOutput{
		Message: "ingested successfully",
	}, nil
}
