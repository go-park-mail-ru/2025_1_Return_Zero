package main

import (
	"flag"
	"fmt"
	"net/http"

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
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Cors middleware.Cors
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host returnzero.ru
// @BasePath /
func main() {
	port := flag.String("p", ":8080", "server port")
	flag.Parse()

	config, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", *port)

	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	r.Use(middleware.Logger)
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)
	r.Use(config.Cors.Middleware)

	r.NotFoundHandler = middleware.NotFoundHandler()

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackRepository.NewTrackMemoryRepository(), artistRepository.NewArtistMemoryRepository(), albumRepository.NewAlbumMemoryRepository()))
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumRepository.NewAlbumMemoryRepository(), artistRepository.NewArtistMemoryRepository()))
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(artistRepository.NewArtistMemoryRepository()))

	r.Handle("/tracks", middleware.Pagination(trackHandler.GetAllTracks)).Methods("GET")
	r.Handle("/albums", middleware.Pagination(albumHandler.GetAllAlbums)).Methods("GET")
	r.Handle("/artists", middleware.Pagination(artistHandler.GetAllArtists)).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	staticFileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./public")))
	r.PathPrefix("/static/").Handler(staticFileHandler)

	err = http.ListenAndServe(*port, r)
	if err != nil {
		fmt.Println(err)
	}
}
