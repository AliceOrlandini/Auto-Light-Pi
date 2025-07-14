package config

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	DbDsn := os.Getenv("DB_DSN")

	psqlDB, err := sql.Open("postgres", DbDsn)
	if err != nil {
		return err
	}

	err = psqlDB.Ping()
	if err != nil {
		return err
	}

	DB = psqlDB
	return nil
}