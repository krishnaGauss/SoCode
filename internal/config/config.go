package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Host     string
	Port     int
	GRPCPort int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:     getEnv("SERVER_HOST", "localhost"),
			Port:     getEnvInt("SERVER_PORT", 8080),
			GRPCPort: getEnvInt("GRPC_PORT", 9090),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			Database: getEnv("DB_NAME", "logs"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value:=os.Getenv(key); value!=""{
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int{
	if value:=os.Getenv(key); value!=""{
		if int_val, err:=strconv.Atoi(value); err==nil{
			return int_val
		}
	}

	return defaultValue
	
}