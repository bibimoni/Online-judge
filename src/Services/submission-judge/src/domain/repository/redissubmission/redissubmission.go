package repository

import (
	"context"

	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
)

type RedisSubmissionRepository interface {
	PulishSubmission(ctx context.Context, res usecase.WSSubmissionResponse) error
	Subscribe(ctx context.Context, channelId string) (<-chan *usecase.WSSubmissionResponse, error)
	GetChannelString(problemId, username, submissionId string) string
	PushSubmissionJob(ctx context.Context, req *isolateservice.SubmissionRequest) error
	PopSubmissionJob(ctx context.Context) (*isolateservice.SubmissionRequest, error)
}
