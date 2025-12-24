package impl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem/utils"
)

const PROBLEM_INFO_FILENAME = "problem.json"

type ProblemServiceImpl struct {
	problemServerAddr string
}

func NewProblemServiceImpl() (*ProblemServiceImpl, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config for PROBLEM SERVICE: %v", err)
	}
	return &ProblemServiceImpl{
		problemServerAddr: cfg.ProblemServerAddr,
	}, nil
}

func NewProblemService() (problem.ProblemService, error) {
	return NewProblemServiceImpl()
}

func (ps *ProblemServiceImpl) Get(ctx context.Context, id string) (*problem.ProblemServiceGetOutput, error) {
	req := common.APIRequest{
		Method:  "GET",
		URL:     ps.problemServerAddr + "get/" + id + "/" + PROBLEM_INFO_FILENAME,
		Timeout: 60 * time.Second,
	}

	result, err := common.SendRequest[problem.ProblemServiceGetOutput](ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("There is an error occured fetch requesting from PROBLEM SERVER")
	}
	if result.ScoringMode == "" || len(result.TestGroups) == 0 {
		result.ScoringMode = "ICPC"
	}

	return result, nil
}

func (ps *ProblemServiceImpl) GetTestCaseDirAddr(problemId string, tcType problem.TestCaseType) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId + "/tests"
	switch tcType {
	case problem.INPUT:
		stringAddr += "/input/"
	case problem.OUTPUT:
		stringAddr += "/output/"
	default:
		return "", fmt.Errorf("Please provide either INPUT or OUTPUT for testcase type")
	}

	err = utils.EnsureProblemDirectory(stringAddr)

	return stringAddr, nil
}

func (ps *ProblemServiceImpl) GetTestCaseAddr(problemId string, tcType problem.TestCaseType, testNum int) (string, error) {
	stringAddr, err := ps.GetTestCaseDirAddr(problemId, tcType)
	if err != nil {
		return "", err
	}
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	strTcType := "input"
	if tcType == problem.OUTPUT {
		strTcType = "output"
	}

	stringAddr += strconv.Itoa(testNum)
	remoteURL := fmt.Sprintf("%sget/%s/tests/%s/%s", cfg.ProblemServerAddr, problemId, strTcType, strconv.Itoa(testNum))

	log := config.GetLogger()
	log.Debug().Msgf("string address: %s, url: %s", stringAddr, remoteURL)

	// return stringAddr, nil
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}

func (ps *ProblemServiceImpl) GetCheckerAddr(problemId string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId
	err = utils.EnsureProblemDirectory(stringAddr)
	if err != nil {
		return "", err
	}
	stringAddr = stringAddr + "/checker"
	remoteURL := fmt.Sprintf("%sget/%s/checker", cfg.ProblemServerAddr, problemId)

	// return stringAddr, nil
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}

func (ps *ProblemServiceImpl) GetInteractorAddr(problemId string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId
	err = utils.EnsureProblemDirectory(stringAddr)
	if err != nil {
		return "", err
	}
	stringAddr = stringAddr + "/interactor"
	remoteURL := fmt.Sprintf("%sget/%s/interactor", cfg.ProblemServerAddr, problemId)
	// return stringAddr, nil
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}

func (ps *ProblemServiceImpl) GetCrossRunAddr(problemId string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId
	err = utils.EnsureProblemDirectory(stringAddr)
	if err != nil {
		return "", err
	}
	stringAddr = stringAddr + "/CrossRun.jar"
	remoteURL := fmt.Sprintf("%sget/%s/CrossRun.jar", cfg.ProblemServerAddr, problemId)
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}
