package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/s3"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/usecase"
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
	defer logger.Sync()

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
	defer conn.Close()

	accessInterceptor := interceptors.NewAccessInterceptor(logger)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
	)

	postgresPool, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		logger.Error("Error connecting to postgres:", zap.Error(err))
		return
	}
	defer postgresPool.Close()

	fmt.Println("config ", cfg.S3.S3ImagesBucket)
	s3, err := s3.InitS3(cfg.S3)
	if err != nil {
		logger.Error("Error initializing S3:", zap.Error(err))
		return
	}

	userRepository := repository.NewUserPostgresRepository(postgresPool)
	userS3Repository := repository.NewS3Repository(s3, cfg.S3.S3ImagesBucket)
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
