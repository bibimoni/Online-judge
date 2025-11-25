package appctx

import (
	poolservice "github.com/bibimoni/Online-judge/submission-judge/src/service/pool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppContext interface {
	GetMainDbConnection() *mongo.Database
	GetPool() poolservice.PoolService
	GetRedis() *redis.Client
}

type appCtx struct {
	database *mongo.Database
	pool     poolservice.PoolService
	rdb      *redis.Client
}

func NewAppContext(
	database *mongo.Database,
	pool poolservice.PoolService,
	rdb *redis.Client,
) *appCtx {
	return &appCtx{
		database,
		pool,
		rdb,
	}
}

func (ctx *appCtx) GetMainDbConnection() *mongo.Database {
	return ctx.database
}

func (ctx *appCtx) GetPool() poolservice.PoolService {
	return ctx.pool
}

func (ctx *appCtx) GetRedis() *redis.Client {
	return ctx.rdb
}
