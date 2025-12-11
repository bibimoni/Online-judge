package components

import (
	repository "contest/src/domain/repository/contest/impl"
	createcontest "contest/src/usecase/create_contest"
	createcontestimpl "contest/src/usecase/create_contest/impl"
	editcontest "contest/src/usecase/edit_contest"
	editcontestimpl "contest/src/usecase/edit_contest/impl"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppContext interface {
	GetMainDbConnection() *mongo.Database
	GetRedis() *redis.Client
	GetCreateContestInteractor() createcontest.CreateContestInteractor
	GetEditContestInteractor() editcontest.EditContestInteractor
}

type appCtx struct {
	database                *mongo.Database
	rdb                     *redis.Client
	createContestInteractor createcontest.CreateContestInteractor
	editContestInteractor   editcontest.EditContestInteractor
}

func (ctx *appCtx) GetMainDbConnection() *mongo.Database { return ctx.database }
func (ctx *appCtx) GetRedis() *redis.Client              { return ctx.rdb }
func (ctx *appCtx) GetCreateContestInteractor() createcontest.CreateContestInteractor {
	return ctx.createContestInteractor
}
func (ctx *appCtx) GetEditContestInteractor() editcontest.EditContestInteractor {
	return ctx.editContestInteractor
}

func NewAppContext(
	database *mongo.Database,
	rdb *redis.Client,
) *appCtx {
	contestRepo := repository.NewContestRepository(database)

	createContestInteractor := createcontestimpl.NewCreateContestInteractor(contestRepo)
	editContestInteractor := editcontestimpl.NewEditContestInteractor(contestRepo)

	return &appCtx{
		database:                database,
		rdb:                     rdb,
		createContestInteractor: createContestInteractor,
		editContestInteractor:   editContestInteractor,
	}
}
