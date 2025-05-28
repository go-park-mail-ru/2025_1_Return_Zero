package postgres

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/db/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectPostgres(cfg config.PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

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

func RunMigrations(cfg config.PostgresConfig) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)

	migrator, err := migrations.NewMigrator(dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	now, exp, info, err := migrator.Info()
	if err != nil {
		return fmt.Errorf("failed to get migration info: %w", err)
	}

	if now < exp {
		fmt.Println("Migration needed, current state:")
		fmt.Println(info)

		err = migrator.Migrate()
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		fmt.Println("Migration successful!")
	} else {
		fmt.Println("No database migration needed")
	}

	return nil
}
