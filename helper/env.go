package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)

}