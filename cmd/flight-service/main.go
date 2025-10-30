package main

import (
	"log"
	"net/http"

	"airline-booking/internal/flight"
	"airline-booking/pkg/config"
	"airline-booking/pkg/db"
	"airline-booking/pkg/kafka"
	"airline-booking/pkg/redis"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to PostgreSQL
	pg, err := db.ConnectPostgres(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pg.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Connect to Redis
	redisClient := redis.NewRedisClient(&cfg.Redis)
	defer redisClient.Close()
	log.Println("✅ Connected to Redis")

	// Initialize Kafka Producer (Sarama)
	producer, err := kafka.NewProducer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka producer: %v", err)
	}
	defer producer.Close()
	log.Println("✅ Connected to Kafka Producer")

	// Initialize Repository and Handler
	repo := flight.NewRepository(pg)
	handler := flight.NewHandler(repo, producer, cfg.Kafka.Topic)

	// Define HTTP routes
	http.HandleFunc("/flights", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetFlights(w, r)
		case http.MethodPost:
			handler.AddFlight(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("🚀 Flight service started successfully — all connections active.")
	log.Println("🌐 Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
