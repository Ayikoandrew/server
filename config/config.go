package config

import (
	"errors"
	"fmt"
	"log/slog"
)

type Config struct {
	DBPassword string `json:"db_password" env:"DB_PASSWORD"`
	DBUser     string `json:"db_user" env:"DB_USER"`
	DBName     string `json:"db_name" env:"DB_NAME"`
	DBPort     string `json:"db_port" env:"DB_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
}

func NewConfig(db_password, db_user, db_name, db_port, db_host string) (*Config, error) {
	if db_password == "" || db_user == "" || db_name == "" || db_port == "" {
		slog.Info("missing required database configuration parameters")
		return nil, errors.New("missing required database configuration parameters")
	}
	return &Config{
		DBPassword: db_password,
		DBUser:     db_user,
		DBName:     db_name,
		DBPort:     db_port,
		DBHost:     db_host,
	}, nil
}

func (c *Config) PGXDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s password=%s user=%s sslmode=disable", c.DBHost, c.DBPort, c.DBName, c.DBPassword, c.DBUser)
}
