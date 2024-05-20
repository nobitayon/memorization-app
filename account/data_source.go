package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dataSource struct {
	DB *sqlx.DB
}

func initDS() (*dataSource, error) {
	log.Printf("initializing data sources\n")

	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost,
		pgPort,
		pgUser,
		pgPassword,
		pgDB,
		pgSSL,
	)
	log.Printf("connecting to postgresql\n")
	db, err := sqlx.Open("postgres", pgConnString)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %n", err)
	}

	return &dataSource{
		DB: db,
	}, nil
}

func (d *dataSource) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing postgresql: %w", err)
	}
	return nil
}
