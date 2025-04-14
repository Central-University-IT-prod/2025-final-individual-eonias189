package main

import (
	"advertising/advertising-service/internal/config"
	"advertising/advertising-service/internal/repo/minio"
	"advertising/advertising-service/internal/repo/postgres"
	"advertising/advertising-service/internal/repo/redis"
	"advertising/advertising-service/internal/service"
	"advertising/advertising-service/internal/transport/rest/v1"
	"advertising/advertising-service/internal/transport/rest/v1/handlers"
	"advertising/pkg/logger"
	minio_helper "advertising/pkg/minio"
	"advertising/pkg/openai"
	pg_helper "advertising/pkg/postgres"
	redis_helper "advertising/pkg/redis"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	shutdownTimeout = time.Second * 5
)

func main() {
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		log.Fatal("get config", err)
	}

	l, err := logger.Get(cfg.LogLevel)
	if err != nil {
		log.Fatal("get logger", err)
	}

	db, err := pg_helper.Connect(ctx, cfg.PostgresConfig)
	if err != nil {
		l.Fatal("connect to postgres", zap.Error(err))
	}

	rdb, err := redis_helper.Connect(ctx, cfg.RedisConfig)
	if err != nil {
		l.Fatal("connect to redis", zap.Error(err))
	}

	minioCli, err := minio_helper.Connect(cfg.MinioConfig)
	if err != nil {
		l.Fatal("connect to minio", zap.Error(err))
	}

	chat := openai.NewChat(cfg.OpenAIConfig)

	timeRepo := redis.NewTimeRepo(rdb)
	clientsRepo := postgres.NewClientRepo(db)
	advertisersRepo := postgres.NewAdvertiserRepo(db)
	mlScoreRepo := postgres.NewMlScoresRepo(db)
	campaignsRepo := postgres.NewCampaignsRepo(db)
	adsRepo := postgres.NewAdsRepo(db)
	clientActionsRepo := postgres.NewClientActionsRepo(db)
	statsRepo := postgres.NewStatsRepo(db)
	staticRepo := minio.NewStaticRepo(minioCli, cfg.StaticBucket)

	timeService := service.NewTimeService(timeRepo)
	advertisersService := service.NewAdvertisersService(advertisersRepo, mlScoreRepo)
	campaignsService := service.NewCampaignsService(campaignsRepo, advertisersRepo, timeRepo, staticRepo, cfg.StaticBaseUrl)
	adsService := service.NewAdsService(adsRepo, clientsRepo, campaignsRepo, clientActionsRepo, timeRepo)
	statsService := service.NewStatsService(statsRepo, campaignsRepo, advertisersRepo)
	aiService := service.NewAIService(chat)

	adsHandler := handlers.NewAdsHandler(adsService)
	advertisersHandler := handlers.NewAdvertisersHandler(advertisersService)
	campaignsHandler := handlers.NewCampaignsHandler(campaignsService)
	clietnsHandler := handlers.NewClientsHandler(clientsRepo)
	statisticsHandler := handlers.NewStatsHandler(statsService)
	timeHandler := handlers.NewTimeHandler(timeService)
	staticHandler := handlers.NewStaticHandler(staticRepo)
	aiHandler := handlers.NewAIHandler(aiService)

	handler := rest.NewHandler(
		adsHandler, advertisersHandler, campaignsHandler,
		clietnsHandler, statisticsHandler, timeHandler,
		aiHandler,
	)

	server, err := rest.NewServer(handler, staticHandler, l)
	if err != nil {
		l.Fatal("get logger", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		l.Info("starting server at port", zap.Int("port", cfg.ServerPort))
		err := server.Start(ctx, cfg.ServerPort)
		if err != nil && err != http.ErrServerClosed {
			l.Error("start server", zap.Error(err))
		}
	}()

	<-sigCh
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		l.Error("shutdown server", zap.Error(err))
	}

	l.Info("server stopped")
}
