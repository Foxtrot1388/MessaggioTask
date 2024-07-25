package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	LogLevel         int
	DB               string
	ConnectionString string
	KafkaHost        string
	KafkaTopic       string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}
}

func GetConfig() *AppConfig {
	return &AppConfig{
		LogLevel: getEnvAsInt("LOG_LEVEL", -4),
		DB:       getEnv("POSTGRES_DB", "messaggio"),
		ConnectionString: fmt.Sprintf("postgres://%s:%s/%s?sslmode=disable&user=%s&password=%s",
			getEnv("POSTGRES_HOST", "localhost"), getEnv("POSTGRES_PORT", "5432"), getEnv("POSTGRES_DB", "messaggio"),
			getEnv("POSTGRES_USER", "root"), getEnv("POSTGRES_PASSWORD", "root")),
		KafkaHost:  getEnv("KAFKA_HOST", "localhost:9092"),
		KafkaTopic: getEnv("KAFKA_TOPIC", "tasks"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
