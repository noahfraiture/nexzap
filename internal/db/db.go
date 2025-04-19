package db

import (
	"context"
	"fmt"
	"log"
	db "nexzap/internal/db/generated"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var dbPool *pgxpool.Pool
var repo *db.Queries

func GetPool() *pgxpool.Pool {
	if dbPool == nil {
		log.Fatal(fmt.Errorf("DB not initialized"))
	}
	return dbPool
}

func GetRepository() *db.Queries {
	if repo == nil {
		log.Fatal(fmt.Errorf("DB not initialized"))
	}
	return repo
}

func Init() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning: Error loading .env file: %v", err)
	}

	password, err := getPassword()
	if err != nil {
		return err
	}
	connStrPgx := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("POSTGRES_USER"),
		password,
		os.Getenv("POSTGRES_HOST"),
		"5432",
		os.Getenv("POSTGRES_DB"),
	)

	dbPool, err = pgxpool.New(context.Background(), connStrPgx)
	if err != nil {
		return err
	}
	repo = db.New(dbPool)

	return nil
}

func Populate() error {
	password, err := getPassword()
	if err != nil {
		return err
	}
	if err := buildDatabase(password); err != nil && err != migrate.ErrNoChange {
		fmt.Println("error during building of database")
		return err
	}
	return nil

}

func buildDatabase(password string) error {
	// Build the database from the migrations files
	connStrMigration := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&search_path=public",
		os.Getenv("POSTGRES_USER"),
		password,
		os.Getenv("POSTGRES_HOST"),
		"5432",
		os.Getenv("POSTGRES_DB"),
	)
	m, err := migrate.New(
		"file://internal/db/migrations",
		connStrMigration,
	)
	if err != nil {
		return err
	}
	return m.Up()
}

func getPassword() (string, error) {
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		passwordFile, ok := os.LookupEnv("POSTGRES_PASSWORD_FILE")
		if !ok {
			return "", fmt.Errorf("no password set")
		}
		data, err := os.ReadFile(passwordFile)
		if err != nil {
			return "", err
		}
		password = strings.TrimSpace(string(data))
	}
	return password, nil
}

// NukeDatabase drops the public schema and all its contents from the database.
// WARNING: This is a destructive operation and will delete all data in the public schema.
func NukeDatabase() error {
	if dbPool == nil {
		return fmt.Errorf("DB not initialized")
	}

	_, err := dbPool.Exec(context.Background(), "DROP SCHEMA IF EXISTS public CASCADE")
	if err != nil {
		return fmt.Errorf("failed to drop public schema: %v", err)
	}

	_, err = dbPool.Exec(context.Background(), "CREATE SCHEMA public")
	if err != nil {
		return fmt.Errorf("failed to recreate public schema: %v", err)
	}

	return nil
}
