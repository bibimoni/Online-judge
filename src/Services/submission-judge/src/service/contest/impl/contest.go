package contestserviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	contestservice "github.com/bibimoni/Online-judge/submission-judge/src/service/contest"
)

const IngestContestSubmissionEndpoint = "internal/ingest-submission"

type ContestServiceImpl struct {
	contestServiceAddr string
	internalSecret     string
}

func NewContestServiceImpl() (*ContestServiceImpl, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config for CONTEST SERVICE : %v", err)
	}
	return &ContestServiceImpl{
		contestServiceAddr: cfg.ContestServerAddr,
		internalSecret:     cfg.InternalSecret,
	}, nil
}

func NewContestService() (*ContestServiceImpl, error) {
	return NewContestServiceImpl()
}

func (csi *ContestServiceImpl) IngestContestSubmission(
	ctx context.Context,
	submissionId string,
	verdict domain.Verdict,
	points float64,
) error {
	req := common.APIRequest{
		Method: "POST",
		URL:    csi.contestServiceAddr + IngestContestSubmissionEndpoint,
		Body: map[string]any{
			"submission_id": submissionId,
			"verdict":       string(verdict),
			"points":        points,
		},
		Headers: map[string]string{
			"X-Internal-Secret": csi.internalSecret,
		},
		Timeout: 10 * time.Second,
	}

	result, err := common.SendRequest[contestservice.IngestContestSubmissionResponse](ctx, req)
	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("there is an error occured fetch requesting from CONTEST SERVICE")
	}

	return nil
}
