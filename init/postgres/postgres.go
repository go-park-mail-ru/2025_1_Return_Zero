package postgres

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectPostgres(cfg config.PostgresConfig) (*sql.DB, error) {
	cfg.POSTGRES_HOST = "localhost"
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.POSTGRES_HOST, cfg.POSTGRES_PORT, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB)

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