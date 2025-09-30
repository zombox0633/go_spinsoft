package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigType struct {
	Port         string
	APIKey       string
	MonGoURL     string
	DataBaseName string
}

func LoadConfig() *ConfigType {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &ConfigType{
		Port:         getEnv("PORT", "8080"),
		MonGoURL:     getEnv("MONGO_URI", ""),
		DataBaseName: getEnv("DB_NAME", ""),
		APIKey:       getEnv("API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
