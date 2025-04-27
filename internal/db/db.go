package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	generated "nexzap/internal/db/generated"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DEFAULT_USER     = "nexzap"
	DEFAULT_HOST     = "localhost"
	DEFAULT_PORT     = "5432"
	DEFAULT_DATABASE = "nexzap"
	DEFAULT_PASSWORD = "nexzap"
)

// Database struct to encapsulate the connection pool and repository
type Database struct {
	pool *pgxpool.Pool
	repo *generated.Queries
}

// NewDatabase initializes the Database struct with connection pooling and retries
func NewDatabase() (*Database, error) {
	var err error
	creds := getCredentials()
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s host=%s",
		creds.user,
		creds.password,
		creds.host,
		creds.port,
		creds.database,
		creds.host,
	)

	// Simple retry mechanism with exponential backoff
	var pool *pgxpool.Pool
	for attempt := range 5 {
		pool, err = pgxpool.New(context.Background(), connStr)
		if err == nil {
			break
		}
		wait := time.Duration(1<<attempt) * time.Second
		log.Printf("Failed to connect to database: %v. Retrying in %v...", err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %v", err)
	}

	poolConfig, _ := pgxpool.ParseConfig(connStr)
	poolConfig.MaxConns = 20
	pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	repository := generated.New(pool)
	return &Database{pool: pool, repo: repository}, nil
}

func (d *Database) GetRepository() *generated.Queries {
	return d.repo
}

func (d *Database) Close() {
	d.pool.Close()
}

func (d *Database) Populate() error {
	creds := getCredentials()
	connStrMigration := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&search_path=public",
		creds.user,
		creds.password,
		creds.host,
		creds.port,
		creds.database,
	)
	m, err := migrate.New(os.Getenv("MIGRATIONS_PATH"), connStrMigration)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}
	return nil
}

func (d *Database) NukeDatabase() error {
	if os.Getenv("ENV") != "dev" {
		return fmt.Errorf("NukeDatabase can only be run in development environment")
	}
	_, err := d.pool.Exec(context.Background(), "DROP SCHEMA IF EXISTS public CASCADE")
	if err != nil {
		return fmt.Errorf("failed to drop public schema: %v", err)
	}
	_, err = d.pool.Exec(context.Background(), "CREATE SCHEMA public")
	if err != nil {
		return fmt.Errorf("failed to recreate public schema: %v", err)
	}
	return nil
}

func (d *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.pool.Ping(ctx)
}

func getCredentials() credentials {
	creds := credentials{
		user:     DEFAULT_USER,
		host:     DEFAULT_HOST,
		port:     DEFAULT_PORT,
		database: DEFAULT_DATABASE,
		password: DEFAULT_PASSWORD,
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		creds.user = user
	}
	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		creds.host = host
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		creds.port = port
	}
	if db := os.Getenv("POSTGRES_DB"); db != "" {
		creds.database = db
	}
	if password, err := getPassword(); err == nil {
		creds.password = password
	}
	return creds
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

type credentials struct {
	user     string
	host     string
	port     string
	database string
	password string
}
