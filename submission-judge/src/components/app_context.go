package appctx

import (
	evalrepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation"
	ei "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation/impl"
	redisrepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission"
	ri "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission/impl"
	sourcecoderepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode"
	sci "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode/impl"
	submissionrepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	si "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/checker"
	checkerimpl "github.com/bibimoni/Online-judge/submission-judge/src/service/checker/impl"
	contestservice "github.com/bibimoni/Online-judge/submission-judge/src/service/contest"
	contestserviceimpl "github.com/bibimoni/Online-judge/submission-judge/src/service/contest/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/interactor"
	interactorimpl "github.com/bibimoni/Online-judge/submission-judge/src/service/interactor/impl"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	ii "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/judge"
	ji "github.com/bibimoni/Online-judge/submission-judge/src/service/judge/impl"
	poolservice "github.com/bibimoni/Online-judge/submission-judge/src/service/pool"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	pi "github.com/bibimoni/Online-judge/submission-judge/src/service/problem/impl"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppContext interface {
	GetMainDbConnection() *mongo.Database
	GetPool() poolservice.PoolService
	GetRedis() *redis.Client
	GetSubmissionRepo() submissionrepo.SubmissionRepository
	GetSourcecodeRepo() sourcecoderepo.SourcecodeRepository
	GetProblemService() problem.ProblemService
	GetEvalRepo() evalrepo.EvaluationRepository
	GetCheckerSerivce() checker.CheckerService
	GetInteractorService() interactor.InteractorService
	GetRedisRepo() redisrepo.RedisSubmissionRepository
	GetIsolateService() isolateservice.IsolateService
	GetJudgeService() judge.JudgeService
	GetDB() *mongo.Database
	GetContestService() contestservice.ContestService
}

type appCtx struct {
	database          *mongo.Database
	pool              poolservice.PoolService
	rdb               *redis.Client
	submissionRepo    submissionrepo.SubmissionRepository
	sourcecodeRepo    sourcecoderepo.SourcecodeRepository
	problemService    problem.ProblemService
	evalRepo          evalrepo.EvaluationRepository
	checkerService    checker.CheckerService
	interactorService interactor.InteractorService
	redisRepo         redisrepo.RedisSubmissionRepository
	isolateService    isolateservice.IsolateService
	judgeService      judge.JudgeService
	contestSvc        contestservice.ContestService
}

func (ctx *appCtx) GetMainDbConnection() *mongo.Database                   { return ctx.database }
func (ctx *appCtx) GetPool() poolservice.PoolService                       { return ctx.pool }
func (ctx *appCtx) GetRedis() *redis.Client                                { return ctx.rdb }
func (ctx *appCtx) GetSubmissionRepo() submissionrepo.SubmissionRepository { return ctx.submissionRepo }
func (ctx *appCtx) GetSourcecodeRepo() sourcecoderepo.SourcecodeRepository { return ctx.sourcecodeRepo }
func (ctx *appCtx) GetProblemService() problem.ProblemService              { return ctx.problemService }
func (ctx *appCtx) GetEvalRepo() evalrepo.EvaluationRepository             { return ctx.evalRepo }
func (ctx *appCtx) GetCheckerSerivce() checker.CheckerService              { return ctx.checkerService }
func (ctx *appCtx) GetInteractorService() interactor.InteractorService     { return ctx.interactorService }
func (ctx *appCtx) GetRedisRepo() redisrepo.RedisSubmissionRepository      { return ctx.redisRepo }
func (ctx *appCtx) GetIsolateService() isolateservice.IsolateService       { return ctx.isolateService }
func (ctx *appCtx) GetJudgeService() judge.JudgeService                    { return ctx.judgeService }
func (ctx *appCtx) GetDB() *mongo.Database                                 { return ctx.database }
func (ctx *appCtx) GetContestService() contestservice.ContestService       { return ctx.contestSvc }

func NewAppContext(
	database *mongo.Database,
	pool poolservice.PoolService,
	rdb *redis.Client,
) *appCtx {
	problemSvc, err := pi.NewProblemService()
	if err != nil {
		config.GetLogger().Panic().Err(err).Msg("Can't not create problem service")
	}
	sourcecodeRepo := sci.NewSourcecodeRepository(database)
	evalRepo := ei.NewEvaluationRepository(database)
	submissionRepo := si.NewSubmissionRepository(database, &evalRepo, &sourcecodeRepo)
	checkerS := checkerimpl.NewCheckerService()
	interactorS := interactorimpl.NewInteractorService()
	redis := ri.NewRedisSubmissionRepository(rdb)
	isolateS, _ := ii.NewIsolateService()
	judgeSvc := ji.NewJudgeServiceImpl(pool, problemSvc, evalRepo, checkerS, interactorS, redis, submissionRepo, sourcecodeRepo, isolateS)
	contestSvc, err := contestserviceimpl.NewContestService()
	if err != nil {
		config.GetLogger().Panic().Err(err).Msg("Can't not create contest service")
	}

	redisRepo := ri.NewRedisSubmissionRepository(rdb)
	return &appCtx{
		database,
		pool,
		rdb,
		submissionRepo,
		sourcecodeRepo,
		problemSvc,
		evalRepo,
		checkerS,
		interactorS,
		redisRepo,
		isolateS,
		judgeSvc,
		contestSvc,
	}
}
