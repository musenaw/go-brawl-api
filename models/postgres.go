package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		Database: "brawl",
		SSLMode:  "disable",
	}
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Open(config PostgresConfig) error {
	var err error
	DB, err = gorm.Open(postgres.Open(config.String()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	return nil
}
