package configs

import (
	"errors"
	"os"
)


func LoadDatabaseConfig() (pg_user string, pg_password string, pg_db string, host string, err error) {
	pg_user = os.Getenv("POSTGRES_USER")
	pg_password = os.Getenv("POSTGRES_PASSWORD")
	pg_db = os.Getenv("POSTGRES_DB")
	host = os.Getenv("HOST")

	if host == "" || pg_user == "" || pg_password == "" || pg_db == "" {
		return "", "", "", "", errors.New("database configuration not found")
	}

	return pg_user, pg_password, pg_db, host, nil
}

func GetPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		return ":8000"
	}

	return ":" + port
}

func GetRabbitMQUrl() string {
	rabbitmq_url := os.Getenv("RABBITMQ_URL")

	if rabbitmq_url == "" {
		return "amqp://guest:guest@localhost:5672/"
	}

	return rabbitmq_url
}
