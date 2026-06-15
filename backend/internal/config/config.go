package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env                string
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
	AllowedOrigins     []string
	FootballDataAPIKey string

	PointsExact   int
	PointsGD      int
	PointsOutcome int
}

func Load() (*Config, error) {
	cfg := &Config{
		Env:                getenv("APP_ENV", "development"),
		Port:               getenv("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTAccessTTL:       getDuration("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL:      getDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
		AllowedOrigins:     strings.Split(getenv("ALLOWED_ORIGINS", "http://localhost:5173"), ","),
		FootballDataAPIKey: os.Getenv("FOOTBALL_DATA_API_KEY"),
		PointsExact:        getInt("POINTS_EXACT", 5),
		PointsGD:           getInt("POINTS_GD", 3),
		PointsOutcome:      getInt("POINTS_OUTCOME", 1),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" || len(cfg.JWTSecret) < 32 {
		return nil, errors.New("JWT_SECRET must be set and at least 32 chars")
	}
	return cfg, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getDuration(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
