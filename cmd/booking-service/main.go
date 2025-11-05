package main

import (
	"airline-booking/internal/booking"
	"airline-booking/pkg/config"
	"airline-booking/pkg/db"
	"airline-booking/pkg/kafka"
	"airline-booking/pkg/redis"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pg, err := db.ConnectPostgres(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pg.Close()
	log.Println("✅ Connected to PostgreSQL")

	redisClient := redis.NewRedisClient(&cfg.Redis)
	defer redisClient.Close()
	log.Println("✅ Connected to Redis")

	producer, err := kafka.NewProducer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("Kafka producer connection failed: %v", err)
	}
	defer producer.Close()
	log.Println("✅ Connected to Kafka")

	repo := booking.NewRepository(pg, redisClient.GetClient())
	handler := booking.NewHandler(repo, producer, cfg.Kafka.Topic)

	http.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.AddBooking(w, r)
		case http.MethodGet:
			handler.GetBookings(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("🚀 Booking service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
