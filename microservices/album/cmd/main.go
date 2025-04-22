package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
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

	port := fmt.Sprintf(":%d", cfg.Services.AlbumService.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Can't start album service:", zap.Error(err))
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

	albumRepository := repository.NewAlbumPostgresRepository(postgresPool)
	albumUsecase := usecase.NewAlbumUsecase(albumRepository)
	albumService := delivery.NewAlbumService(albumUsecase)
	albumProto.RegisterAlbumServiceServer(server, albumService)

	logger.Info("Album service started", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting album service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("Album service stopped")
}
