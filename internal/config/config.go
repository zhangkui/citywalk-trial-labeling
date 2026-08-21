package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment    string
	ListenAddr     string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	ShutdownAfter  time.Duration
	MySQLHost      string
	MySQLPort      int
	MySQLDatabase  string
	MySQLUser      string
	MySQLPassword  string
	RedisHost      string
	RedisPort      int
	RedisPassword  string
	JWTSecret      string
	AdminUsername  string
	AdminPassword  string
	DataDir        string
	CORSAllowOrigin string
}

func Load() (*Config, error) {
	_ = loadDotEnv(filepath.Join(mustWD(), ".env"))
	cfg := &Config{
		Environment:    getEnv("APP_ENV", "development"),
		ListenAddr:     getEnv("APP_LISTEN_ADDR", ":8080"),
		ReadTimeout:    getDurationEnv("APP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:   getDurationEnv("APP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:    getDurationEnv("APP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownAfter:  getDurationEnv("APP_SHUTDOWN_TIMEOUT", 30*time.Second),
		MySQLHost:      getEnv("MYSQL_HOST", "mysql"),
		MySQLPort:      getIntEnv("MYSQL_PORT", 3306),
		MySQLDatabase:  getEnv("MYSQL_DATABASE", "citywalk"),
		MySQLUser:      getEnv("MYSQL_USER", "citywalk"),
		MySQLPassword:  getEnv("MYSQL_PASSWORD", "citywalk2024"),
		RedisHost:      getEnv("REDIS_HOST", "redis"),
		RedisPort:      getIntEnv("REDIS_PORT", 6379),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		JWTSecret:      getEnv("JWT_SECRET", "change-me"),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "Admin123!"),
		DataDir:        getEnv("DATA_DIR", "./data"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "*"),
	}
	if cfg.JWTSecret == "change-me" {
		return nil, fmt.Errorf("JWT_SECRET must be configured")
	}
	return cfg, nil
}

func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local", c.MySQLUser, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDatabase)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}

func mustWD() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")
		_ = os.Setenv(key, value)
	}
	return scanner.Err()
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return parsed
}
