package utils

import (
	"os"

	// "github.com/joho/godotenv"
)

func EnvVariable(key string) string {
	// err := godotenv.Load(".env")
	// CheckErr(err)
	return os.Getenv(key)
}
