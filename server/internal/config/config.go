package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	MySQLDSN             string
	JWTSecret            string
	JWTExpiry            time.Duration
	UploadDir            string
	MaxUploadMB          int
	PasswordResetExpiry  time.Duration
	LogPasswordResetLink bool
}

func Load() (*Config, error) {
	jwtExp := getDuration("JWT_EXPIRY", 24*time.Hour)
	resetExp := getDuration("PASSWORD_RESET_EXPIRY", time.Hour)
	maxMB := getInt("MAX_UPLOAD_MB", 10)

	c := &Config{
		Port:                 getenv("PORT", "8080"),
		MySQLDSN:             os.Getenv("MYSQL_DSN"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		JWTExpiry:            jwtExp,
		UploadDir:            getenv("UPLOAD_DIR", "./uploads"),
		MaxUploadMB:          maxMB,
		PasswordResetExpiry:  resetExp,
		LogPasswordResetLink: getenv("LOG_PASSWORD_RESET_LINK", "") == "true",
	}
	if c.MySQLDSN == "" {
		return nil, fmt.Errorf("MYSQL_DSN is required (e.g. user:pass@tcp(127.0.0.1:3306)/tofex?parseTime=true&loc=UTC)")
	}
	if len(c.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	return c, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
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
