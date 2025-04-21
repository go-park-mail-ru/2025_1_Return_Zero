package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	grpc "github.com/go-park-mail-ru/2025_1_Return_Zero/init/microservices"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/redis"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/s3"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	albumHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/delivery/http"
	albumRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/repository"
	albumUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/usecase"
	artistHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/delivery/http"
	artistRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/repository"
	artistUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/usecase"
	authRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	trackHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/delivery/http"
	trackRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/repository"
	trackUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/usecase"
	trackFileRepo "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/trackFile/repository"
	userHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/repository"
	userUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/usecase"
	userFileRepo "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile/repository"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host returnzero.ru
// @BasePath /api/v1
func main() {
	logger, err := logger.NewZapLogger()
	if err != nil {
		logger.Error("Error creating logger:", zap.Error(err))
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config:", zap.Error(err))
		return
	}

	redisPool := redis.NewRedisPool(cfg.Redis)
	defer redisPool.Close()

	postgresConn, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		logger.Error("Error connecting to Postgres:", zap.Error(err))
		return
	}
	defer postgresConn.Close()

	s3, err := s3.InitS3(cfg.S3)
	if err != nil {
		logger.Error("Error initializing S3:", zap.Error(err))
		return
	}

	r := mux.NewRouter()
	logger.Info("Server starting on port %s...", zap.String("port", fmt.Sprintf(":%d", cfg.Port)))

	r.PathPrefix("/api/v1/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/v1/docs/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	clients, err := grpc.InitGrpc(&cfg.Services, logger)
	if err != nil {
		logger.Error("Error initializing gRPC clients:", zap.Error(err))
		return
	}

	artistClient := artistProto.NewArtistServiceClient(clients.ArtistClient)

	newUserUsecase := userUsecase.NewUserUsecase(userRepository.NewUserPostgresRepository(postgresConn), authRepository.NewAuthRedisRepository(redisPool), userFileRepo.NewS3Repository(s3, cfg.S3.S3ImagesBucket))

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Auth(newUserUsecase))
	r.Use(middleware.CorsMiddleware(cfg.Cors))
	r.Use(middleware.CSRFMiddleware(cfg.CSRF))

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackRepository.NewTrackPostgresRepository(postgresConn), artistRepository.NewArtistPostgresRepository(postgresConn), albumRepository.NewAlbumPostgresRepository(postgresConn), trackFileRepo.NewS3Repository(s3, cfg.S3.S3TracksBucket, cfg.S3.S3Duration), userRepository.NewUserPostgresRepository(postgresConn)), cfg)
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumRepository.NewAlbumPostgresRepository(postgresConn), artistRepository.NewArtistPostgresRepository(postgresConn)), cfg)
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(&artistClient), cfg)
	userHandler := userHttp.NewUserHandler(newUserUsecase)

	r.HandleFunc("/api/v1/tracks", trackHandler.GetAllTracks).Methods("GET")
	r.HandleFunc("/api/v1/tracks/{id}", trackHandler.GetTrackByID).Methods("GET")
	r.HandleFunc("/api/v1/tracks/{id}/stream", trackHandler.CreateStream).Methods("POST")
	r.HandleFunc("/api/v1/streams/{id}", trackHandler.UpdateStreamDuration).Methods("PUT", "PATCH")

	r.HandleFunc("/api/v1/albums", albumHandler.GetAllAlbums).Methods("GET")

	r.HandleFunc("/api/v1/artists", artistHandler.GetAllArtists).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id}", artistHandler.GetArtistByID).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id}/tracks", trackHandler.GetTracksByArtistID).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id}/albums", albumHandler.GetAlbumsByArtistID).Methods("GET")

	r.HandleFunc("/api/v1/auth/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/logout", userHandler.Logout).Methods("POST")
	r.HandleFunc("/api/v1/auth/check", userHandler.CheckUser).Methods("GET")

	r.HandleFunc("/api/v1/user/me/avatar", userHandler.UploadAvatar).Methods("POST")
	r.HandleFunc("/api/v1/user/me", userHandler.ChangeUserData).Methods("PUT")
	r.HandleFunc("/api/v1/user/me", userHandler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/v1/user/{username}", userHandler.GetUserData).Methods("GET")

	r.HandleFunc("/api/v1/user/{username}/history", trackHandler.GetLastListenedTracks).Methods("GET")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("Error starting server:", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	logger.Info("Composer server stopped")
}
