package contestservice

import (
	"context"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
)

type ContestService interface {
	IngestContestSubmission(ctx context.Context, submission_id string, verdict domain.Verdict, points float64) error
	IsProblemInActiveContest(ctx context.Context, problemId string) (bool, error)
}

type ProblemInActiveContestResponse struct {
	Success bool                               `json:"success"`
	Data    ProblemInActiveContestResponseData `json:"data"`
}

type ProblemInActiveContestResponseData struct {
	Locked bool `json:"locked"`
}

type IngestContestSubmissionResponse struct {
	Success bool                                `json:"success"`
	Data    IngestContestSubmissionResponseData `json:"data"`
}

type IngestContestSubmissionResponseData struct {
	Message string `json:"message"`
}
