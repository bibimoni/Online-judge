package components

import (
	contestrepo "contest/src/domain/repository/contest"
	repository "contest/src/domain/repository/contest/impl"
	"contest/src/infrastructure/config"
	contestservice "contest/src/service/contest"
	contestserviceimpl "contest/src/service/contest/impl"
	problemserviceimpl "contest/src/service/problem/impl"
	contestusecase "contest/src/usecase/contest"
	contestinteractor "contest/src/usecase/contest/interactor"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppContext interface {
	GetMainDbConnection() *mongo.Database
	GetRedis() *redis.Client
	GetContestInteractor() contestusecase.ContestInteractor
	GetContestRepository() contestrepo.ContestRepository
	GetContestService() contestservice.ContestService
}

type appCtx struct {
	database          *mongo.Database
	rdb               *redis.Client
	contestInteractor contestusecase.ContestInteractor
	contestRepo       contestrepo.ContestRepository
	contestService    contestservice.ContestService
}

func (ctx *appCtx) GetMainDbConnection() *mongo.Database { return ctx.database }
func (ctx *appCtx) GetRedis() *redis.Client              { return ctx.rdb }
func (ctx *appCtx) GetContestInteractor() contestusecase.ContestInteractor {
	return ctx.contestInteractor
}
func (ctx *appCtx) GetContestRepository() contestrepo.ContestRepository { return ctx.contestRepo }
func (ctx *appCtx) GetContestService() contestservice.ContestService    { return ctx.contestService }

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

	return &appCtx{
		database:          database,
		rdb:               rdb,
		contestInteractor: contestInteractor,
		contestRepo:       contestRepo,
		contestService:    contestService,
	}
}
