package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server settings
	ServerPort string

	// Database settings
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Redis settings
	RedisHost string
	RedisPort int
}

// LoadConfig loads the configuration from environment variables and secrets
func LoadConfig() (*Config, error) {
	// Set defaults
	cfg := &Config{
		ServerPort: "8080",
		DBHost:     "postgres",
		DBPort:     5432,
		RedisHost:  "redis",
		RedisPort:  6379,
	}

	// Server settings
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.ServerPort = port
	}

	// Database settings
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.DBHost = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.DBPort = p
		}
	}

	// Read secrets
	dbUser, err := readFileOrEnv("/run/secrets/db_user", "DB_USER", "postgres")
	if err != nil {
		return nil, err
	}
	cfg.DBUser = dbUser

	dbPassword, err := readFileOrEnv("/run/secrets/db_password", "DB_PASSWORD", "")
	if err != nil {
		return nil, err
	}
	cfg.DBPassword = dbPassword

	dbName, err := readFileOrEnv("/run/secrets/db_name", "DB_NAME", "postgres")
	if err != nil {
		return nil, err
	}
	cfg.DBName = dbName

	// Redis settings
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.RedisHost = host
	}

	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.RedisPort = p
		}
	}

	return cfg, nil
}

// GetDBConnString returns the database connection string
func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

// GetRedisConnString returns the Redis connection string
func (c *Config) GetRedisConnString() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}

// readFileOrEnv reads from a file, then environment variable, then fallback
func readFileOrEnv(filePath, envVar, fallback string) (string, error) {
	// Try reading from file first
	content, err := os.ReadFile(filePath)
	if err == nil {
		return string(content), nil
	}

	// If file doesn't exist or can't be read, try environment variable
	if value := os.Getenv(envVar); value != "" {
		return value, nil
	}

	// If no value found and fallback is empty, return error
	if fallback == "" {
		return "", fmt.Errorf("no value for %s found in file %s or environment", envVar, filePath)
	}

	// Return fallback value
	return fallback, nil
}
