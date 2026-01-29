package components

import (
	contestrepo "contest/src/domain/repository/contest"
	contestsubmissionrepo "contest/src/domain/repository/contest-submission"
	contestsubmissionrepoimpl "contest/src/domain/repository/contest-submission/impl"
	repository "contest/src/domain/repository/contest/impl"
	"contest/src/infrastructure/config"
	contestservice "contest/src/service/contest"
	contestsubmissionservice "contest/src/service/contest-submission"
	contestsubmissionserviceimpl "contest/src/service/contest-submission/impl"
	contestserviceimpl "contest/src/service/contest/impl"
	problemserviceimpl "contest/src/service/problem/impl"
	contestusecase "contest/src/usecase/contest"
	contestsubmissionusecase "contest/src/usecase/contest-submission"
	contestsubmissioninteractor "contest/src/usecase/contest-submission/interactor"
	contestinteractor "contest/src/usecase/contest/interactor"
	contestantusecase "contest/src/usecase/contestant"
	contestantinteractor "contest/src/usecase/contestant/interactor"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppContext interface {
	GetMainDbConnection() *mongo.Database
	GetRedis() *redis.Client
	GetContestInteractor() contestusecase.ContestInteractor
	GetContestRepository() contestrepo.ContestRepository
	GetContestService() contestservice.ContestService
	GetContestantInteractor() contestantusecase.ContestantInteractor
	GetContestSubmissionInteractor() contestsubmissionusecase.ContestSubmissionInteractor
}

type appCtx struct {
	database                    *mongo.Database
	rdb                         *redis.Client
	contestInteractor           contestusecase.ContestInteractor
	contestRepo                 contestrepo.ContestRepository
	contestService              contestservice.ContestService
	contestantInteractor        contestantusecase.ContestantInteractor
	contestsubmissionRepo       contestsubmissionrepo.ContestSubmissionRepository
	contestsubmissionService    contestsubmissionservice.ContestSubmissionService
	contestsubmissionInteractor contestsubmissionusecase.ContestSubmissionInteractor
}

func (ctx *appCtx) GetMainDbConnection() *mongo.Database                { return ctx.database }
func (ctx *appCtx) GetRedis() *redis.Client                             { return ctx.rdb }
func (ctx *appCtx) GetContestRepository() contestrepo.ContestRepository { return ctx.contestRepo }
func (ctx *appCtx) GetContestService() contestservice.ContestService    { return ctx.contestService }
func (ctx *appCtx) GetContestInteractor() contestusecase.ContestInteractor {
	return ctx.contestInteractor
}
func (ctx *appCtx) GetContestantInteractor() contestantusecase.ContestantInteractor {
	return ctx.contestantInteractor
}
func (ctx *appCtx) GetContestSubmissionInteractor() contestsubmissionusecase.ContestSubmissionInteractor {
	return ctx.contestsubmissionInteractor
}
func NewAppContext(
	database *mongo.Database,
	rdb *redis.Client,
) *appCtx {
	contestRepo := repository.NewContestRepository(database)
	problemService, err := problemserviceimpl.NewProblemService()
	if err != nil {
		config.GetLogger().Fatal().Err(err).Msg("failed to create problem service")
	}
	contestService := contestserviceimpl.NewContestService(contestRepo, problemService)
	contestInteractor := contestinteractor.NewContestInteractor(contestRepo, contestService)
	contestsubmissionRepo := contestsubmissionrepoimpl.NewContestSubmissionRepository(database)
	contestsubmissionservice, err := contestsubmissionserviceimpl.NewContestSubmissionService(
		contestRepo,
		contestsubmissionRepo,
	)
	if err != nil {
		config.GetLogger().Fatal().Err(err).Msg("failed to create contestsubmission service")
	}
	contestantInteractor := contestantinteractor.NewContestantInteractor(
		contestService,
		contestRepo,
		contestsubmissionservice,
	)
	contestsubmissionInteractor := contestsubmissioninteractor.NewContestSubmissionInteractor(
		contestsubmissionservice,
	)

	return &appCtx{
		database:                    database,
		rdb:                         rdb,
		contestInteractor:           contestInteractor,
		contestRepo:                 contestRepo,
		contestService:              contestService,
		contestantInteractor:        contestantInteractor,
		contestsubmissionRepo:       contestsubmissionRepo,
		contestsubmissionService:    contestsubmissionservice,
		contestsubmissionInteractor: contestsubmissionInteractor,
	}
}
