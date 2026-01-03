package components

import (
	repository "contest/src/domain/repository/contest/impl"
	contestusecase "contest/src/usecase/contest"
	contestinteractor "contest/src/usecase/contest/interactor"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppContext interface {
	GetMainDbConnection() *mongo.Database
	GetRedis() *redis.Client
	GetContestInteractor() contestusecase.ContestInteractor
}

type appCtx struct {
	database          *mongo.Database
	rdb               *redis.Client
	contestInteractor contestusecase.ContestInteractor
}

func (ctx *appCtx) GetMainDbConnection() *mongo.Database { return ctx.database }
func (ctx *appCtx) GetRedis() *redis.Client              { return ctx.rdb }
func (ctx *appCtx) GetContestInteractor() contestusecase.ContestInteractor {
	return ctx.contestInteractor
}

func NewAppContext(
	database *mongo.Database,
	rdb *redis.Client,
) *appCtx {
	contestRepo := repository.NewContestRepository(database)
	contestInteractor := contestinteractor.NewContestInteractor(contestRepo)

	return &appCtx{
		database:          database,
		rdb:               rdb,
		contestInteractor: contestInteractor,
	}
}
