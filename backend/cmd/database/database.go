package database

import (
	"fmt"
	"log"

	"github.com/codepnw/ticket-api/cmd/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDatabase(cfg *config.EnvConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(`
		host=%s user=%s dbname=%s password=%s sslmode=%s port=%d`,
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBName,
		cfg.DBPassword,
		cfg.DBSSLMode,
		cfg.DBPort,
	)

	db, err := sqlx.Connect(cfg.DBDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %e", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping check database: %e", err)
	}

	log.Println("connected to database!")
	return db, nil
}
