package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	DefaultOffset = 0
	DefaultLimit  = 10
)

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host localhost:8080
// TODO: change host to the production host
// @BasePath /
func main() {
	r := mux.NewRouter()

	tracksHandler := &TracksHandler{
		Model: models.NewTracksModel(),
	}

	albumsHandler := &AlbumsHandler{
		Model: models.NewAlbumsModel(),
	}

	artistsHandler := &ArtistsHandler{
		Model: models.NewArtistsModel(),
	}

	fmt.Println("Server starting on port 8080...")

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/tracks", tracksHandler.List).Methods("GET")
	r.HandleFunc("/albums", albumsHandler.List).Methods("GET")
	r.HandleFunc("/artists", artistsHandler.List).Methods("GET")

	userApi := NewMyHandler()
	r.HandleFunc("/login", userApi.loginHandler).Methods("POST")
	r.HandleFunc("/logout", userApi.logoutHandler).Methods("POST")
	r.HandleFunc("/signup", userApi.signupHandler).Methods("POST")
	r.HandleFunc("/user", userApi.checkSession).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
