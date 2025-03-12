package main

import (
	"flag"
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
// @host returnzero.ru
// TODO: change host to the production host
// @BasePath /

func main() {
	port := flag.String("p", ":8080", "server port")
	flag.Parse()

	cors := &Cors{
		AllowedOrigins:   []string{"http://returnzero.ru", "http://127.0.0.1:3000", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

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

	fmt.Printf("Server starting on port %s...\n", *port)

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/tracks", tracksHandler.List).Methods("GET")
	r.HandleFunc("/albums", albumsHandler.List).Methods("GET")
	r.HandleFunc("/artists", artistsHandler.List).Methods("GET")

	userApi := NewMyHandler()
	r.HandleFunc("/login", userApi.loginHandler).Methods("POST")
	r.HandleFunc("/logout", userApi.logoutHandler).Methods("POST")
	r.HandleFunc("/signup", userApi.signupHandler).Methods("POST")
	r.HandleFunc("/user", userApi.checkSession).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	r.HandleFunc("/static/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/"+mux.Vars(r)["path"])
	})

	err := http.ListenAndServe(*port, cors.Middleware(r))
	if err != nil {
		fmt.Println(err)
	}
}
