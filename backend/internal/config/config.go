package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Storage StorageConfig
	User    UserConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type StorageConfig struct {
	Path          string
	MaxUploadSize int64
}

type UserConfig struct {
	DefaultUserID int
}

var AppConfig *Config

func Init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "nexushub_personal"),
		},
		Storage: StorageConfig{
			Path:          getEnv("STORAGE_PATH", "./storage"),
			MaxUploadSize: 100 * 1024 * 1024, // 100MB
		},
		User: UserConfig{
			DefaultUserID: 1,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
