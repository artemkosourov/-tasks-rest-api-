package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	DBUser       string
	DBPassword   string
	DBName       string
	DBHost       string
	DBPort       string
	JWTSecret    string
	RefreshTTL   time.Duration
	AccessTTL    time.Duration
	RateLimit    string
	Port         string
	KafkaBrokers []string
	KafkaTopic   string
}

func LoadConfig() Config {
	port := getEnv("PORT", "8080")

	if port[0] != ':' {
		port = ":" + port
	}

	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(kafkaBrokers, ",")

	return Config{
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "taskdb"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		JWTSecret:    getEnv("JWT_SECRET", "secret_key"),
		RefreshTTL:   7 * 24 * time.Hour,
		AccessTTL:    15 * time.Minute,
		RateLimit:    getEnv("RATE_LIMIT", "5-M"),
		Port:         port,
		KafkaBrokers: brokers,
		KafkaTopic:   getEnv("KAFKA_TOPIC", "task-events"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
