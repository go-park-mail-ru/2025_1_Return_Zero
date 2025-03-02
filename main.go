package main

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
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

	fmt.Println("Server starting on port 8081...")
	r.HandleFunc("/docs/", httpSwagger.WrapHandler)
	r.HandleFunc("/tracks", tracksHandler.List).Methods("GET")
	r.HandleFunc("/albums", albumsHandler.List).Methods("GET")
	r.HandleFunc("/artists", artistsHandler.List).Methods("GET")
  
	api := NewMyHandler()
	r.HandleFunc("/login", api.loginHandler).Methods("POST")
	r.HandleFunc("/logout", api.logoutHandler).Methods("POST")
	r.HandleFunc("/signup", api.signupHandler).Methods("POST")
	r.HandleFunc("/user", api.checkSession).Methods("GET")
	
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Println(err)
	}
}
