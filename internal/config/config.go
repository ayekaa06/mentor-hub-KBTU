// Package config загружает и хранит конфигурацию приложения.
// Источники (приоритет по убыванию): переменные окружения → .env файл → defaults.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config — корневая структура конфигурации.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig — подключение к PostgreSQL.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig — параметры JWT токенов.
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load читает конфигурацию и возвращает Config.
func Load() (*Config, error) {
	setDefaults()

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// .env не обязателен — переменные окружения тоже работают
	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			AccessTokenTTL:  viper.GetDuration("JWT_ACCESS_TTL"),
			RefreshTokenTTL: viper.GetDuration("JWT_REFRESH_TTL"),
		},
	}

	return cfg, nil
}

// DSN возвращает строку подключения к PostgreSQL для pgx.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

func setDefaults() {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "10s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_ACCESS_TTL", "15m")
	viper.SetDefault("JWT_REFRESH_TTL", "168h") // 7 дней
}
