package impl

import (
	"bytes"
	"context"
	"errors"
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
		err = js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, 0, 0, 0, 0, JudgementFailedMessage)
		if err != nil {
			i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
		}
		return err
	}

	defer js.pService.Put(i)
	defer i.Logger.Debug().Msgf("Returning isolate to pool, number in pool will be: %d", js.pService.Len())

	i.Logger.Debug().Msgf("Took out an isolate, number of isolate remains in the pool is: %d", js.pService.Len())
	// Prepare all the nessessary files
	err = js.Prep(ctx, i, lang, req, problemInfo)
	if err != nil {
		i.Logger.Debug().Msgf("Error: %v", err)
		// i.Logger.Debug().Msgf("Judgement failed or CompilationError, return the isolate, number of isolate in the pool is: %d", js.pService.Len())
		return err
	}

	i.Logger.Debug().Msgf("Preparation success!, keep using the isolate, number of isolate in the pool is: %d", js.pService.Len())

	switch req.SubmissionType {
	case domain.SubmissionType(domain.ICPC):
		err = js.JudgeICPC(ctx, i, lang, req, problemInfo)
	default:
		// js.pService.Put(i)
		i.Logger.Error().Err(err).Msgf("Other submission type is not supported")
		err = judge.UnsupportedSubmissionType
	}

	return err
}

// This function will help copy/create the nessessary files into the isolate working directory
func (js *JudgeServiceImpl) Prep(ctx context.Context, i *domain.Isolate, lang pkg.Language, req *isolateservice.SubmissionRequest, problemInfo *problem.ProblemServiceGetOutput) error {
	i.Logger.Info().Msgf("Assigned to submission with id: %s", (*req).SubmissionId)

	var (
		errBuf bytes.Buffer
		err    error
		vert   *judge.RunVerdict
	)

	// always remember to return the isolate instance
	// defer func() {
	// 	if err != nil {
	// 		judgeutils.ReturnIsolateIfFail(js.pService, i, err)
	// 	}
	// }()
	//
	_, err = utils.CreateSubmissionSourceFile(i, req.Sourcecode, req.SubmissionId, lang.DefaultFileName())
	if err != nil {
		return err
	}

	i.Logger.Info().Msgf("Created source file inside the isolate working directory")

	lang.Compile(i, req, &errBuf)
	vert, err = judgeutils.CheckRunStatus(i, req.SubmissionId)
	msg := judgeutils.GetCompileMessage(vert, errBuf.String())

	i.Logger.Info().Msgf("Compile message: %s", errBuf.String())
	if err != nil {
		err = js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if err != nil {
			i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
		}
		return err
	}

	switch vert.Status {
	case "RE", "SG", "TO", "XX":
		err = js.updateFinal(ctx, req.EvalId, domain.COMPILATION_ERROR, vert.Time, vert.MaxRss, 0, 0, msg)
		if err != nil {
			i.Logger.Error().Msgf("Database error, can't update verdict: %v", err)
		}
		err = judge.CompilationError
		i.Logger.Debug().Msgf("Compile error: %v", err)
	default:
		// This is just to detect if the program failed to compile via the information given by the meta file,
		// this is basically hardcoding and i might have to find a way to make this cleaner
		if vert.Status != "" || vert.ExitCode != 0 {
			err = js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
			if err != nil {
				i.Logger.Error().Msgf("Database error, can't update verdict: %v", err)
			}
			err = judge.JugdgementFailed
		}

	}
	i.Logger.Debug().Msgf("Compile error: %v", err)
	if err != nil {
		return err
	}

	err = js.prepChecker(ctx, i, req, vert)

	// Prepare the interactor file
	if problemInfo.IsInteractive {
		err = js.prepInteractor(ctx, i, req, vert)
	}

	return nil
}

func (js *JudgeServiceImpl) prepChecker(ctx context.Context, i *domain.Isolate, req *isolateservice.SubmissionRequest, vert *judge.RunVerdict) error {
	checkerLocation, err := js.problemService.GetCheckerAddr(req.ProblemId)
	if err != nil {
		err = js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if err != nil {
			i.Logger.Error().Msgf("Database error, can't update verdict: %v", err)
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
	interactorLocation, err := js.problemService.GetInteractorAddr(req.ProblemId)
	if err != nil {
		err := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if err != nil {
			i.Logger.Error().Msgf("Database error, can't update verdict: %v", err)
		}
		return err
	}
	crossrunLocation, err := js.problemService.GetCrossRunAddr(req.ProblemId)
	if err != nil {
		err := js.updateFinal(ctx, req.EvalId, domain.JUDGEMENT_FAILED, vert.Time, vert.MaxRss, 0, 0, vert.Message)
		if err != nil {
			i.Logger.Error().Msgf("Database error, can't update verdict: %v", err)
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
		js.OnFail(ctx, i, req.EvalId, *curCpu, *curMem, tc-1, JudgementFailedMessage)
		return true, err
	}

	outaddr := utils.GetSubmissionDir(i, req.SubmissionId) + "/output_" + strconv.Itoa(tc)

	fout, err := os.Create(outaddr)
	if err != nil {
		i.Logger.Panic().Msgf("Error occured when trying to create new output file: %v", err)
		return true, err
	}
	defer fout.Close()

	fin, err := os.Open(tcInputAddr)
	if err != nil {
		i.Logger.Panic().Msgf("Error occured when trying to read input file: %v", err)
		return true, err
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

	i.Logger.Debug().Msgf("Start to run the code, config is: %v", rc)
	if problemInfo.IsInteractive {
		runCmd, err := lang.RunCmdStrNoStream(i, &rc, req)
		if err != nil {
			return true, nil
		}
		interactorAddr := judgeutils.GetSubmissionInteractorAddr(i, req)
		crossrunAddr := judgeutils.GetSubmissionCrossRunJarAddr(i, req)
		reportAddr := judgeutils.GetSubmissionReportFileAddr(i, req)
		iivert, _, msg, err := js.interactorService.RunInteractor(crossrunAddr, interactorAddr, tcInputAddr, outaddr, tcAnsAddtr, reportAddr, runCmd)
		if err != nil {
			return true, nil
		}
		ivert = iivert
		i.Logger.Debug().Msgf("Interactor response: vert: %v, msg: %s", ivert, msg)
	} else {
		err = lang.Run(i, &rc, req)
	}

	vert, err := judgeutils.CheckRunStatus(i, req.SubmissionId)
	i.Logger.Debug().Msgf("Run Status from MetaFile: %v", vert)
	if err != nil {
		js.OnFail(ctx, i, req.EvalId, *curCpu, *curMem, tc-1, JudgementFailedMessage)
		return true, err
	}

	*curCpu = max(*curCpu, vert.Time)
	*curMem = max(*curMem, vert.MaxRss)

	// checkerLocation, e := js.problemService.GetCheckerAddr(req.ProblemId)
	// err = e
	// if err != nil {
	// 	js.OnFail(ctx, i, req.EvalId, *curCpu, *curMem, tc-1, JudgementFailedMessage)
	// 	return true, err
	// }

	// If this is an interactive problem, the verdict must be determine here first
	if problemInfo.IsInteractive && ivert != domain.ACCEPTED {
		msg, err := judgeutils.ReadInteractiveReportFile(i, req)
		err = js.updateCase(ctx, req.EvalId, ivert, vert.Time, vert.MaxRss, msg, 1, *curCpu, *curMem, tc-1)
		err = js.updateFinal(ctx, req.EvalId, ivert, *curCpu, *curMem, tc-1, 0, msg)
		if err != nil {
			i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
		}
		// js.pService.Put(i)
		return true, nil
	}

	checkerLocation := judgeutils.GetSubmissionCheckerAddr(i, req)
	cvert, msg, e := js.checkVerdict(vert, checkerLocation, tcInputAddr, outaddr, tcAnsAddtr)
	err = e
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
		i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
		return true, err
	}

	if cvert != domain.ACCEPTED {
		err = js.updateFinal(ctx, req.EvalId, cvert, *curCpu, *curMem, tc-1, 0, msg)
		if err != nil {
			i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
		}
		// js.pService.Put(i)
		return true, nil
	}
	return false, nil
}
func (js *JudgeServiceImpl) JudgeICPC(
	ctx context.Context,
	i *domain.Isolate,
	lang pkg.Language,
	req *isolateservice.SubmissionRequest,
	problemInfo *problem.ProblemServiceGetOutput,
) error {
	var (
		err error
	)
	// defer func() {
	// 	if err != nil {
	// 		judgeutils.ReturnIsolateIfFail(js.pService, i, err)
	// 	}
	// }()

	var (
		curCpuTime     float64       = 0
		curMemoryUsage memory.Memory = 0
	)

	for tc := 1; tc <= problemInfo.TestNum; tc += 1 {
		done, e := js.RunCase(ctx, i, lang, req, problemInfo, tc, &curCpuTime, &curMemoryUsage)
		// err = e
		if done || e != nil {
			// js.pService.Put(i)
			return err
		}
	}

	err = js.updateFinal(ctx, req.EvalId, domain.ACCEPTED, curCpuTime, curMemoryUsage, problemInfo.TestNum, 1, AcceptedMessage)
	if err != nil {
		i.Logger.Panic().Msgf("Database error, can't update verdict: %v", err)
	}
	// js.pService.Put(i)
	i.Logger.Debug().Msgf("Judgement success!, return the isolate, number of isolate in the pool is: %d", js.pService.Len())

	return nil
}

func (js *JudgeServiceImpl) checkVerdict(vert *judge.RunVerdict, checkerAddr, inputAddr, outputAddr, answerAddr string) (domain.Verdict, string, error) {
	config.GetLogger().Debug().Msgf("Status is: %s", vert.Status)
	switch vert.Status {
	case "TO":
		return domain.TIME_LIMIT_EXCEEDED, vert.Message, nil
	case "XX":
		return domain.JUDGEMENT_FAILED, vert.Message, nil
	}

	if vert.CgOomKilled == 1 {
		config.GetLogger().Debug().Msgf("MLE detected via cg-oom-killed")
		return domain.MEMORY_LIMIT_EXCEEDED, vert.Message, nil
	}

	switch vert.Status {
	case "RE", "SG":
		return domain.RUNTIME_ERROR, vert.Message, nil
	}

	// this return the message and the exit code, which must be use later
	// TODO: Do something with exit code and checker message
	cvert, _, msg, err := js.checkerService.RunChecker(checkerAddr, inputAddr, outputAddr, answerAddr)
	if err != nil {
		return "", vert.Message, err
	}
	return cvert, msg, nil
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
