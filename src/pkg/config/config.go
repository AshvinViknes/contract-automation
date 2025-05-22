package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// Server
	ServerPort   string
	ServerHost   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret        string
	JWTExpiration    time.Duration
	JWTRefreshToken  bool
	RefreshSecretKey string

	// Environment
	Env            string
	LogLevel       string
	AllowedOrigins []string

	// Setup
	SetupToken string
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, separator)
	}
	return defaultValue
}

func New() *Config {
	return &Config{
		// Server defaults
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		ServerHost:   getEnv("SERVER_HOST", "0.0.0.0"),
		ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),

		// Database defaults
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "mysecretpassword"),
		DBName:     getEnv("DB_NAME", "smart_contracts"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT defaults
		JWTSecret:        getEnv("JWT_SECRET", "your-default-secret-key-change-in-production"),
		JWTExpiration:    getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		JWTRefreshToken:  getEnvAsBool("JWT_REFRESH_TOKEN", true),
		RefreshSecretKey: getEnv("REFRESH_SECRET_KEY", "your-refresh-secret-key-change-in-production"),

		// Environment defaults
		Env:            getEnv("APP_ENV", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}, ","),

		// Setup defaults
		SetupToken: getEnv("SETUP_TOKEN", "your-secure-setup-token-change-in-production"),
	}
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
}
