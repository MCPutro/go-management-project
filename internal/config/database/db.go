package database

import (
	"database/sql"
	"fmt"
	"github.com/MCPutro/go-management-project/internal/config"
	_ "github.com/lib/pq"
	"time"
)

func NewPostgresDB(config *config.DatabaseConfig) (*sql.DB, error) {

	db, err := sql.Open(config.PostgresSql.Name, config.PostgresSql.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	fmt.Println("Database connected successfully")

	return db, nil

}
