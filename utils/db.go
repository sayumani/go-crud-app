package utils

import (
	"database/sql"
	"fmt"
	"log"
)

//CreateDbConnection creates a db connection
func CreateDbConnection(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	var err error
	DB, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Failed to create a connection %w", ErrDBConnectionError)
	}
	return DB, nil
}
