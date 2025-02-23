package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-park-mail-ru/2025_1_Return_Zero/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Return Zero API
// @version 1.0
// @description This is the API server for Return Zero music app.
// @host localhost:8080
// TODO: change host to the production host
// @BasePath /
func main() {
	fmt.Println("Server starting on port 8080...")
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)
	http.HandleFunc("/tracks", getTracksHandler)
	http.HandleFunc("/albums", getAlbumsHandler)
	http.HandleFunc("/artists", getArtistsHandler)
	http.ListenAndServe(":8080", nil)
}
