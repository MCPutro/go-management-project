package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/MCPutro/go-management-project/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(config *config.DatabaseConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = sql.Open(config.PostgresSql.Name, config.PostgresSql.DSN)
		if err != nil {
			return nil, err
		}

		// test connection
		err = db.Ping()
		if err != nil {
			log.Println("failed to connect to database, retry in 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			db.SetMaxIdleConns(5)
			db.SetMaxOpenConns(100)
			db.SetConnMaxLifetime(60 * time.Minute)
			db.SetConnMaxIdleTime(10 * time.Minute)

			log.Println("successfully connected to database")

			return db, nil
		}
	}

	return nil, errors.New("failed to connect to database, please check your database configuration")

}
