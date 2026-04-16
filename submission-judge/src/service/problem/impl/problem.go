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

const ProblemInfoFilename = "problem.json"

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

func (ps *ProblemServiceImpl) GetLatestVersion(ctx context.Context, problemId string) (string, error) {
	req := common.APIRequest{
		Method:  "GET",
		URL:     ps.problemServerAddr + "latest-version?problemId=" + problemId,
		Timeout: 30 * time.Second,
	}

	result, err := common.SendRequest[problem.LatestVersionResponse](ctx, req)
	if err != nil {
		return "", err
	}
	if result == nil || !result.Success {
		return "", fmt.Errorf("failed to get latest version for problem %s", problemId)
	}

	return result.Data, nil
}

func (ps *ProblemServiceImpl) Get(ctx context.Context, id string) (*problem.ProblemServiceGetOutput, error) {
	version, err := ps.GetLatestVersion(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version for problem %s: %v", id, err)
	}

	req := common.APIRequest{
		Method:  "GET",
		URL:     ps.problemServerAddr + "get/" + id + "/" + version + "/" + ProblemInfoFilename,
		Timeout: 60 * time.Second,
	}

	result, err := common.SendRequest[problem.ProblemServiceGetOutput](ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("there is an error occured fetch requesting from PROBLEM SERVER")
	}
	if result.ScoringMode == "" || len(result.TestGroups) == 0 {
		result.ScoringMode = "ICPC"
	}

	return result, nil
}

func (ps *ProblemServiceImpl) GetTestCaseDirAddr(problemId, version string, tcType problem.TestCaseType) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId + "/" + version + "/tests"
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

func (ps *ProblemServiceImpl) GetTestCaseAddr(problemId, version string, tcType problem.TestCaseType, testNum int) (string, error) {
	stringAddr, err := ps.GetTestCaseDirAddr(problemId, version, tcType)
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
	remoteURL := fmt.Sprintf("%sget/%s/%s/tests/%s/%s", cfg.ProblemServerAddr, problemId, version, strTcType, strconv.Itoa(testNum))

	log := config.GetLogger()
	log.Debug().Msgf("string address: %s, url: %s", stringAddr, remoteURL)

	// return stringAddr, nil
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}

func (ps *ProblemServiceImpl) GetCheckerAddr(problemId, version string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId + "/" + version
	err = utils.EnsureProblemDirectory(stringAddr)
	if err != nil {
		return "", err
	}
	stringAddr = stringAddr + "/checker"
	remoteURL := fmt.Sprintf("%sget/%s/%s/checker", cfg.ProblemServerAddr, problemId, version)

	// return stringAddr, nil
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}

func (ps *ProblemServiceImpl) GetInteractorAddr(problemId, version string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId + "/" + version
	err = utils.EnsureProblemDirectory(stringAddr)
	if err != nil {
		return "", err
	}
	stringAddr = stringAddr + "/interactor"
	remoteURL := fmt.Sprintf("%sget/%s/%s/interactor", cfg.ProblemServerAddr, problemId, version)
	// return stringAddr, nil
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}

func (ps *ProblemServiceImpl) GetCrossRunAddr(problemId, version string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	stringAddr := cfg.ProblemsDir + "/" + problemId + "/" + version
	err = utils.EnsureProblemDirectory(stringAddr)
	if err != nil {
		return "", err
	}
	stringAddr = stringAddr + "/CrossRun.jar"
	remoteURL := fmt.Sprintf("%sget/%s/%s/CrossRun.jar", cfg.ProblemServerAddr, problemId, version)
	return utils.GetFileWithCache(context.Background(), stringAddr, remoteURL)
}
