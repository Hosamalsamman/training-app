package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Pool *pgxpool.Pool
var DB *gorm.DB

func Connect() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// pgx
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		url.QueryEscape(password),
		host,
		port,
		dbname,
	)

	var err error

	Pool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("pgx connection failed:", err)
	}

	// GORM
	gormDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
	)

	DB, err = gorm.Open(postgres.Open(gormDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("gorm connection failed:", err)
	}

	log.Printf(
		"Connected to PostgreSQL (%s/%s)",
		host,
		dbname,
	)
}
