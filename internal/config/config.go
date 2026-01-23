package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	HTTP  HTTPConfig
}

type AppConfig struct {
	Env string
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type HTTPConfig struct {
	Port          string
	BasicAuthUser string
	BasicAuthPass string
}

func Load() *Config {
	cfg := &Config{
		App: AppConfig{
			Env: getEnv("APP_ENV", "local"),
		},
		DB: DBConfig{
			DSN: getEnv(
				"DB_DSN",
				"postgres://postgres:postgres@localhost:5432/productdb?sslmode=disable",
			),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Enabled:  getEnvAsBool("REDIS_ENABLED", true),
		},
		HTTP: HTTPConfig{
			Port:          getEnv("HTTP_PORT", "8080"),
			BasicAuthUser: getEnv("BASIC_AUTH_USER", "admin"),
			BasicAuthPass: getEnv("BASIC_AUTH_PASS", "admin123"),
		},
	}

	log.Printf("config loaded (env=%s)", cfg.App.Env)
	return cfg
}

func getEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return defaultVal
}
