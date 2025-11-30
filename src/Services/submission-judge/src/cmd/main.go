package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appctx "github.com/bibimoni/Online-judge/submission-judge/src/components"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/database"
	"github.com/bibimoni/Online-judge/submission-judge/src/router"
	pi "github.com/bibimoni/Online-judge/submission-judge/src/service/pool/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	si "github.com/bibimoni/Online-judge/submission-judge/src/service/store/impl"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Can't load config")
	}
	log := config.NewLogger(cfg.LogLevel)

	client, err := database.GetMongoDbClient(cfg.Database.Uri)
	if err != nil {
		log.Fatal().Err(err).Msgf("Can't not load mongoDB")
	}
	defer client.Disconnect(context.Background())

	pool, err := pi.NewPoolSerivce()

	if err != nil {
		log.Fatal().Err(err).Msgf("Can't initialize new pool service")
	}

	redis, err := database.GetRedisClient()
	if err != nil {
		log.Fatal().Err(err).Msgf("Can't initialize redis client")
	}

	appCtx := appctx.NewAppContext(client.Database(cfg.Database.Name), pool, redis)

	store.DefaultStore = si.NewStoreWithDefaultLangs()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	v1 := r.Group("/api/v1")

	gin.SetMode(gin.DebugMode)

	router.RegisterRouter(v1, appCtx)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port

	log.Info().Msgf("Submission-Judge server is listening on: %s", serverAddr)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("Failed to start Submission-Judge server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
