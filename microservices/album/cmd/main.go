package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/s3"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "album_service")

	accessInterceptor := interceptors.NewAccessInterceptor(logger, metrics)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
		grpc.MaxRecvMsgSize(50*1024*1024), // 50 MB
		grpc.MaxSendMsgSize(50*1024*1024), // 50 MB
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

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.AlbumPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.AlbumPort))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	albumRepository := repository.NewAlbumPostgresRepository(postgresPool, metrics)
	s3Repository := repository.NewS3Repository(s3, cfg.S3.S3ImagesBucket, metrics)
	albumUsecase := usecase.NewAlbumUsecase(albumRepository, s3Repository)
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
