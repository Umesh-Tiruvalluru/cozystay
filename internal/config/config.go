package config

import (
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/logger"
)

type Config struct {
	Port string
	Env  string
	DB   struct {
		Dsn          string // Database String
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Logger    *logger.AppLogger
	JwtSecret string
}
