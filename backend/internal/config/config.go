package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Dsn string
}

type ServerConfig struct {
	Listen string
}

type SecretConfig struct {
	JwtSecret  string
	HashSecret string
}

type Config struct {
	Db       DbConfig
	Server   ServerConfig
	Secret   SecretConfig
	LogLevel int
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbHost := getEnvRequired("DB_HOST")
	dbPort := getEnvRequired("DB_PORT")
	dbUser := getEnvRequired("DB_USER")
	dbPassword := getEnvRequired("DB_PASSWORD")
	dbName := getEnvRequired("DB_NAME")

	dbDsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	logLevelStr := getEnv("LOG_LEVEL", "0")
	logLevel, err := strconv.Atoi(logLevelStr)
	if err != nil {
		logLevel = 0
	}

	return &Config{
		Db: DbConfig{
			Dsn: dbDsn,
		},
		Server: ServerConfig{
			Listen: getEnv("SERVER_LISTEN", ":8080"),
		},
		Secret: SecretConfig{
			JwtSecret:  getEnvRequired("JWT_SECRET"),
			HashSecret: getEnvRequired("HASH_SECRET"),
		},
		LogLevel: logLevel,
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("Не найден ключ " + key)
	}
	return value
}
