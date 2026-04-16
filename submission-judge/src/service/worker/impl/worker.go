package workerimpl

import (
	"context"
	"sync"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	evaluationRepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation"
	redisRepo "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	contestservice "github.com/bibimoni/Online-judge/submission-judge/src/service/contest"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/judge"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/worker"
)

type WorkerServiceImpl struct {
	redisRepo      redisRepo.RedisSubmissionRepository
	judgeService   judge.JudgeService
	problemService problem.ProblemService
	quit           chan struct{}
	wg             sync.WaitGroup
	workerCount    int
	contestService contestservice.ContestService
	evalRepo       evaluationRepo.EvaluationRepository
}

func NewWorkerServiceImpl(
	redis redisRepo.RedisSubmissionRepository,
	judgeService judge.JudgeService,
	problemService problem.ProblemService,
	workerCount int,
	contestService contestservice.ContestService,
	evalRepo evaluationRepo.EvaluationRepository,
) *WorkerServiceImpl {
	return &WorkerServiceImpl{
		redisRepo:      redis,
		judgeService:   judgeService,
		problemService: problemService,
		quit:           make(chan struct{}),
		workerCount:    workerCount,
		contestService: contestService,
		evalRepo:       evalRepo,
	}
}
func NewWorkerService(
	redis redisRepo.RedisSubmissionRepository,
	judgeService judge.JudgeService,
	problemService problem.ProblemService,
	workerCount int,
	contestService contestservice.ContestService,
	evalRepo evaluationRepo.EvaluationRepository,
) worker.WorkerService {
	return NewWorkerServiceImpl(redis, judgeService, problemService, workerCount, contestService, evalRepo)
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

			ttl := time.Duration(problemInfo.TimeLimit*problemInfo.TestNum*2) * time.Millisecond
			status, err := ws.redisRepo.SetNX(ctx, req.SubmissionId, "locked", ttl)
			if err != nil {
				config.GetLogger().Error().Err(err).Msg("Error getting submission status from redis")
				continue
			}
			if !status {
				config.GetLogger().Info().Msgf("Submission %s is already being processed, skipping...", req.SubmissionId)
				continue
			}
			config.GetLogger().Info().Msgf("Locking submission %s for %v", req.SubmissionId, ttl)

			err = ws.judgeService.JudgeStart(ctx, lang, req, problemInfo)
			if err != nil {
				config.GetLogger().Error().Err(err).Msg("Error processing submission")
			}

			// if this is from a contest (upsert!)
			if req.ContestId != "" {
				config.GetLogger().Info().Msgf("Ingesting contest submission for submission id: %s", req.SubmissionId)
				eval, err := ws.evalRepo.GetEvalBySubmissionIdNoBson(ctx, req.SubmissionId)
				if err != nil {
					config.GetLogger().Error().Err(err).Msg("Error fetching submission after judging")
					continue
				}
				if domain.SubmissionType(problemInfo.ScoringMode) == domain.ICPC {
					ws.contestService.IngestContestSubmission(ctx, req.SubmissionId, eval.Verdict, float64(eval.Points))
				} else if domain.SubmissionType(problemInfo.ScoringMode) == domain.IOI {
					ws.contestService.IngestContestSubmission(ctx, req.SubmissionId, eval.Verdict, float64(eval.Score))
				} else {
					// placeholder
					ws.contestService.IngestContestSubmission(ctx, req.SubmissionId, eval.Verdict, float64(eval.Points))
				}
			}
		}
	}
}
