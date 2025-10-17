package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/config"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/logger"
	_ "github.com/lib/pq"
)

var (
	cfg config.Config
	db  *sql.DB
)

func main() {
	var err error

	loadFlags()

	cfg.Logger = logger.NewAppLogger(cfg.Env)


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

func loadFlags() {
	flag.StringVar(&cfg.Port, "port", "4000", "API Server Port")
	flag.StringVar(&cfg.Env, "env", "development",
		"Environment (Production | Staging | Production)")

	flag.StringVar(&cfg.DB.Dsn, "db-dsn",
		"postgres://postgres:3011@localhost:5432/bookings?sslmode=disable", "PostgreSQL DSN")

	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25,
		"PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25,
		"PostgreSQL max open idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m",
		"PostgreSQL max connection idle time")


	flag.Parse()
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

	fmt.Println("DB connected successfully")

	return db, nil
}
