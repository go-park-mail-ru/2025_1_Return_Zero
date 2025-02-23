package main 

import (
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/login", api.loginHandler)
	r.HandleFunc("/logout", api.logoutHandler)
	r.HandleFunc("/signup", api.signupHandler)

	http.ListenAndServe(":8080", r)
}