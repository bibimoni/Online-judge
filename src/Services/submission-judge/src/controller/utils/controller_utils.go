package controller_utils

import (
	appctx "github.com/bibimoni/Online-judge/submission-judge/src/components"
	ei "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation/impl"
	ri "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission/impl"
	sci "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode/impl"
	si "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission/impl"
	checkerimpl "github.com/bibimoni/Online-judge/submission-judge/src/service/checker/impl"
	interactorimpl "github.com/bibimoni/Online-judge/submission-judge/src/service/interactor/impl"
	ii "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/impl"
	ji "github.com/bibimoni/Online-judge/submission-judge/src/service/judge/impl"
	pi "github.com/bibimoni/Online-judge/submission-judge/src/service/problem/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission/interactor"

	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
)

func InitInteractor(appContext appctx.AppContext) (*interactor.SubmissionInteractor, error) {
	log := config.GetLogger()

	db := appContext.GetMainDbConnection()
	submissionRepo := si.NewSubmissionRepository(db)
	sourcecodeRepo := sci.NewSourcecodeRepository(db)
	problemSvc, err := pi.NewProblemService()
	evalRepo := ei.NewEvaluationRepository(db)
	checkerS := checkerimpl.NewCheckerService()
	interactorS := interactorimpl.NewInteractorService()
	redis := ri.NewRedisSubmissionRepository(appContext.GetRedis())
	isolateS, _ := ii.NewIsolateService()
	judgeSvc := ji.NewJudgeServiceImpl(appContext.GetPool(), problemSvc, evalRepo, checkerS, interactorS, redis, submissionRepo, sourcecodeRepo, isolateS)
	if err != nil {
		log.Error().Msgf("Can't initialize submit request, got error : %v", err)
		return nil, err
	}
	return interactor.NewSubmissionInteractor(submissionRepo, sourcecodeRepo, problemSvc, judgeSvc, evalRepo), nil
}
