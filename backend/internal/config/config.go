package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	APIPort        string
	BcryptCost     int
	CORSOrigins    []string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	
	if port == "" {
		port = os.Getenv("API_PORT")
	}
	
	if port == "" {
		port = "8080"
	}

	bcryptCost := 12
	if v := os.Getenv("BCRYPT_COST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("BCRYPT_COST must be an integer: %w", err)
		}
		bcryptCost = parsed
	}

	corsOrigins := []string{"http://localhost:3000", "http://localhost:5173"}
	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		corsOrigins = strings.Split(v, ",")
		for i := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
		}
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		APIPort:     port,
		BcryptCost:  bcryptCost,
		CORSOrigins: corsOrigins,
	}, nil
}
