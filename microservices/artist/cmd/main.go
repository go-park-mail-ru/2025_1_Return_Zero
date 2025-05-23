package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/interceptors"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
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

	port := fmt.Sprintf(":%d", cfg.Services.ArtistService.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Can't start artist service:", zap.Error(err))
		return
	}
	defer conn.Close()

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "artist_service")

	accessInterceptor := interceptors.NewAccessInterceptor(logger, metrics)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor.UnaryServerInterceptor()),
	)

	postgresPool, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		logger.Error("Error connecting to postgres:", zap.Error(err))
		return
	}
	defer postgresPool.Close()

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.ArtistPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.ArtistPort))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	artistRepository := repository.NewArtistPostgresRepository(postgresPool, metrics)
	artistUsecase := usecase.NewArtistUsecase(artistRepository)
	artistService := delivery.NewArtistService(artistUsecase)
	artistProto.RegisterArtistServiceServer(server, artistService)

	logger.Info("Artist service started", zap.String("port", port))

	err = server.Serve(conn)
	if err != nil {
		logger.Fatal("Error starting artist service:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("Artist service stopped")
}
