package contestsubmissionrepo

import (
	domain "contest/src/domain/entity"
	"context"
	"time"
)

type ContestSubmissionRepository interface {
	Create(
		ctx context.Context,
		contestId string,
		username string,
		contestProblemId string,
		submissionType domain.ParticipantType,
		submitAt time.Time,
		submissionId string,
	) (string, error)
	UpsertFromJudgeEvent(
		ctx context.Context,
		submissionId string,
		verdict domain.Verdict,
		points float64,
	) (string, error)
	ListByContest(
		ctx context.Context,
		contestId string,
		includeVirtual bool,
		includeUnrated bool,
	) ([]domain.ContestSubmission, error)
	SetIgnored(ctx context.Context, contestSubmissionId string, ignored bool) error
	SetIgnoreUser(ctx context.Context, contestId string, username string, ignored bool) error
}
