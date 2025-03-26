package main

import (
	"flag"
	"fmt"
	"net/http"

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
)

func main() {
	port := flag.String("p", ":8080", "server port")
	flag.Parse()

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", *port)

	r.Use(middleware.Logger)
	r.Use(middleware.RequestId)
	r.Use(middleware.AccessLog)

	r.NotFoundHandler = middleware.NotFoundHandler()

	trackHandler := trackHttp.NewTrackHandler(trackUsecase.NewUsecase(trackRepository.NewTrackMemoryRepository()))
	albumHandler := albumHttp.NewAlbumHandler(albumUsecase.NewUsecase(albumRepository.NewAlbumMemoryRepository()))
	artistHandler := artistHttp.NewArtistHandler(artistUsecase.NewUsecase(artistRepository.NewArtistMemoryRepository()))

	r.Handle("/tracks", middleware.Pagination(trackHandler.GetAllTracks)).Methods("GET")
	r.Handle("/albums", middleware.Pagination(albumHandler.GetAllAlbums)).Methods("GET")
	r.Handle("/artists", middleware.Pagination(artistHandler.GetAllArtists)).Methods("GET")

	err := http.ListenAndServe(*port, r)
	if err != nil {
		fmt.Println(err)
	}
}
