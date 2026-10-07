package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Filename  string
	Route     string
	Port      string
	Threshold int
}

func Load() (*Config, error) {
	threshold, err := strconv.Atoi(getEnv("THRESHOLD", "60"))
	if err != nil || threshold <= 0 {
		return nil, fmt.Errorf("invalid threshold value %q: must be a positive integer", os.Getenv("THRESHOLD"))
	}
	return &Config{
		Filename:  getEnv("FILENAME", "timestamps.log"),
		Route:     getEnv("ROUTE", "/"),
		Port:      getEnv("PORT", "8000"),
		Threshold: threshold,
	}, nil
}

func (c *Config) ServerAddr() string {
	return ":" + c.Port
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
