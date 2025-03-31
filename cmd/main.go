package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	authRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth/repository"
	userHttp "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/repository"
	userUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/redis"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
)

func main() {
	port := flag.String("p", ":8084", "server port")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	redisConn, err := redis.ConnectRedis(cfg.Redis)
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	defer redisConn.Close()

	db, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	r := mux.NewRouter()
	fmt.Printf("Server starting on port %s...\n", *port)

	userHandler := userHttp.NewUserHandler(userUsecase.NewUserUsecase(userRepository.NewUserPostgresRepository(db), authRepository.NewAuthRedisRepository(redisConn)))

	r.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/logout", userHandler.Logout).Methods("POST")
	r.HandleFunc("/user", userHandler.CheckUser).Methods("GET")
	err = http.ListenAndServe(*port, r)
	if err != nil {
		fmt.Println(err)
	}
}
