package contestsubmissionserviceimpl

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	contestsubmissionrepo "contest/src/domain/repository/contest-submission"
	"contest/src/infrastructure/config"
	contestsubmissionservice "contest/src/service/contest-submission"
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	InternalJudgeSubmitEndpoint = "internal/contest/submit"
)

type ContestSubmissionServiceImpl struct {
	contestSubmissionServiceAddr string
	contestrepo                  contestrepo.ContestRepository
	contestsubmissionrepo        contestsubmissionrepo.ContestSubmissionRepository
}

func NewContestSubmissionServiceImpl(contestrepo contestrepo.ContestRepository, contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository) (*ContestSubmissionServiceImpl, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &ContestSubmissionServiceImpl{
		contestSubmissionServiceAddr: cfg.JudgeServerAddr,
		contestrepo:                  contestrepo,
		contestsubmissionrepo:        contestsubmissionrepo,
	}, nil
}

func NewContestSubmissionService(contestrepo contestrepo.ContestRepository, contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository) (*ContestSubmissionServiceImpl, error) {
	return NewContestSubmissionServiceImpl(contestrepo, contestsubmissionrepo)
}

func (css *ContestSubmissionServiceImpl) SubmitContestSubmissionToJudge(
	ctx context.Context,
	contestId string,
	problemLabel string,
	code string,
	language string,
	username string,
	submissionType domain.ParticipantType,
) (string, error) {
	contestProblems, err := css.contestrepo.GetProblemByLabels(ctx, contestId, []string{problemLabel})
	if err != nil {
		return "", err
	}

	// since we only request one problem label, we can safely get the first element
	contestProblem := contestProblems[0]
	submitAt := time.Now()
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	// submit to judge service
	req := common.APIRequest{
		Method: "POST",
		URL:    css.contestSubmissionServiceAddr + InternalJudgeSubmitEndpoint,
		Headers: map[string]string{
			"X-Internal-Secret": cfg.InternalSecret,
		},
		Body: map[string]string{
			"problem_id":      strconv.FormatUint(contestProblem.ProblemId, 10),
			"code":            code,
			"username":        username,
			"contest_id":      contestId,
			"language":        language,
			"submit_at":       submitAt.Format(time.RFC3339),
			"submission_type": string(submissionType),
		},
		Timeout: 30 * time.Second,
	}

	result, err := common.SendRequest[contestsubmissionservice.SubmitSubmissionResponse](ctx, req)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", fmt.Errorf("there is an error occured fetch requesting from JUDGE SERVICE")
	}

	contestSubmissionId, err := css.contestsubmissionrepo.Create(
		ctx,
		contestId,
		username,
		contestProblem.ProblemId,
		submissionType,
		submitAt,
		result.Data.ID,
	)

	if err != nil {
		return "", err
	}

	return contestSubmissionId, nil
}

func (css *ContestSubmissionServiceImpl) UpsertContestSubmissionFromJudgeEvent(ctx context.Context, submissionId string, verdict domain.Verdict, points float64) error {
	return css.contestsubmissionrepo.UpdateFromJudgeEvent(ctx, submissionId, verdict, points)
}