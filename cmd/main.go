package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	albumHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/delivery/http"
	albumRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/repository"
	albumUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/usecase"
	artistHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/delivery/http"
	artistRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/repository"
	artistUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/usecase"
	trackHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/delivery/http"
	trackRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/repository"
	trackUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/usecase"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host returnzero.ru
// @BasePath /
func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", config.Port)

	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	r.Use(middleware.Logger)
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)
	r.Use(config.Cors.Middleware)

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackRepository.NewTrackMemoryRepository(), artistRepository.NewArtistMemoryRepository(), albumRepository.NewAlbumMemoryRepository()), config)
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumRepository.NewAlbumMemoryRepository(), artistRepository.NewArtistMemoryRepository()), config)
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(artistRepository.NewArtistMemoryRepository()), config)

	r.HandleFunc("/tracks", trackHandler.GetAllTracks).Methods("GET")
	r.HandleFunc("/albums", albumHandler.GetAllAlbums).Methods("GET")
	r.HandleFunc("/artists", artistHandler.GetAllArtists).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	staticFileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./public")))
	r.PathPrefix("/static/").Handler(staticFileHandler)

	err = http.ListenAndServe(config.Port, r)
	if err != nil {
		fmt.Println(err)
	}
}
