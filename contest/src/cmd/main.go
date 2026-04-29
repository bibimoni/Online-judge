package main

import (
	"contest/services/scheduler"
	"contest/src/components"
	"contest/src/infrastructure/config"
	"contest/src/infrastructure/database"
	"contest/src/router"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	redis, err := database.GetRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Fatal().Err(err).Msgf("Can't initialize redis client")
	}

	appCtx := components.NewAppContext(client.Database(cfg.Database.Name), redis)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	v1 := r.Group("/api/v1")
	gin.SetMode(gin.DebugMode)
	router.RegisterRouter(v1, appCtx)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Info().Msgf("Contest server is listening on: %s", serverAddr)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		scheduler.Scheduler(redis, appCtx.GetContestRepository())
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("Failed to start Contest server: %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exiting")
}
