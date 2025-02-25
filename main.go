package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/login", api.loginHandler).Methods("POST")
	r.HandleFunc("/logout", api.logoutHandler).Methods("POST")
	r.HandleFunc("/signup", api.signupHandler).Methods("POST")
	r.HandleFunc("/user", api.checkSession).Methods("GET")
	http.ListenAndServe(":8080", r)
}
