package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	DB    *sqlx.DB
	Cache *redis.Client
	Ctx   context.Context
}

func NewRepository(db *sqlx.DB, cache *redis.Client) *Repository {
	return &Repository{
		DB:    db,
		Cache: cache,
		Ctx:   context.Background(),
	}
}

// AddBooking inserts a new booking into DB and caches it in Redis to prevent duplicates.
func (r *Repository) AddBooking(b Booking) error {
	cacheKey := fmt.Sprintf("booking:%s:%d", b.Passenger, b.FlightID)

	// Check if user already booked this flight (from cache)
	exists, err := r.Cache.Exists(r.Ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Redis check failed: %v", err)
	}
	if exists > 0 {
		return fmt.Errorf("duplicate booking detected for user %s and flight %d", b.Passenger, b.FlightID)
	}

	// Insert booking into DB
	query := `
		INSERT INTO bookings (flight_id, passenger, seats, total_price, status)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = r.DB.Exec(query, b.FlightID, b.Passenger, b.Seats, b.TotalPrice, b.Status)
	if err != nil {
		return fmt.Errorf("failed to insert booking: %w", err)
	}

	// Cache the booking for 1 hour to prevent duplicate submissions
	data, _ := json.Marshal(b)
	err = r.Cache.Set(r.Ctx, cacheKey, data, time.Hour).Err()
	if err != nil {
		log.Printf("Failed to cache booking in Redis: %v", err)
	} else {
		log.Printf("Booking cached in Redis: %s", cacheKey)
	}

	return nil
}

// GetAllBookings retrieves all bookings, using Redis cache if available.
func (r *Repository) GetAllBookings() ([]Booking, error) {
	cacheKey := "bookings:all"

	// Try fetching from Redis cache first
	cachedData, err := r.Cache.Get(r.Ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedBookings []Booking
		if err := json.Unmarshal([]byte(cachedData), &cachedBookings); err == nil {
			log.Println("Bookings served from Redis cache")
			return cachedBookings, nil
		}
	}

	// Fetch from database if cache miss
	rows, err := r.DB.Queryx(`SELECT id, flight_id, passenger, seats, total_price, status FROM bookings`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings: %w", err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.StructScan(&b); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	// Store in Redis cache for faster future access (30s TTL)
	if len(bookings) > 0 {
		data, _ := json.Marshal(bookings)
		err = r.Cache.Set(r.Ctx, cacheKey, data, 30*time.Second).Err()
		if err != nil {
			log.Printf("Failed to cache bookings list: %v", err)
		} else {
			log.Println("Cached bookings list in Redis")
		}
	}

	return bookings, nil
}
