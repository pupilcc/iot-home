package util

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var loadEnvOnce sync.Once

func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	loadEnvOnce.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
			return
		}
		log.Printf("Loaded .env file successfully")
	})

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("Environment variable [%s] not found", key)
	return ""
}

func GetIntEnv(key string) int {
	val, _ := strconv.Atoi(GetEnv(key))
	return val
}
