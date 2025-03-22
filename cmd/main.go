package main

import (
	"flag"
	"fmt"
	"net/http"
	auth "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/auth"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("p", ":8080", "server port")
	flag.Parse()

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", *port)

	authApi := auth.NewAuthHandler()
	r.HandleFunc("/signup", authApi.Signup).Methods("POST")
	r.HandleFunc("/login", authApi.Login).Methods("POST")
	r.HandleFunc("/logout", authApi.Logout).Methods("POST")
	r.HandleFunc("/check_user", authApi.CheckUser).Methods("GET")
	
	err := http.ListenAndServe(*port, r)
	if err != nil {
		fmt.Println(err)
	}
}
