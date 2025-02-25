package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	DefaultOffset = 0
	DefaultLimit  = 10
)

type application struct {
	models *models.Models
}

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host localhost:8080
// TODO: change host to the production host
// @BasePath /
func main() {
	app := &application{
		models: models.NewModels(),
	}

	fmt.Println("Server starting on port 8080...")
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)
	http.HandleFunc("/tracks", app.getTracks)
	http.HandleFunc("/albums", app.getAlbums)
	http.HandleFunc("/artists", app.getArtists)
	http.ListenAndServe(":8080", nil)
}
