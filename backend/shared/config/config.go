package config

import (
	"os"
	"strconv"
	"time"
)

// Config アプリケーション設定
type Config struct {
	// Server
	Port          int
	Environment   string
	JWTSecretKey  string
	JWTExpiration time.Duration

	// Database
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	// Redis
	RedisAddr string
	RedisPass string

	// Elasticsearch
	ElasticsearchURL string

	// AWS
	AWSRegion      string
	AWSEndpointURL string
	AWSAccessKeyID string
	AWSSecretKey   string

	// Cognito
	CognitoUserPoolID string
	CognitoClientID   string
	CognitoRegion     string
}

// LoadConfig 設定をロード
func LoadConfig() *Config {
	return &Config{
		Port:              getEnvInt("PORT", 8000),
		Environment:       getEnv("ENVIRONMENT", "development"),
		JWTSecretKey:      getEnv("JWT_SECRET_KEY", "your-secret-key"),
		JWTExpiration:     time.Duration(getEnvInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnvInt("DB_PORT", 3306),
		DBName:            getEnv("DB_NAME", "ecsite"),
		DBUser:            getEnv("DB_USER", "ecuser"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:         getEnv("REDIS_PASSWORD", ""),
		ElasticsearchURL:  getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		AWSEndpointURL:    getEnv("AWS_ENDPOINT_URL", ""),
		AWSAccessKeyID:    getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey:      getEnv("AWS_SECRET_ACCESS_KEY", ""),
		CognitoUserPoolID: getEnv("COGNITO_USER_POOL_ID", ""),
		CognitoClientID:   getEnv("COGNITO_CLIENT_ID", ""),
		CognitoRegion:     getEnv("COGNITO_REGION", "us-east-1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}
