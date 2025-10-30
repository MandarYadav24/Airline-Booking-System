package booking

import (
	"encoding/json"
	"log"
	"net/http"

	"airline-booking/pkg/kafka"
)

type Handler struct {
	Repo     *Repository
	Producer *kafka.Producer
	Topic    string
}

func NewHandler(repo *Repository, producer *kafka.Producer, topic string) *Handler {
	return &Handler{Repo: repo, Producer: producer, Topic: topic}
}

// AddBooking handles booking creation
func (h *Handler) AddBooking(w http.ResponseWriter, r *http.Request) {
	var b Booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.Repo.AddBooking(b); err != nil {
		http.Error(w, "Failed to save booking", http.StatusInternalServerError)
		log.Println("DB error:", err)
		return
	}

	// Publish Kafka event
	event, _ := json.Marshal(b)
	err := h.Producer.SendMessage(h.Topic, "booking_created", string(event))
	if err != nil {
		log.Printf("Kafka publish error: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking created successfully"})
}

// GetBookings returns all bookings
func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.Repo.GetAllBookings()
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
