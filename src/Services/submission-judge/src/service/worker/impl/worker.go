package workerimpl

import (
	"context"
	redisRepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/judge"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/worker"
	"sync"
	"time"
)

type WorkerServiceImpl struct {
	redisRepo      redisRepo.RedisSubmissionRepository
	judgeService   judge.JudgeService
	problemService problem.ProblemService
	quit           chan struct{}
	wg             sync.WaitGroup
	workerCount    int
}

func NewWorkerServiceImpl(
	redis redisRepo.RedisSubmissionRepository,
	judgeService judge.JudgeService,
	problemService problem.ProblemService,
	workerCount int,
) *WorkerServiceImpl {
	return &WorkerServiceImpl{
		redisRepo:      redis,
		judgeService:   judgeService,
		problemService: problemService,
		quit:           make(chan struct{}),
		workerCount:    workerCount,
	}
}
func NewWorkerService(
	redis redisRepo.RedisSubmissionRepository,
	judgeService judge.JudgeService,
	problemService problem.ProblemService,
	workerCount int,
) worker.WorkerService {
	return NewWorkerServiceImpl(redis, judgeService, problemService, workerCount)
}
func (ws *WorkerServiceImpl) Start() {
	for range ws.workerCount {
		ws.wg.Add(1)
		go ws.worker()
	}
}
func (ws *WorkerServiceImpl) Stop() {
	close(ws.quit)
	ws.wg.Wait()
}
func (ws *WorkerServiceImpl) worker() {
	defer ws.wg.Done()
	for {
		select {
		case <-ws.quit:
			return
		default:
			ctx := context.Background()
			req, err := ws.redisRepo.PopSubmissionJob(ctx)
			if err != nil {
				config.GetLogger().Error().Err(err).Msg("Redis connection error, retrying...")
				time.Sleep(time.Second)
				continue
			}
			lang, err := store.DefaultStore.Get(req.LanguageId)
			if err != nil {
				config.GetLogger().Error().Err(err).Msg("Unknown language")
				continue
			}
			problemInfo, err := ws.problemService.Get(ctx, req.ProblemId)
			if err != nil {
				config.GetLogger().Error().Err(err).Msg("Unknown problem")
				continue
			}
			err = ws.judgeService.JudgeStart(ctx, lang, req, problemInfo)
			if err != nil {
				config.GetLogger().Error().Err(err).Msg("Error processing submission")
			}
		}
	}
}
