package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/s3"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
)

func main() {
	logger, err := loggerPkg.NewZapLogger()
	if err != nil {
		logger.Error("Error creating logger:", zap.Error(err))
		return
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Error syncing logger:", zap.Error(err))
		}
	}()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config:", zap.Error(err))
		return
	}

	port := fmt.Sprintf(":%d", cfg.Services.PlaylistService.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Can't start album service:", zap.Error(err))
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection:", zap.Error(err))
		}
	}()

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "playlist_service")

	accessInterceptor := interceptors.NewAccessInterceptor(logger, metrics)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
	)

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.PlaylistPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.PlaylistPort))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	postgresPool, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		logger.Error("Error connecting to postgres:", zap.Error(err))
		return
	}
	defer func() {
		if err := postgresPool.Close(); err != nil {
			logger.Error("Error closing Postgres:", zap.Error(err))
		}
	}()

	s3, err := s3.InitS3(cfg.S3)
	if err != nil {
		logger.Error("Error initializing S3:", zap.Error(err))
		return
	}

	playlistRepository := repository.NewPlaylistPostgresRepository(postgresPool, metrics)
	playlistS3Repository := repository.NewPlaylistS3Repository(s3, cfg.S3.S3ImagesBucket)
	playlistUsecase := usecase.NewPlaylistUsecase(playlistRepository, playlistS3Repository)
	playlistService := delivery.NewPlaylistService(playlistUsecase, playlistS3Repository)
	playlistProto.RegisterPlaylistServiceServer(server, playlistService)

	logger.Info("Playlist service started", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting playlist service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("Playlist service stopped")
}
