package config

import (
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	HTTP        HTTP
	Postgres    Postgres
	ExternalAPI API
}

type HTTP struct {
	Host string
	Port string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type API struct {
	URL string
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &AppConfig{
		HTTP: HTTP{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},

		Postgres: Postgres{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},

		ExternalAPI: API{
			URL: os.Getenv("EXTERNAL_API_URL"),
		},
	}

	return config, nil
}
