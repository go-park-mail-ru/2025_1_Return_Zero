package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/redis"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
)

func main() {
	logger, err := loggerPkg.NewZapLogger()
	if err != nil {
		logger.Error("Error creating logger:", zap.Error(err))
		return
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config:", zap.Error(err))
		return
	}

	port := fmt.Sprintf(":%d", cfg.Services.AuthService.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Can't start auth service:", zap.Error(err))
		return
	}
	defer conn.Close()

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "auth_service")

	accessInterceptor := interceptors.NewAccessInterceptor(logger, metrics)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
	)

	redisPool := redis.NewRedisPool(cfg.Redis)
	defer redisPool.Close()

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.AuthPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.AuthPort))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	authRepository := repository.NewAuthRedisRepository(redisPool, metrics)
	authUsecase := usecase.NewAuthUsecase(authRepository)
	authService := delivery.NewAuthService(authUsecase)
	authProto.RegisterAuthServiceServer(server, authService)

	logger.Info("Auth service started on port %s...", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting auth service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("Auth service stopped")
}
