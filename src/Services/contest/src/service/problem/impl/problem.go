package problemserviceimpl

import (
	"contest/src/common"
	"contest/src/infrastructure/config"
	problemservice "contest/src/service/problem"
	"context"
	"fmt"
	"time"
)

const ProblemInfoFilename = "problem.json"

type ProblemServiceImpl struct {
	problemServiceAddr string
}

func NewProblemServiceImpl() (*ProblemServiceImpl, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config for PROBLEM SERVICE: %v", err)
	}
	return &ProblemServiceImpl{
		problemServiceAddr: cfg.ProblemServerAddr,
	}, nil
}

func NewProblemService() (problemservice.ProblemService, error) {
	return NewProblemServiceImpl()
}

func (ps *ProblemServiceImpl) Get(ctx context.Context, problemId string) (*problemservice.ProblemServiceGetOutput, error) {
	req := common.APIRequest{
		Method:  "GET",
		URL:     ps.problemServiceAddr + "get/" + problemId + "/" + ProblemInfoFilename,
		Timeout: 60 * time.Second,
	}

	result, err := common.SendRequest[problemservice.ProblemServiceGetOutput](ctx, req)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("there is an error occured fetch requesting from PROBLEM SERVICE")
	}

	return result, nil
}

func (ps *ProblemServiceImpl) GetProblemMaxScores(problem *problemservice.ProblemServiceGetOutput) float64 {
	totalScore := 0.0
	if len(problem.TestGroups) == 0 {
		for _, score := range problem.TestMaxScores {
			totalScore += score
		}
		return totalScore
	}

	for _, group := range problem.TestGroups {
		totalScore += group.MaxScore
	}
	return totalScore
}
