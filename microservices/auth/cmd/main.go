package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/redis"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/usecase"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
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

	accessInterceptor := interceptors.NewAccessInterceptor(logger)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
	)

	redisPool := redis.NewRedisPool(cfg.Redis)
	defer redisPool.Close()

	authRepository := repository.NewAuthRedisRepository(redisPool)
	authUsecase := usecase.NewAuthUsecase(authRepository)
	authService := delivery.NewAuthService(authUsecase)
	authProto.RegisterAuthServiceServer(server, authService)

	logger.Info("Auth service started on port %s...", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting artist service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("Auth service stopped")
}
