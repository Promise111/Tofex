package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadDotEnv reads a .env file into the process environment.
// Existing environment variables are not overwritten.
// Search order: ENV_FILE path, ./.env, ./backend/.env (repo root).
func LoadDotEnv() {
	if p := os.Getenv("ENV_FILE"); p != "" {
		if err := godotenv.Load(p); err != nil {
			log.Printf("env: could not load ENV_FILE=%q: %v", p, err)
		}
		return
	}
	for _, p := range []string{".env", "server/.env", "backend/.env"} {
		if err := godotenv.Load(p); err == nil {
			log.Printf("env: loaded %s", p)
			return
		}
	}
}
