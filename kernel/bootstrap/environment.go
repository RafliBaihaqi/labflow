package bootstrap

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// loadEnvFile loads the given .env file if it exists. A missing file is not an
// error (container deployments inject config through the environment instead).
// A malformed file IS an error so configuration mistakes surface loudly.
func loadEnvFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return godotenv.Load(path)
}

// LoadEnvVariables loads ".env" from the working directory when present.
func LoadEnvVariables() {
	if err := loadEnvFile(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
