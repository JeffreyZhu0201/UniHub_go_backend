package config

import (
	"os"

	"github.com/spf13/viper"
)

// Config holds application configuration loaded from yaml.
type Config struct {
	Server struct {
		Port string
		Mode string
	}
	Database struct {
		Driver   string
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		Source   string
	}
	Redis struct {
		Addr     string
		Password string
	}
	JWT struct {
		Secret          string
		ExpirationHours int
	}
}

// Load loads configuration from CONFIG_PATH env or defaults to configs/config.yaml.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	if envPath := getenv("CONFIG_PATH", ""); envPath != "" {
		v.SetConfigFile(envPath)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// getenv wraps os.Getenv with default value.
func getenv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
