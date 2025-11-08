package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Filename  string
	Address   string
	Route     string
	Port      string
	Threshold int
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	cfg := &Config{
		Filename:  getEnv("FILENAME", "timestamps.log"),
		Address:   getEnv("ADDRESS", "localhost"),
		Route:     getEnv("ROUTE", "/"),
		Port:      getEnv("PORT", "8000"),
		Threshold: 60,
	}

	thresholdStr := getEnv("THRESHOLD", "60")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		return nil, fmt.Errorf("invalid threshold value: %w", err)
	}
	cfg.Threshold = threshold

	return cfg, nil
}

// ServerAddr returns the full server address
func (c *Config) ServerAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
