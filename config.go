package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ConfigDB struct {
	User     string
	Password string
	Dbname   string
	Port     string
}

type Config struct {
    DB ConfigDB
}

func loadConfig() (Config, error) {
	var config Config
	var exists bool

	// Открываем env-шку
	if err := godotenv.Load(); err != nil {
		return config, fmt.Errorf("ошибка с .env: %v", err)
	}

	// Получаем значения переменных окружения
	// Если чего то нет, вызываем ошибку
	if config.DB.User, exists = os.LookupEnv("user"); !exists {
		return config, fmt.Errorf("не установлена переменная окруженния user")
	}

	if config.DB.Password, exists = os.LookupEnv("password"); !exists {
		return config, fmt.Errorf("не установлена переменная окруженния password")
	}

	if config.DB.Dbname, exists = os.LookupEnv("dbname"); !exists {
		return config, fmt.Errorf("не установлена переменная окруженния dbname")
	}

	if config.DB.Port, exists = os.LookupEnv("port"); !exists {
		return config, fmt.Errorf("не установлена переменная окруженния port")
	}

	return config, nil
}