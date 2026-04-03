package impl

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	// "fmt"
	"os"
	"strconv"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	evalrepository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation"
	redisrepository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission"
	screpository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode"
	subrepository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/checker"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/interactor"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/utils"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/judge"
	judgeutils "github.com/bibimoni/Online-judge/submission-judge/src/service/judge/utils"
	poolservice "github.com/bibimoni/Online-judge/submission-judge/src/service/pool"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"

	// "github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	isubmission_utils "github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission/utils"
)

type JudgeServiceImpl struct {
	pService          poolservice.PoolService
	problemService    problem.ProblemService
	evalRepo          evalrepository.EvaluationRepository
	checkerService    checker.CheckerService
	interactorService interactor.InteractorService
	redisRepo         redisrepository.RedisSubmissionRepository
	submissionRepo    subrepository.SubmissionRepository
	sourcecodeRepo    screpository.SourcecodeRepository
	isolateservice    isolateservice.IsolateService
}

var CompilationError = errors.New("Compile error!")
var JudgementFailedMessage = "failed to evaluate submission"
var AcceptedMessage = "Accepted"

func NewJudgeServiceImpl(
	pService poolservice.PoolService,
	problemService problem.ProblemService,
	evalRepo evalrepository.EvaluationRepository,
	checkerS checker.CheckerService,
	interactorS interactor.InteractorService,
	redis redisrepository.RedisSubmissionRepository,
	submissionRepo subrepository.SubmissionRepository,
	sourcecodeRepo screpository.SourcecodeRepository,
	isolateservice isolateservice.IsolateService,
) *JudgeServiceImpl {
	js := &JudgeServiceImpl{
		pService:          pService,
		problemService:    problemService,
		evalRepo:          evalRepo,
		checkerService:    checkerS,
		interactorService: interactorS,
		redisRepo:         redis,
		submissionRepo:    submissionRepo,
		sourcecodeRepo:    sourcecodeRepo,
		isolateservice:    isolateservice,
	}

	return js
}

func NewJudgeService(
	pService poolservice.PoolService,
	problemService problem.ProblemService,
	evalRepo evalrepository.EvaluationRepository,
	checkerS checker.CheckerService,
	interactorS interactor.InteractorService,
	redis redisrepository.RedisSubmissionRepository,
	submissionRepo subrepository.SubmissionRepository,
	sourcecodeRepo screpository.SourcecodeRepository,
	isoisolateservice isolateservice.IsolateService,
) judge.JudgeService {
	return NewJudgeServiceImpl(pService, problemService, evalRepo, checkerS, interactorS, redis, submissionRepo, sourcecodeRepo, isoisolateservice)
}

// If this worked, i might have to remove the unnecessary parameter
func (js *JudgeServiceImpl) Judge(ctx context.Context, req *isolateservice.SubmissionRequest, problemInfo *problem.ProblemServiceGetOutput) error {
	return js.redisRepo.PushSubmissionJob(ctx, req)
}

func (js *JudgeServiceImpl) JudgeStart(ctx context.Context, lang pkg.Language, req *isolateservice.SubmissionRequest, problemInfo *problem.ProblemServiceGetOutput) error {
	// update PENDING to Websocket
	err := js.updateWS(ctx, req.EvalId)
	if err != nil {
		config.GetLogger().Panic().Msgf("redis stopped working: %v", err)
		return err
	}

	// This should be here, inside judgeStart
	i, err := js.pService.Get()
	if err != nil {
		// Very unlikely, happen when channel is closed
		updateErr := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, 0, 0, 0, 0, JudgementFailedMessage)
		if updateErr != nil {
			i.Logger.Panic().Err(err).Msgf("Database error")
		}
		return err
	}

	defer js.pService.Put(i)
	defer i.Logger.Debug().Msgf("Returning isolate to pool, number in pool will be: %d", js.pService.Len())

	i.Logger.Info().Msgf("Judging submission: %s, isolates remaining: %d",
		req.SubmissionId, js.pService.Len())
	// Prepare all the nessessary files
	err = js.Prep(ctx, i, lang, req, problemInfo)
	if err != nil {
		i.Logger.Debug().Msgf("Preparation error: %v", err)
		return err
	}

	i.Logger.Debug().Msgf("Preparation successful, starting test case execution")

	switch req.SubmissionType {
	case domain.SubmissionType(domain.ICPC):
		err = js.JudgeICPC(ctx, i, lang, req, problemInfo)
	default:
		// js.pService.Put(i)
		i.Logger.Error().Msgf("Unsupported submission type: %v", req.SubmissionType)

	}

	return err
}

// This function will help copy/create the nessessary files into the isolate working directory
func (js *JudgeServiceImpl) Prep(ctx context.Context, i *domain.Isolate, lang pkg.Language, req *isolateservice.SubmissionRequest, problemInfo *problem.ProblemServiceGetOutput) error {
	i.Logger.Info().Msgf("Assigned to submission with id: %s", (*req).SubmissionId)
	_, err := utils.CreateSubmissionSourceFile(i, req.Sourcecode, req.SubmissionId, lang.DefaultFileName())
	if err != nil {
		i.Logger.Error().Err(err).Msgf("Error creating source file")
		return err
	}
	i.Logger.Debug().Msgf("Source file created successfully!")

	var compileOutput bytes.Buffer
	lang.Compile(i, req, &compileOutput)

	compileVerdict, err := judgeutils.CheckRunStatus(i, req.SubmissionId)
	i.Logger.Info().Msgf("Compile message: %s", compileOutput.String())
	if err != nil {
		i.Logger.Error().Err(err).Msgf("Failed to read compile meta file")
		updateErr := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, 0, 0, 0, 0, JudgementFailedMessage)
		if updateErr != nil {
			i.Logger.Panic().Err(updateErr).Msgf("Database error")
		}
		return err
	}

	compileMsg := judgeutils.GetCompileMessage(compileVerdict, compileOutput.String())
	i.Logger.Debug().Msgf("Compile output: %s", compileOutput.String())

	if !js.isCompilationSuccessful(compileVerdict) {
		i.Logger.Info().Msgf("Compilation failed: status=%s, exitcode=%d", compileVerdict.Status, compileVerdict.ExitCode)
		updateErr := js.updateFinal(ctx, req.EvalId, domain.COMPILATION_ERROR, compileVerdict.Time, compileVerdict.MaxRss, 0, 0, compileMsg)
		if updateErr != nil {
			i.Logger.Panic().Err(updateErr).Msgf("Database error")
		}
		return judge.CompilationError
	}

	err = js.prepChecker(ctx, i, req, compileVerdict)
	if err != nil {
		i.Logger.Error().Err(err).Msgf("Error preparing checker")
		return err
	}

	// Prepare the interactor file
	if problemInfo.IsInteractive {
		err = js.prepInteractor(ctx, i, req, compileVerdict)
		if err != nil {
			i.Logger.Error().Err(err).Msgf("Error preparing interactor")
			return err
		}
	}

	i.Logger.Info().Msgf("Preparation completed successfully!")
	return nil
}

func (js *JudgeServiceImpl) isCompilationSuccessful(vert *judge.RunVerdict) bool {
	if vert.ExitCode == 0 && (vert.Status == "OK" || vert.Status == "") {
		return true
	}
	return false
}

func (js *JudgeServiceImpl) prepChecker(ctx context.Context, i *domain.Isolate, req *isolateservice.SubmissionRequest, vert *judge.RunVerdict) error {
	checkerLocation, err := js.problemService.GetCheckerAddr(req.ProblemId, req.ProblemVersion)
	if err != nil {
		updateErr := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if updateErr != nil {
			i.Logger.Error().Err(updateErr).Msg("Database error")
		}
		return err
	}
	err = utils.CopyChecker(i, (*req).SubmissionId, checkerLocation)
	if err != nil {
		return err
	}
	i.Logger.Info().Msgf("Checker file copied successfully!")
	return err
}

func (js *JudgeServiceImpl) prepInteractor(ctx context.Context, i *domain.Isolate, req *isolateservice.SubmissionRequest, vert *judge.RunVerdict) error {
	interactorLocation, err := js.problemService.GetInteractorAddr(req.ProblemId, req.ProblemVersion)
	if err != nil {
		updateErr := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if updateErr != nil {
			i.Logger.Error().Err(updateErr).Msg("Database error")
		}
		return err
	}
	crossrunLocation, err := js.problemService.GetCrossRunAddr(req.ProblemId, req.ProblemVersion)
	if err != nil {
		updateErr := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if updateErr != nil {
			i.Logger.Error().Err(updateErr).Msg("Database error")
		}
		return err
	}

	err = utils.CopyInteractor(i, (*req).SubmissionId, interactorLocation)
	if err != nil {
		return err
	}
	err = utils.CopyCrossRun(i, (*req).SubmissionId, crossrunLocation)
	if err != nil {
		return err
	}
	i.Logger.Info().Msgf("Interactor files (interactor, crossrun) copied successfully!")
	return err
}

func (js *JudgeServiceImpl) OnFail(
	ctx context.Context,
	i *domain.Isolate,
	evalId string,
	curCpu float64,
	curMem memory.Memory,
	tcSuccess int,
	msg string,
) {
	if err := js.updateFinal(ctx, evalId, domain.JUDGEMENT_FAILED, curCpu, curMem, tcSuccess, 0, msg); err != nil {
		i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
	}
}

func (js *JudgeServiceImpl) RunCase(
	ctx context.Context,
	i *domain.Isolate,
	lang pkg.Language,
	req *isolateservice.SubmissionRequest,
	problemInfo *problem.ProblemServiceGetOutput,
	tc int,
	curCpu *float64,
	curMem *memory.Memory,
) (done bool, err error) {
	var ivert domain.Verdict
	tcInputAddr, err := js.problemService.GetTestCaseAddr(req.ProblemId, problem.TestCaseType(problem.INPUT), tc)
	if err != nil {
		js.OnFail(ctx, i, req.EvalId, *curCpu, *curMem, tc-1, JudgementFailedMessage)
		return true, err
	}

	tcAnsAddtr, err := js.problemService.GetTestCaseAddr(req.ProblemId, problem.TestCaseType(problem.OUTPUT), tc)
	if err != nil {
		return nil, "", "", err
	}

	outaddr := utils.GetSubmissionDir(i, req.SubmissionId) + "/output_" + strconv.Itoa(tc)

	fout, err := os.Create(outaddr)
	if err != nil {
		i.Logger.Panic().Msgf("Error occured when trying to create new output file: %v", err)
		return nil, "", "", err
	}
	defer fout.Close()

	fin, err := os.Open(tcInputAddr)
	if err != nil {
		i.Logger.Panic().Msgf("Error occured when trying to read input file: %v", err)
		return nil, "", "", err
	}
	defer fin.Close()

	rc := domain.RunConfig{
		TimeLimit:    time.Millisecond * time.Duration(problemInfo.TimeLimit),
		MemoryLimit:  memory.Memory(problemInfo.MemoryLimit),
		Meta:         true,
		Stdout:       fout,
		Stdin:        fin,
		MaxProcesses: 1,
	}

	var interactorVerdict domain.Verdict
	if problemInfo.IsInteractive {
		interactorVerdict, err = js.runInteractive(i, lang, req, &rc, tcInputAddr, outaddr, tcAnsAddtr)
		if err != nil {
			return nil, "", "", err
		}
	} else {
		lang.Run(i, &rc, req)
	}

	vert, err := judgeutils.CheckRunStatus(i, req.SubmissionId)
	if err != nil {
		return nil, "", "", err
	}
	return vert, outaddr, interactorVerdict, nil
}

func (js *JudgeServiceImpl) runInteractive(
	i *domain.Isolate,
	lang pkg.Language,
	req *isolateservice.SubmissionRequest,
	rc *domain.RunConfig,
	tcInputAddr, outaddr, tcAnsAddtr string,
) (domain.Verdict, error) {
	runCmd, err := lang.RunCmdStrNoStream(i, rc, req)
	if err != nil {
		return "", err
	}
	interactorAddr := judgeutils.GetSubmissionInteractorAddr(i, req)
	crossrunAddr := judgeutils.GetSubmissionCrossRunJarAddr(i, req)
	reportAddr := judgeutils.GetSubmissionReportFileAddr(i, req)

	verdict, _, msg, err := js.interactorService.RunInteractor(crossrunAddr, interactorAddr, tcInputAddr, outaddr, tcAnsAddtr, reportAddr, runCmd)
	if err != nil {
		js.OnFail(ctx, i, req.EvalId, *curCpu, *curMem, tc-1, JudgementFailedMessage)
		return true, err
	}

	curSuccess := tc
	if cvert != domain.ACCEPTED {
		curSuccess -= 1
	}
	err = js.updateCase(ctx, req.EvalId, cvert, vert.Time, vert.MaxRss, msg, 1, *curCpu, *curMem, curSuccess)

	if err != nil {
		result.Verdict = domain.JUDGEMENT_FAILED
		result.Message = JudgementFailedMessage
		result.ShouldStop = true
		return result
	}

	result.Verdict = verdict
	result.Message = msg
	result.Score = score

	result.ShouldStop = (verdict != domain.ACCEPTED && verdict != domain.PARTIAL_RESULT && verdict != domain.POINTS)
	return result
}

func (js *JudgeServiceImpl) JudgeICPC(
	ctx context.Context,
	i *domain.Isolate,
	lang pkg.Language,
	req *isolateservice.SubmissionRequest,
	problemInfo *problem.ProblemServiceGetOutput,
) error {
	var (
		curCpuTime     float64       = 0
		curMemoryUsage memory.Memory = 0
	)

	for tc := 1; tc <= problemInfo.TestNum; tc += 1 {
		result, err := js.RunCase(ctx, i, lang, req, problemInfo, tc, &curCpuTime, &curMemoryUsage)
		if err != nil {
			return err
		}

		curSuccess := tc
		if result.Verdict != domain.ACCEPTED {
			curSuccess = tc - 1
		}

		updateErr := js.updateCase(ctx, req.EvalId, result.Verdict, result.Time, result.Memory, result.Message, 1, curCpuTime, curMemoryUsage, curSuccess)
		if updateErr != nil {
			i.Logger.Panic().Err(updateErr).Msgf("Database error")
			return updateErr
		}

		if result.ShouldStop {
			updateErr := js.updateFinal(ctx, req.EvalId, result.Verdict, curCpuTime, curMemoryUsage, curSuccess, 0, result.Message)
			if updateErr != nil {
				i.Logger.Panic().Err(updateErr).Msgf("Database error")
				return updateErr
			}
			i.Logger.Debug().Msgf("ICPC Judging stopped at test case %d with verdict %s", tc, result.Verdict)
			return nil
		}
	}

	err := js.updateFinal(ctx, req.EvalId, domain.ACCEPTED, curCpuTime, curMemoryUsage, problemInfo.TestNum, 1, AcceptedMessage)
	if err != nil {
		i.Logger.Panic().Err(err).Msgf("Database error")
		return err
	}
	i.Logger.Debug().Msgf("ICPC Judging completed, All %d test cases passed!", problemInfo.TestNum)
	return nil
}

func (js *JudgeServiceImpl) JudgeIOI(
	ctx context.Context,
	i *domain.Isolate,
	lang pkg.Language,
	req *isolateservice.SubmissionRequest,
	problemInfo *problem.ProblemServiceGetOutput,
) error {
	i.Logger.Error().Msgf("IOI mode: Judging %d test cases with scoring", problemInfo.TestNum)
	var (
		curCpuTime     float64       = 0
		curMemoryUsage memory.Memory = 0
	)

	testScores := make(map[int]float64)
	testMaxScores := make(map[int]float64)
	verdicts := make(map[int]domain.Verdict)

	for tc := 1; tc <= problemInfo.TestNum; tc += 1 {
		idx := tc - 1
		result, err := js.RunCase(ctx, i, lang, req, problemInfo, tc, &curCpuTime, &curMemoryUsage)
		if err != nil {
			i.Logger.Error().Err(err).Msgf("Error running test case %d", tc)
			testScores[idx] = 0.0
			verdicts[idx] = domain.JUDGEMENT_FAILED
			continue
		}

		maxScore := 100.0
		if idx < len(problemInfo.TestMaxScores) {
			maxScore = problemInfo.TestMaxScores[idx]
		}

		actualScore := (result.Score / 100.0) * maxScore
		testScores[idx] = actualScore
		testMaxScores[idx] = maxScore
		verdicts[idx] = result.Verdict

		i.Logger.Debug().Msgf("Test case %d: verdict=%s, score=%.2f/%.2f", tc, result.Verdict, actualScore, maxScore)
		updateErr := js.updateCaseFloat(ctx, req.EvalId, result.Verdict, result.Time, result.Memory, result.Message, actualScore, maxScore, curCpuTime, curMemoryUsage, tc)
		if updateErr != nil {
			i.Logger.Panic().Err(updateErr).Msgf("Database error")
			return updateErr
		}
	}

	totalScore, totalMaxScore := CalculateTotalScore(problemInfo, testScores, testMaxScores, verdicts, i)
	var finalVerdict domain.Verdict
	if totalScore >= totalMaxScore {
		finalVerdict = domain.ACCEPTED
	} else {
		finalVerdict = domain.PARTIAL_RESULT
	}

	nSuccess := 0
	for tc := 1; tc <= problemInfo.TestNum; tc += 1 {
		idx := tc - 1
		if verdicts[idx] == domain.ACCEPTED {
			nSuccess += 1
		}
	}

	message := fmt.Sprintf("Total Score: %.2f/%.2f", min(totalScore, totalMaxScore), totalMaxScore)
	err := js.updateFinalFloat(ctx, req.EvalId, finalVerdict, curCpuTime, curMemoryUsage, nSuccess, totalScore, totalMaxScore, message)
	if err != nil {
		i.Logger.Panic().Err(err).Msgf("Database error")
		return err
	}

	i.Logger.Info().Msgf("IOI: %s - %s", finalVerdict, message)
	return nil
}

func (js *JudgeServiceImpl) checkRuntimeErrors(vert *judge.RunVerdict) domain.Verdict {
	switch vert.Status {
	case "TO":
		return domain.TIME_LIMIT_EXCEEDED
	case "XX":
		return domain.JUDGEMENT_FAILED
	}

	if vert.CgOomKilled == 1 {
		config.GetLogger().Debug().Msgf("MLE detected via cg-oom-killed")
		return domain.MEMORY_LIMIT_EXCEEDED
	}

	switch vert.Status {
	case "RE", "SG":
		return domain.RUNTIME_ERROR
	}
	return ""
}

func (js *JudgeServiceImpl) checkVerdict(vert *judge.RunVerdict, checkerAddr, inputAddr, outputAddr, answerAddr string) (domain.Verdict, string, float64, error) {
	config.GetLogger().Debug().Msgf("Status is: %s", vert.Status)
	switch vert.Status {
	case "TO":
		return domain.TIME_LIMIT_EXCEEDED, vert.Message, 0.0, nil
	case "XX":
		return domain.JUDGEMENT_FAILED, vert.Message, 0.0, nil
	}

	if vert.CgOomKilled == 1 {
		config.GetLogger().Debug().Msgf("MLE detected via cg-oom-killed")
		return domain.MEMORY_LIMIT_EXCEEDED, vert.Message, 0.0, nil
	}

	switch vert.Status {
	case "RE", "SG":
		return domain.RUNTIME_ERROR, vert.Message, 0.0, nil
	}

	// this return the message and the exit code, which must be use later
	// TODO: Do something with exit code and checker message
	cvert, _, msg, score, err := js.checkerService.RunChecker(checkerAddr, inputAddr, outputAddr, answerAddr)
	if err != nil {
		return "", vert.Message, 0.0, err
	}
	return cvert, msg, score, nil
}

func (js *JudgeServiceImpl) updateFinal(
	ctx context.Context,
	evalId string,
	verdict domain.Verdict,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
	points int,
	message string,
) error {
	err := js.evalRepo.UpdateFinal(
		ctx,
		evalId,
		verdict,
		cpuTime,
		memoryUsage,
		nsucess,
		points,
		message,
	)
	if err != nil {
		return err
	}
	return js.updateWS(ctx, evalId)
}

func (js *JudgeServiceImpl) updateCase(
	ctx context.Context,
	evalId string,
	verdictCase domain.Verdict,
	cpuTimeCase float64,
	memoryUsageCase memory.Memory,
	outputCase string,
	pointsCase int,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
) error {
	err := js.evalRepo.UpdateCase(
		ctx,
		evalId,
		verdictCase,
		cpuTimeCase,
		memoryUsageCase,
		outputCase,
		pointsCase,
		cpuTime,
		memoryUsage,
		nsucess,
	)
	if err != nil {
		return err
	}

	return js.updateWS(ctx, evalId)
}

func (js *JudgeServiceImpl) updateWS(ctx context.Context, evalId string) error {
	eval, err := js.evalRepo.GetEval(ctx, evalId)
	if err != nil {
		return err
	}

	wsUpdate, err := isubmission_utils.GetSubmissionWithoutSourceCode(
		ctx,
		js.evalRepo,
		js.sourcecodeRepo,
		js.submissionRepo,
		eval.SubmissionId.Hex(),
	)

	if err != nil {
		return err
	}

	return js.redisRepo.PulishSubmission(ctx, *wsUpdate)
}

func (js *JudgeServiceImpl) updateFinalFloat(
	ctx context.Context,
	evalId string,
	verdict domain.Verdict,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
	score float64,
	maxScore float64,
	message string,
) error {
	err := js.evalRepo.UpdateFinalFloat(
		ctx,
		evalId,
		verdict,
		cpuTime,
		memoryUsage,
		nsucess,
		score,
		maxScore,
		message,
	)
	if err != nil {
		return err
	}
	return js.updateWS(ctx, evalId)
}

func (js *JudgeServiceImpl) updateCaseFloat(
	ctx context.Context,
	evalId string,
	verdictCase domain.Verdict,
	cpuTimeCase float64,
	memoryUsageCase memory.Memory,
	outputCase string,
	scoreCase float64,
	maxScoreCase float64,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
) error {
	err := js.evalRepo.UpdateCaseFloat(
		ctx,
		evalId,
		verdictCase,
		cpuTimeCase,
		memoryUsageCase,
		outputCase,
		scoreCase,
		maxScoreCase,
		cpuTime,
		memoryUsage,
		nsucess,
	)
	if err != nil {
		return err
	}

	return js.updateWS(ctx, evalId)
}
