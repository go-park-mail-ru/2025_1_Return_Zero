package main

import (
	"flag"
	"fmt"
	"net/http"

	authRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth/repository"
	userHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/repository"
	userUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("p", ":8080", "server port")
	flag.Parse()

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", *port)

	userHandler := userHttp.NewUserHandler(userUsecase.NewUserUsecase(userRepository.NewUserMemoryRepository(), authRepository.NewAuthMemoryRepository()))

	r.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/logout", userHandler.Logout).Methods("POST")
	r.HandleFunc("/user", userHandler.CheckUser).Methods("GET")
	err := http.ListenAndServe(*port, r)
	if err != nil {
		fmt.Println(err)
	}
}
