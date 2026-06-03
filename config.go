package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Addr     string
	Password string
	Channel  string
}

type MailtrapConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type AppConfig struct {
	Redis    RedisConfig
	Mailtrap MailtrapConfig
}

func NewConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("💡 Info: No physical .env file found. Using default container environments.")
	}

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	return &AppConfig{
		Redis: RedisConfig{
			Addr:     redisHost + ":" + redisPort,
			Password: os.Getenv("REDIS_PASSWORD"),
			Channel:  "order-notification",
		},
		Mailtrap: MailtrapConfig{
			Host:     getEnv("MAILTRAP_HOST", "sandbox.smtp.mailtrap.io"),
			Port:     getEnv("MAILTRAP_PORT", "2525"),
			Username: os.Getenv("MAILTRAP_USERNAME"),
			Password: os.Getenv("MAILTRAP_PASSWORD"),
			From:     getEnv("MAIL_FROM", "store@roofiahmad-homelabs.my.id"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
