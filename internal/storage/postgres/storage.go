package postgres

import (
	"database/sql"
	"fmt"
	"os"
)

func NewDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USERNAME")
	dbHost := os.Getenv("DB_HOST")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbPort == "" {
		return nil, fmt.Errorf("DB_USERNAME, DB_PASSWORD, DB_NAME или DB_PORT не установлены")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("[NewDB]: Ошибка при открытии подключения к БД: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[NewDB]: Ошибка подключения к бд: %v", err)
	}

	return db, nil
}
