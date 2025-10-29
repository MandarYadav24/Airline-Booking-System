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
	Kafka    KafkaConfig
	Redis    RedisConfig
}

/*---------------Postgres-----------------*/
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

/*-------------------- Kafka --------------------*/
type KafkaConfig struct {
	Brokers          []string
	Topic            string
	GroupID          string `mapstructure:"groupId"`
	AutoOffsetReset  string `mapstructure:"autoOffsetReset"`
	EnableAutoCommit bool   `mapstructure:"enableAutoCommit"`
	Producer         struct {
		Retries        int
		BatchTimeoutMS int `mapstructure:"batchTimeoutMS"`
	}
	Consumer struct {
		InitialOffset string `mapstructure:"initialOffset"`
	}
}

/*
type KafkaConfig struct {
	Brokers  []string
	Topic    string
	GroupID  string
	Security struct {
		EnableTLS bool
		Username  string
		Password  string
	}
	Producer struct {
		Retries        int
		BatchTimeoutMS int
	}
	Consumer struct {
		InitialOffset string
	}
}
*/
/*-------------------- Redis --------------------*/
type RedisConfig struct {
	Address      string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	ReadTimeout  string
	WriteTimeout string
	DialTimeout  string
}

func LoadConfig() (*Config, error) {
	absHome := os.Getenv("ABS_HOME")
	if absHome == "" {
		return nil, fmt.Errorf("ABS_HOME environment variable is not set")
	}

	configDir := filepath.Join(absHome, "configs")

	//var cfg Config {}
	cfg := &Config{}

	// --- Postgres ---
	dbViper := viper.New()
	dbViper.SetConfigName("dbconfig")
	dbViper.SetConfigType("yaml")
	dbViper.AddConfigPath(configDir)
	if err := dbViper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading dbconfig.yml: %v", err)
	}
	if err := dbViper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		log.Fatalf("Error decoding Postgres config: %v", err)
	}

	// --- Kafka ---
	kafkaViper := viper.New()
	kafkaViper.SetConfigName("kafkaconfig")
	kafkaViper.SetConfigType("yaml")
	kafkaViper.AddConfigPath(configDir)
	if err := kafkaViper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading kafkaconfig.yml: %v", err)
	}
	if err := kafkaViper.UnmarshalKey("kafka", &cfg.Kafka); err != nil {
		log.Fatalf("Error decoding Kafka config: %v", err)
	}

	// Inside LoadConfig() function, after loading Kafka config:
	if len(cfg.Kafka.Brokers) == 0 {
		log.Printf("Warning: No Kafka brokers configured")
	}
	//log.Printf("Loaded Kafka config: brokers=%v, topic=%s", cfg.Kafka.Brokers, cfg.Kafka.Topic)

	// --- Redis ---
	redisViper := viper.New()
	redisViper.SetConfigName("redisconfig")
	redisViper.SetConfigType("yaml")
	redisViper.AddConfigPath(configDir)
	if err := redisViper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading redisconfig.yml: %v", err)
	}
	if err := redisViper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		log.Fatalf("Error decoding Redis config: %v", err)
	}

	return cfg, nil

	/*
		// -------------------- Load Postgres --------------------
		viper.SetConfigName("dbconfig")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configDir)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading dbconfig.yml: %v", err)
		}
		if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
			log.Fatalf("Error decoding Postgres config: %v", err)
		}

		// -------------------- Load Kafka --------------------
		viper.SetConfigName("kafkaconfig")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading kafkaconfig.yml: %v", err)
		}
		if err := viper.UnmarshalKey("kafka", &cfg.Kafka); err != nil {
			log.Fatalf("Error decoding Kafka config: %v", err)
		}

		// -------------------- Load Redis --------------------
		viper.SetConfigName("redisconfig")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading redisconfig.yml: %v", err)
		}
		if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
			log.Fatalf("Error decoding Redis config: %v", err)
		}
	*/
	//return &cfg, nil
}
