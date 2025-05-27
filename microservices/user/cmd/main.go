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
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
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

	port := fmt.Sprintf(":%d", cfg.Services.UserService.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Can't start user service:", zap.Error(err))
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection:", zap.Error(err))
		}
	}()

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "user_service")

	accessInterceptor := interceptors.NewAccessInterceptor(logger, metrics)

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.UserPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.UserPort))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
	)

	postgresPool, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		logger.Error("Error connecting to postgres:", zap.Error(err))
		return
	}
	defer func() {
		if err := postgresPool.Close(); err != nil {
			logger.Error("Error closing postgres pool:", zap.Error(err))
		}
	}()

	fmt.Println("config ", cfg.S3.S3ImagesBucket)
	s3, err := s3.InitS3(cfg.S3)
	if err != nil {
		logger.Error("Error initializing S3:", zap.Error(err))
		return
	}

	userRepository := repository.NewUserPostgresRepository(postgresPool, metrics)
	userS3Repository := repository.NewS3Repository(s3, cfg.S3.S3ImagesBucket, metrics)
	userUsecase := usecase.NewUserUsecase(userRepository, userS3Repository)
	userService := delivery.NewUserService(userUsecase)
	userProto.RegisterUserServiceServer(server, userService)

	logger.Info("User service started on port %s...", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting user service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("User service stopped")
}
