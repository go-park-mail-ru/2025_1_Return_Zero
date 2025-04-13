package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/init/postgres"
)

func main() {
	sqlFile := flag.String("file", "data.sql", "path to sql file")
	flag.Parse()
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg.Postgres.POSTGRES_HOST = "localhost"
	db, err := postgres.ConnectPostgres(cfg.Postgres)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	sql, err := os.ReadFile(*sqlFile)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Println(err)
	}

}
