package common

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port               string
	Environment        string
	DatabaseURL        string
	AccesssTokenTTL    int
	RefreshTokenTTL    int
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func Load(envPath string) (*Config, error) {
	port := getEnv("PORT")
	env := getEnv("ENVIRONMENT")
	dbUrl := getEnv("DATABASE_URL")
	accessTokenSecret := getEnv("ACCESS_TOKEN_SECRET")
	refreshTokenSecret := getEnv("REFRESH_TOKEN_SECRET")

	// sessions time to live in seconds
	accesssTokenTTL, err := strconv.Atoi(getEnv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, errors.New("ACCESS_TOKEN_TTL have to be an int")
	}
	refreshTokenTTL, err := strconv.Atoi(getEnv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, errors.New("REFRESH_TOKEN_TTL have to be an int")
	}

	// extra validation
	if dbUrl == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if accessTokenSecret == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if refreshTokenSecret == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	tempPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("port needs to be ")
	}

	if tempPort <= 0 || tempPort > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got %s", port)
	}

	return &Config{
		Port:               port,
		Environment:        env,
		DatabaseURL:        dbUrl,
		AccesssTokenTTL:    accesssTokenTTL,
		RefreshTokenTTL:    refreshTokenTTL,
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
	}, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE)URL is required")
	}

	tempPort, err := strconv.Atoi(c.Port)
	if err != nil {
		return errors.New("port needs to be ")
	}

	if tempPort <= 0 || tempPort > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %s", c.Port)
	}

	return nil
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}
