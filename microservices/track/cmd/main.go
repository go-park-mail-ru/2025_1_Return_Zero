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
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
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

	port := fmt.Sprintf(":%d", cfg.Services.TrackService.Port)
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
	metrics := metrics.NewMetrics(reg, "track_service")

	accessInterceptor := interceptors.NewAccessInterceptor(logger, metrics)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
		grpc.MaxRecvMsgSize(50*1024*1024), // 50 MB
		grpc.MaxSendMsgSize(50*1024*1024), // 50 MB
	)

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.TrackPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.TrackPort))
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

	trackRepository := repository.NewTrackPostgresRepository(postgresPool, metrics)
	trackS3Repository := repository.NewTrackS3Repository(s3, cfg.S3.S3TracksBucket, cfg.S3.S3Duration, metrics)
	trackUsecase := usecase.NewTrackUsecase(trackRepository, trackS3Repository)
	trackService := delivery.NewTrackService(trackUsecase)
	trackProto.RegisterTrackServiceServer(server, trackService)

	logger.Info("Track service started", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting track service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("Track service stopped")
}
