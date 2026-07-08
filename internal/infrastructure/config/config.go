package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName               string   `mapstructure:"APP_NAME"`
	AppVersion            string   `mapstructure:"APP_VERSION"`
	AppEnv                string   `mapstructure:"APP_ENV"`
	AppPort               string   `mapstructure:"APP_PORT"`
	AppDebug              bool     `mapstructure:"APP_DEBUG"`
	AppCorsAllowedOrigins []string `mapstructure:"APP_CORS_ALLOWED_ORIGINS"`

	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_NAME"`
	DbSslMode  string `mapstructure:"DB_SSLMODE"`
	DbTz       string `mapstructure:"DB_TZ"`

	JwtSecret            string `mapstructure:"JWT_SECRET"`
	JwtExpiryHour        int    `mapstructure:"JWT_EXPIRY_HOUR"`
	JwtRefreshExpiryHour int    `mapstructure:"JWT_REFRESH_EXPIRY_HOUR"`
}

func LoadConfig() (*Config, error) {
	// Tentukan file dan tipe config
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Baca file konfigurasi
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file not found, reading from environment variables")
	}

	// Ambil value raw untuk menampung parsing CORS origins (karena di .env berupa string comma-separated)
	var raw struct {
		AppName               string `mapstructure:"APP_NAME"`
		AppVersion            string `mapstructure:"APP_VERSION"`
		AppEnv                string `mapstructure:"APP_ENV"`
		AppPort               string `mapstructure:"APP_PORT"`
		AppDebug              bool   `mapstructure:"APP_DEBUG"`
		AppCorsAllowedOrigins string `mapstructure:"APP_CORS_ALLOWED_ORIGINS"`
		DbHost                string `mapstructure:"DB_HOST"`
		DbPort                string `mapstructure:"DB_PORT"`
		DbUser                string `mapstructure:"DB_USER"`
		DbPassword            string `mapstructure:"DB_PASSWORD"`
		DbName                string `mapstructure:"DB_NAME"`
		DbSslMode             string `mapstructure:"DB_SSLMODE"`
		DbTz                  string `mapstructure:"DB_TZ"`
		JwtSecret             string `mapstructure:"JWT_SECRET"`
		JwtExpiryHour         int    `mapstructure:"JWT_EXPIRY_HOUR"`
		JwtRefreshExpiryHour  int    `mapstructure:"JWT_REFRESH_EXPIRY_HOUR"`
	}

	if err := viper.Unmarshal(&raw); err != nil {
		return nil, err
	}

	var origins []string
	if raw.AppCorsAllowedOrigins != "" {
		origins = strings.Split(raw.AppCorsAllowedOrigins, ",")
		for i, o := range origins {
			origins[i] = strings.TrimSpace(o)
		}
	}

	return &Config{
		AppName:               raw.AppName,
		AppVersion:            raw.AppVersion,
		AppEnv:                raw.AppEnv,
		AppPort:               raw.AppPort,
		AppDebug:              raw.AppDebug,
		AppCorsAllowedOrigins: origins,
		DbHost:                raw.DbHost,
		DbPort:                raw.DbPort,
		DbUser:                raw.DbUser,
		DbPassword:            raw.DbPassword,
		DbName:                raw.DbName,
		DbSslMode:             raw.DbSslMode,
		DbTz:                  raw.DbTz,
		JwtSecret:             raw.JwtSecret,
		JwtExpiryHour:         raw.JwtExpiryHour,
		JwtRefreshExpiryHour:  raw.JwtRefreshExpiryHour,
	}, nil
}
