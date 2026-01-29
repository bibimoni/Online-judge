package interactor

import (
	"context"
	"slices"
	"strconv"

	// "strconv"

	evalRepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation"
	scr "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode"
	repository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	sr "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/judge"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission"
	isubmission_utils "github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission/utils"
	usecasews "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
)

type SubmissionInteractor struct {
	submissionRepo sr.SubmissionRepository
	sourcecodeRepo scr.SourcecodeRepository
	problemService problem.ProblemService
	judgeService   judge.JudgeService
	evalRepo       evalRepo.EvaluationRepository
}

func NewSubmissionInteractor(
	sr sr.SubmissionRepository,
	scr scr.SourcecodeRepository,
	ps problem.ProblemService,
	js judge.JudgeService,
	er evalRepo.EvaluationRepository,
) *SubmissionInteractor {
	return &SubmissionInteractor{
		submissionRepo: sr,
		sourcecodeRepo: scr,
		problemService: ps,
		judgeService:   js,
		evalRepo:       er,
	}
}

func (si *SubmissionInteractor) SubmitSubmission(ctx context.Context, input *usecase.SubmitSubmissionInput) (output *usecase.SubmitSubmissionResponse, err error) {
	log := config.GetLogger()
	log.Info().Msgf("User %s submitted a solution in %s, for problem with problem id: %s", input.Username, input.LanguageId, input.ProblemId)

	problemInfo, err := si.problemService.Get(ctx, input.ProblemId)
	if err != nil {
		return nil, err
	}

	params := repository.CreateSubmissionInput{
		ProblemId: strconv.FormatInt(problemInfo.ProblemId, 10),
		Username:  input.Username,
		Type:      input.SubmissionType,
	}
	submissionId, err := si.submissionRepo.CreateSubmission(ctx, params)
	if err != nil {
		log.Debug().Msgf("error happened when trying to create new submission: %v", err)
		return nil, err
	}

	_, err = si.sourcecodeRepo.CreateSourcecode(ctx, input.Code, input.LanguageId, submissionId)
	if err != nil {
		return nil, err
	}

	evalId, err := si.evalRepo.CreateEval(ctx, submissionId, problemInfo.TimeLimit, memory.Memory(problemInfo.MemoryLimit), problemInfo.TestNum)
	if err != nil {
		return nil, err
	}

	req := isolateservice.SubmissionRequest{
		SubmissionId:   submissionId,
		Username:       input.Username,
		Sourcecode:     input.Code,
		SubmissionType: input.SubmissionType,
		ProblemId:      input.ProblemId,
		LanguageId:     input.LanguageId,
		EvalId:         evalId,
	}

	log.Info().Msgf("Enqueue submission, id: %s. With eval id: %s", submissionId, evalId)
	err = si.judgeService.Judge(ctx, &req, problemInfo)

	return &usecase.SubmitSubmissionResponse{
		Message: "Submit successfully!",
		ID:      submissionId,
	}, nil
}

func (si *SubmissionInteractor) RejudgeSubmission(ctx context.Context, input *usecase.RejudgeSubmissionInput) (*usecase.RejudgeSubmissionOutput, error) {
	validSubmissionRequests, err := isubmission_utils.GetSubmissionRequests(
		ctx,
		si.evalRepo,
		si.sourcecodeRepo,
		si.submissionRepo,
		input.SubmissionIds,
	)

	if err != nil {
		return nil, err
	}

	resp := &usecase.RejudgeSubmissionOutput{
		RejudgeSuccessSubmissionIds: make([]string, 0, len(validSubmissionRequests)),
	}

	var distinctProblemIds = make([]string, 0)
	for _, req := range validSubmissionRequests {
		distinctProblemIds = append(distinctProblemIds, req.ProblemId)
	}

	slices.Sort(distinctProblemIds)
	distinctProblemIds = slices.Compact(distinctProblemIds)

	var problemInfos = make(map[string]*problem.ProblemServiceGetOutput)
	for _, pid := range distinctProblemIds {
		pinfo, err := si.problemService.Get(ctx, pid)
		if err != nil {
			return nil, err
		}
		problemInfos[pid] = pinfo
	}

	for _, req := range validSubmissionRequests {
		resp.RejudgeSuccessSubmissionIds = append(resp.RejudgeSuccessSubmissionIds, req.SubmissionId)
		si.judgeService.Judge(ctx, &req, problemInfos[req.ProblemId])
	}

	config.GetLogger().Info().Msgf("Rejudge %d submissions", len(validSubmissionRequests))

	return resp, nil
}

func (si *SubmissionInteractor) GetSubmission(ctx context.Context, input *usecase.GetSubmissionInput) (*usecase.GetSubmissionOutput, error) {
	return isubmission_utils.GetSubmission(
		ctx,
		si.evalRepo,
		si.sourcecodeRepo,
		si.submissionRepo,
		input.SubmissionId,
	)
}

func (si *SubmissionInteractor) GetProblemSubmission(ctx context.Context, input *usecase.GetProblemSubmissionInput) (*usecase.GetProblemSubmissionOutput, error) {
	strs, err := si.submissionRepo.FindAllProblemSubmissionIds(ctx, input.ProblemId)
	if err != nil {
		return nil, err
	}

	config.GetLogger().Debug().Msgf("strings: %v", strs)
	var problemSubmissions []usecasews.WSSubmissionResponse
	for _, str := range strs {
		item, err := isubmission_utils.GetSubmissionWithoutSourceCode(ctx, si.evalRepo, si.sourcecodeRepo, si.submissionRepo, str)
		if err != nil {
			return nil, err
		}
		problemSubmissions = append(problemSubmissions, *item)
	}

	return &usecase.GetProblemSubmissionOutput{
		Submissions: problemSubmissions,
	}, nil
}

func (si *SubmissionInteractor) InternalContestSubmitSubmission(ctx context.Context, input *usecase.InternalContestSubmitSubmissionInput) (*usecase.SubmitSubmissionResponse, error) {
	log := config.GetLogger()
	log.Info().Msgf("User %s submitted a solution in %s, for problem with problem id: %s", input.Username, input.LanguageId, input.ProblemId)

	problemInfo, err := si.problemService.Get(ctx, input.ProblemId)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get problem info")
		return nil, err
	}

	params := repository.CreateSubmissionInput{
		ProblemId: strconv.FormatInt(problemInfo.ProblemId, 10),
		Username:  input.Username,
		Type:      input.SubmissionType,
	}
	submissionId, err := si.submissionRepo.CreateSubmissionWithTimestamp(ctx, params, input.SubmitAt)
	if err != nil {
		log.Debug().Msgf("error happened when trying to create new submission: %v", err)
		return nil, err
	}

	_, err = si.sourcecodeRepo.CreateSourcecode(ctx, input.Code, input.LanguageId, submissionId)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create source code")
		return nil, err
	}

	evalId, err := si.evalRepo.CreateEval(ctx, submissionId, problemInfo.TimeLimit, memory.Memory(problemInfo.MemoryLimit), problemInfo.TestNum)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create evaluation")
		return nil, err
	}

	req := isolateservice.SubmissionRequest{
		SubmissionId:   submissionId,
		Username:       input.Username,
		Sourcecode:     input.Code,
		SubmissionType: input.SubmissionType,
		ProblemId:      input.ProblemId,
		LanguageId:     input.LanguageId,
		EvalId:         evalId,
		ContestId:      input.ContestId,
	}

	log.Info().Msgf("Enqueue submission, id: %s. With eval id: %s", submissionId, evalId)
	err = si.judgeService.Judge(ctx, &req, problemInfo)

	if err != nil {
		log.Error().Msgf("InternalContestSubmitSubmission error: %v", err)
		return nil, err
	}

	return &usecase.SubmitSubmissionResponse{
		Message: "Submit successfully!",
		ID:      submissionId,
	}, nil
}
