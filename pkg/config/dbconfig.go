package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	// Implementation to load configuration from YAML file
	absHome := os.Getenv("ABS_HOME")
	if absHome == "" {
		return nil, fmt.Errorf("ABS_HOME environment variable is not set")
	}

	viper.SetConfigName("dbconfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(absHome, "configs"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	return cfg, nil
}
