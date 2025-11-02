package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/config"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/logger"
	_ "github.com/lib/pq"
)

var (
	cfg config.Config
	db  *sql.DB
)

func main() {
	cfg.Logger = logger.NewAppLogger(cfg.Env)

	if os.Getenv("DOCKER_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			cfg.Logger.Fatal("Failed to load environmental variables", "error", err)
		}
	}

	cfg.Port = os.Getenv("PORT")
	cfg.Env = os.Getenv("ENV")
	cfg.JwtSecret = os.Getenv("JWT_SECRET")
	cfg.DB.Dsn = os.Getenv("DSN")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if cfg.DB.Dsn == "" {
		cfg.DB.Dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	if err != nil {
		cfg.Logger.Error("Error converting string to int", "Error", err)
	}
	cfg.DB.MaxIdleConns = maxIdleConns
	cfg.DB.MaxIdleTime = os.Getenv("MAX_IDLE_TIME")
	maxOpenConns, err := strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))
	if err != nil {
		cfg.Logger.Error("Error converting string to int", "Error", err)
	}
	cfg.DB.MaxOpenConns = maxOpenConns

	db, err = ConnectDB()
	if err != nil {
		cfg.Logger.Fatal("Failed to connect to database", "error", err)
	}

	defer db.Close()

	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: routes(),
	}

	cfg.Logger.Info("Application starting...",
		"env", cfg.Env,
		"port", cfg.Port,
		"version", "1.0.0",
	)

	if err := srv.ListenAndServe(); err != nil {
		cfg.Logger.Fatal("Server failed to start", "error", err)
	}
}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
