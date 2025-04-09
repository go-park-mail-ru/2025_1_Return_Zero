package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
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
	genreRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/genre/repository"
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
)

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host returnzero.ru
// @BasePath /api/v1
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	redisConn, err := redis.ConnectRedis(cfg.Redis)
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	defer redisConn.Close()

	postgresConn, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer postgresConn.Close()

	s3, err := s3.InitS3(cfg.S3)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", cfg.Port)

	r.PathPrefix("/api/v1/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/v1/docs/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	newUserUsecase := userUsecase.NewUserUsecase(userRepository.NewUserPostgresRepository(postgresConn), authRepository.NewAuthRedisRepository(redisConn), userFileRepo.NewS3Repository(s3, cfg.S3.S3_IMAGES_BUCKET))

	r.Use(middleware.Logger)
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Auth(newUserUsecase))
	r.Use(cfg.Cors.Middleware)

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackRepository.NewTrackPostgresRepository(postgresConn), artistRepository.NewArtistPostgresRepository(postgresConn), albumRepository.NewAlbumPostgresRepository(postgresConn), trackFileRepo.NewS3Repository(s3, cfg.S3.S3_TRACKS_BUCKET, cfg.S3.S3_DURATION)), cfg)
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumRepository.NewAlbumPostgresRepository(postgresConn), artistRepository.NewArtistPostgresRepository(postgresConn), genreRepository.NewGenrePostgresRepository(postgresConn)), cfg)
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(artistRepository.NewArtistPostgresRepository(postgresConn)), cfg)
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
	r.HandleFunc("/api/v1/user/{username}/avatar", userHandler.GetUserAvatar).Methods("GET")
	r.HandleFunc("/api/v1/user/{username}/avatar", userHandler.UploadAvatar).Methods("POST")

	err = http.ListenAndServe(cfg.Port, r)
	if err != nil {
		fmt.Println(err)
	}
}
