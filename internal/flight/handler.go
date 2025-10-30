package flight

import (
	"encoding/json"
	"log"
	"net/http"

	"airline-booking/pkg/kafka"
)

// Handler holds dependencies for flight HTTP routes.
type Handler struct {
	Repo     *Repository
	Producer *kafka.Producer
	Topic    string
}

// NewHandler creates a new flight handler.
func NewHandler(repo *Repository, producer *kafka.Producer, topic string) *Handler {
	return &Handler{
		Repo:     repo,
		Producer: producer,
		Topic:    topic,
	}
}

// GetFlights returns all available flights.
func (h *Handler) GetFlights(w http.ResponseWriter, r *http.Request) {
	flights, err := h.Repo.GetAllFlights()
	if err != nil {
		http.Error(w, "Failed to fetch flights", http.StatusInternalServerError)
		log.Printf("Error fetching flights: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flights)
}

// AddFlight adds a new flight and publishes an event to Kafka.
func (h *Handler) AddFlight(w http.ResponseWriter, r *http.Request) {
	var f Flight
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert into Postgres
	if err := h.Repo.AddFlight(f); err != nil {
		http.Error(w, "Failed to add flight", http.StatusInternalServerError)
		log.Printf("DB insert error: %v", err)
		return
	}

	// Publish Kafka event
	eventData, _ := json.Marshal(f)
	err := h.Producer.SendMessage(h.Topic, "flight_created", string(eventData))
	if err != nil {
		log.Printf("Failed to publish Kafka message: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Flight added successfully",
	})
}
