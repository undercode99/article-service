package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort    string
	AappMode   string
	ElasticUrl string
	Cache      *RedisConfig
	Database   *DatabaseConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type DatabaseConfig struct {
	Dsn string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Dsn: getEnvString("DATABASE_DSN", ""),
	}
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     getEnvString("REDIS_HOST", "localhost"),
		Port:     getEnvString("REDIS_PORT", "6379"),
		Password: getEnvString("REDIS_PASSWORD", ""),
		DB:       getEnvInt("REDIS_DB", 0),
	}
}

func NewConfig() *Config {
	return &Config{
		AppPort:    getEnvString("APP_PORT", "8080"),
		AappMode:   getEnvString("APP_MODE", "development"),
		ElasticUrl: getEnvString("ELASTIC_URL", "http://localhost:9200"),
		Cache:      NewRedisConfig(),
		Database:   NewDatabaseConfig(),
	}
}

func (c *Config) AppModeIsProduction() bool {
	return c.AappMode == "production"
}

func getEnvString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
