package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/s3"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	albumHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/delivery/http"
	albumRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/repository"
	albumUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/usecase"
	artistHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/delivery/http"
	artistRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/repository"
	artistUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/usecase"
	genreRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/genre/repository"
	trackHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/delivery/http"
	trackRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/repository"
	trackUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/usecase"
	trackFileRepo "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/trackFile/repository"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host returnzero.ru
// @BasePath /
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	s3, err := s3.InitS3(cfg.S3)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", cfg.Port)

	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	r.Use(middleware.Logger)
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)
	r.Use(cfg.Cors.Middleware)

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackRepository.NewTrackPostgresRepository(db), artistRepository.NewArtistPostgresRepository(db), albumRepository.NewAlbumPostgresRepository(db), trackFileRepo.NewS3Repository(s3, cfg.S3.S3_TRACKS_BUCKET, cfg.S3.S3_DURATION)), cfg)
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumRepository.NewAlbumPostgresRepository(db), artistRepository.NewArtistPostgresRepository(db), genreRepository.NewGenrePostgresRepository(db)), cfg)
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(artistRepository.NewArtistPostgresRepository(db)), cfg)

	r.HandleFunc("/tracks", trackHandler.GetAllTracks).Methods("GET")
	r.HandleFunc("/tracks/{id}", trackHandler.GetTrackByID).Methods("GET")
	r.HandleFunc("/albums", albumHandler.GetAllAlbums).Methods("GET")
	r.HandleFunc("/artists", artistHandler.GetAllArtists).Methods("GET")
	r.HandleFunc("/artists/{id}", artistHandler.GetArtistByID).Methods("GET")
	r.HandleFunc("/artists/{id}/tracks", trackHandler.GetTracksByArtistID).Methods("GET")
	r.HandleFunc("/artists/{id}/albums", albumHandler.GetAlbumsByArtistID).Methods("GET")

	err = http.ListenAndServe(cfg.Port, r)
	if err != nil {
		fmt.Println(err)
	}
}
