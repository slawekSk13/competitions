package models

import (
	"competition-app/config"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB(cfg *config.Config) error {
	var err error
	
	log.Printf("Connecting to PostgreSQL at %s:%d...", cfg.DBHost, cfg.DBPort)
	
	DB, err = sql.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connection established")
	return nil
}
