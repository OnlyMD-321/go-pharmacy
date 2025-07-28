package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	FirebaseCredentialsPath string
	PostgresDSN           string
}

var AppConfig Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found. Using system environment variables.")
	}

	AppConfig = Config{
		Port:                  getEnv("APP_PORT", "8080"),
		FirebaseCredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", "firebaseServiceAccountKey.json"),
		PostgresDSN:           getEnv("POSTGRES_DSN", ""),
	}

	log.Println("✅ Configuration loaded")
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
