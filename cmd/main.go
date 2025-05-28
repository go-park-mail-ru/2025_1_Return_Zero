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
	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	grpc "github.com/go-park-mail-ru/2025_1_Return_Zero/init/microservices"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/redis"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	albumHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/delivery/http"
	albumUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/usecase"
	artistHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/delivery/http"
	artistUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	jamHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam/delivery/http"
	jamRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam/repository"
	jamUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam/usecase"
	playlistHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/playlist/delivery/http"
	playlistUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/playlist/usecase"
	trackHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/delivery/http"
	trackUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/usecase"
	userHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/delivery/http"
	userUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/usecase"

	labelHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/delivery/http"
	labelRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/repository"
	labelUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/usecase"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/metrics"
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
	defer func() {
		if err := redisPool.Close(); err != nil {
			logger.Error("Error closing Redis:", zap.Error(err))
		}
	}()

	postgresConn, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		logger.Error("Error connecting to Postgres:", zap.Error(err))
		return
	}
	defer func() {
		if err := postgresConn.Close(); err != nil {
			logger.Error("Error closing Postgres:", zap.Error(err))
		}
	}()

	err = postgres.RunMigrations(cfg.Postgres)
	if err != nil {
		logger.Error("Error running migrations:", zap.Error(err))
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

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "api")
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf(":%d", cfg.Prometheus.ApiPort)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.Prometheus.ApiPort))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.Error(err))
		}
	}()

	artistClient := artistProto.NewArtistServiceClient(clients.ArtistClient)
	albumClient := albumProto.NewAlbumServiceClient(clients.AlbumClient)
	trackClient := trackProto.NewTrackServiceClient(clients.TrackClient)
	playlistClient := playlistProto.NewPlaylistServiceClient(clients.PlaylistClient)
	authClient := authProto.NewAuthServiceClient(clients.AuthClient)
	userClient := userProto.NewUserServiceClient(clients.UserClient)

	labelRepository := labelRepository.NewLabelPostgresRepository(postgresConn)
	labelUsecase := labelUsecase.NewLabelUsecase(labelRepository, userClient, artistClient, albumClient, trackClient)
	labelHandler := labelHttp.NewLabelHandler(labelUsecase, cfg)

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Auth(&authClient, &userClient))
	r.Use(middleware.CorsMiddleware(cfg.Cors))
	// r.Use(middleware.CSRFMiddleware(cfg.CSRF))
	r.Use(middleware.MetricsMiddleware(metrics))

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackClient, artistClient, albumClient, playlistClient, userClient), cfg)
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumClient, artistClient), cfg)
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(artistClient, userClient), cfg)
	userHandler := userHttp.NewUserHandler(userUsecase.NewUserUsecase(&userClient, &authClient, &artistClient, &trackClient, &playlistClient))
	playlistHandler := playlistHttp.NewPlaylistHandler(playlistUsecase.NewUsecase(&playlistClient, &userClient), cfg)
	jamHandler := jamHttp.NewJamHandler(jamUsecase.NewUsecase(jamRepository.NewJamRedisRepository(redisPool), userClient), cfg)

	r.HandleFunc("/api/v1/tracks", trackHandler.GetAllTracks).Methods("GET")
	r.HandleFunc("/api/v1/tracks/{id:[0-9]+}", trackHandler.GetTrackByID).Methods("GET")
	r.HandleFunc("/api/v1/tracks/{id:[0-9]+}/stream", trackHandler.CreateStream).Methods("POST")
	r.HandleFunc("/api/v1/tracks/{id:[0-9]+}/like", trackHandler.LikeTrack).Methods("POST")
	r.HandleFunc("/api/v1/tracks/search", trackHandler.SearchTracks).Methods("GET")
	r.HandleFunc("/api/v1/streams/{id:[0-9]+}", trackHandler.UpdateStreamDuration).Methods("PUT", "PATCH")
	r.HandleFunc("/api/v1/selection/{selection}", trackHandler.GetSelectionTracks).Methods("GET")

	r.HandleFunc("/api/v1/albums", albumHandler.GetAllAlbums).Methods("GET")
	r.HandleFunc("/api/v1/albums/{id:[0-9]+}", albumHandler.GetAlbumByID).Methods("GET")
	r.HandleFunc("/api/v1/albums/search", albumHandler.SearchAlbums).Methods("GET")
	r.HandleFunc("/api/v1/albums/{id:[0-9]+}/tracks", trackHandler.GetTracksByAlbumID).Methods("GET")
	r.HandleFunc("/api/v1/albums/{id:[0-9]+}/like", albumHandler.LikeAlbum).Methods("POST")

	r.HandleFunc("/api/v1/artists", artistHandler.GetAllArtists).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id:[0-9]+}", artistHandler.GetArtistByID).Methods("GET")
	r.HandleFunc("/api/v1/artists/search", artistHandler.SearchArtists).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id:[0-9]+}/tracks", trackHandler.GetTracksByArtistID).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id:[0-9]+}/albums", albumHandler.GetAlbumsByArtistID).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id:[0-9]+}/like", artistHandler.LikeArtist).Methods("POST")

	r.HandleFunc("/api/v1/playlists", playlistHandler.CreatePlaylist).Methods("POST")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}", playlistHandler.UpdatePlaylist).Methods("PUT")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}", playlistHandler.RemovePlaylist).Methods("DELETE")
	r.HandleFunc("/api/v1/playlists/to-add", playlistHandler.GetPlaylistsToAdd).Methods("GET")
	r.HandleFunc("/api/v1/playlists/me", playlistHandler.GetCombinedPlaylistsForCurrentUser).Methods("GET")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}/tracks", playlistHandler.AddTrackToPlaylist).Methods("POST")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}/tracks/{trackId:[0-9]+}", playlistHandler.RemoveTrackFromPlaylist).Methods("DELETE")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}/tracks", trackHandler.GetPlaylistTracks).Methods("GET")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}", playlistHandler.GetPlaylistByID).Methods("GET")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}/like", playlistHandler.LikePlaylist).Methods("POST")
	r.HandleFunc("/api/v1/playlists/search", playlistHandler.SearchPlaylists).Methods("GET")

	r.HandleFunc("/api/v1/auth/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/logout", userHandler.Logout).Methods("POST")
	r.HandleFunc("/api/v1/auth/check", userHandler.CheckUser).Methods("GET")

	r.HandleFunc("/api/v1/user/me/avatar", userHandler.UploadAvatar).Methods("POST")
	r.HandleFunc("/api/v1/user/me", userHandler.ChangeUserData).Methods("PUT")
	r.HandleFunc("/api/v1/user/me", userHandler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/v1/user/{username:[a-zA-Z0-9_]+}", userHandler.GetUserData).Methods("GET")
	r.HandleFunc("/api/v1/user/me/history", trackHandler.GetLastListenedTracks).Methods("GET")
	r.HandleFunc("/api/v1/user/{username:[a-zA-Z0-9_]+}/artists", artistHandler.GetFavoriteArtists).Methods("GET")
	r.HandleFunc("/api/v1/user/{username:[a-zA-Z0-9_]+}/tracks", trackHandler.GetFavoriteTracks).Methods("GET")
	r.HandleFunc("/api/v1/user/{username:[a-zA-Z0-9_]+}/playlists", playlistHandler.GetProfilePlaylists).Methods("GET")
	r.HandleFunc("/api/v1/user/me/albums", albumHandler.GetFavoriteAlbums).Methods("GET")

	r.HandleFunc("/api/v1/label", labelHandler.CreateLabel).Methods("POST")
	r.HandleFunc("/api/v1/label", labelHandler.UpdateLabel).Methods("PUT")
	r.HandleFunc("/api/v1/label/{id:[0-9]+}", labelHandler.GetLabel).Methods("GET")

	r.HandleFunc("/api/v1/label/artist", labelHandler.CreateArtist).Methods("POST")
	r.HandleFunc("/api/v1/label/artist", labelHandler.EditArtist).Methods("PUT")
	r.HandleFunc("/api/v1/label/artists", labelHandler.GetArtists).Methods("GET")
	r.HandleFunc("/api/v1/label/artist", labelHandler.DeleteArtist).Methods("DELETE")

	r.HandleFunc("/api/v1/label/album", labelHandler.CreateAlbum).Methods("POST")
	r.HandleFunc("/api/v1/label/album", labelHandler.DeleteAlbum).Methods("DELETE")
	r.HandleFunc("/api/v1/label/albums", labelHandler.GetAlbumsByLabelID).Methods("GET")

	r.HandleFunc("/api/v1/jams", jamHandler.CreateRoom).Methods("POST")
	r.HandleFunc("/api/v1/jams/{id}", jamHandler.WSHandler).Methods("GET")

	r.Handle("/api/v1/metrics", promhttp.Handler())

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

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Error shutting down server:", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("Composer server stopped")
}
