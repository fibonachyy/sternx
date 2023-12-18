package service

import (
	"fmt"
	"time"
)

type Config struct {
	JWTDuration       time.Duration
	TokenSymmetricKey string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		JWTDuration: 15 * time.Minute, // Default JWT duration of 15 minutes
	}
}
func validateConfig(config Config) error {
	if config.TokenSymmetricKey == "" {
		return fmt.Errorf("provide a TokenSymmetricKey in the config file")
	}
	return nil
}
