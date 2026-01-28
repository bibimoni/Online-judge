package contestsubmissionservice

import (
	domain "contest/src/domain/entity"
	"context"
)

type ContestSubmissionService interface {
	SubmitContestSubmissionToJudge(ctx context.Context, contestId string, problemLabel string, code string, languageId string, username string, submissionType domain.ParticipantType) (string, error)
	UpsertContestSubmissionFromJudgeEvent(ctx context.Context, submissionId string, verdict string, points float64) (string, error)
}
type SubmitSubmissionResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
