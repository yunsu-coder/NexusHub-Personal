package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Storage StorageConfig
	User    UserConfig
	JWT     JWTConfig
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

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

var AppConfig *Config

func Init() error {
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
			MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 100*1024*1024), // 100MB
		},
		User: UserConfig{
			DefaultUserID: 1,
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-change-in-production"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 168), // 7 days
		},
	}

	// Validate configuration
	if err := AppConfig.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	return nil
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	// Validate database config
	if c.DB.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	if c.DB.Port == "" {
		return fmt.Errorf("database port cannot be empty")
	}
	if c.DB.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}
	if c.DB.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	// Validate storage config
	if c.Storage.Path == "" {
		return fmt.Errorf("storage path cannot be empty")
	}
	if c.Storage.MaxUploadSize <= 0 {
		return fmt.Errorf("max upload size must be positive")
	}
	if c.Storage.MaxUploadSize > 1024*1024*1024 { // 1GB
		log.Printf("Warning: max upload size is very large: %d bytes", c.Storage.MaxUploadSize)
	}

	// Validate JWT config
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}
	if c.JWT.Secret == "default-secret-change-in-production" {
		log.Println("WARNING: Using default JWT secret! Please change it in production!")
	}
	if len(c.JWT.Secret) < 32 {
		log.Printf("WARNING: JWT secret is too short (%d chars). Recommended: at least 32 characters", len(c.JWT.Secret))
	}
	if c.JWT.ExpireHours <= 0 {
		return fmt.Errorf("JWT expire hours must be positive")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Warning: Invalid int64 value for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
