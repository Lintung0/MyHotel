package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	JWTExpHours int
}

func LoadConfig() *Config{
	err := godotenv.Load()
	if err != nil {
		log.Println("Perhatian: file .env tidak ditemukan, menggunakan environment variables sistem.")
	}

	expHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		expHours = 24
	}

	return &Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET_KEY"),
		JWTExpHours: expHours,
	}
}