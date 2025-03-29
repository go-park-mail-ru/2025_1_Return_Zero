package postgres

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectPostgres(cfg config.PostgresConfig) (*sql.DB, error) {
	dsn := cfg.DSN

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Postgres")

	return db, nil
}
