package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	CurrencyAPI struct {
		Address string `env:"CURRENCY_API_URL" env-default:"http://nationalbank.kz/rss/get_rates.cfm"`
	}

	Postgres struct {
		Address     string `env:"PG_URL" env-default:"postgresql://testuser:12345@localhost:5432/test"`
		MaxPoolSize int    `env:"PG_MAX_POOL_SIZE" env-default:"10"`
		// See https://pkg.go.dev/github.com/jackc/pgx#pkg-variables.
		LogLevel int `env:"PG_LOG_LEVEL" env-default:"5"`
	}

	Logger struct {
		// See https://pkg.go.dev/golang.org/x/exp/slog#Level.
		Level int `env:"LOGGER_LEVEL" env-default:"-4"`
	}

	HTTPServer struct {
		Addr            string        `env:"HTTP_ADDR" env-default:"0.0.0.0:8000"`
		RateLimit       int           `env:"HTTP_RATE_LIMIT" env-default:"100"`
		RateWindow      time.Duration `env:"HTTP_RATE_WINDOW" env-default:"30s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUT_DOWN_TIMEOUT" env-default:"10s"`
	}

	Config struct {
		PG          Postgres
		Logger      Logger
		HTTPServer  HTTPServer
		CurrencyAPI CurrencyAPI
	}
)

func ParseConfig() (Config, error) {
	cfg := new(Config)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return *cfg, fmt.Errorf("config - ParseConfig: %w", err)
	}
	return *cfg, nil
}
