package contestsubmissionservice

import (
	domain "contest/src/domain/entity"
	"context"
)

type ContestSubmissionService interface {
	SubmitContestSubmissionToJudge(ctx context.Context, contestId string, problemLabel string, code string, language string, username string, submissionType domain.ParticipantType) (string, error)
	UpsertContestSubmissionFromJudgeEvent(ctx context.Context, submissionId string, verdict domain.Verdict, points float64) error
}
type SubmitSubmissionResponse struct {
	Success bool                         `json:"success"`
	Data    SubmitSubmissionResponseData `json:"data"`
}

type SubmitSubmissionResponseData struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
