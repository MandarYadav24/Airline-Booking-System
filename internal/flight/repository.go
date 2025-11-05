package flight

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
}

func NewRepository(db *sqlx.DB, cache *redis.Client) *Repository {
	return &Repository{DB: db, Cache: cache}
}

// GetAllFlights fetches all flights from the database.
func (r *Repository) GetAllFlights() ([]Flight, error) {
	ctx := context.Background()
	cacheKey := "flights:all"

	// Check redis cache first
	val, err := r.Cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedFlights []Flight
		if jsonErr := json.Unmarshal([]byte(val), &cachedFlights); jsonErr == nil {
			log.Printf("Flights fetched from cache")
			return cachedFlights, nil
		}
	}

	// If not in cache, fetch from database
	query := `SELECT id, airline, source, destination, departure, arrival, price, available_seats FROM flights`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query flights: %v", err)
	}
	defer rows.Close()

	var flights []Flight
	for rows.Next() {
		var f Flight
		if err := rows.Scan(&f.ID, &f.Airline, &f.Source, &f.Destination, &f.Departure, &f.Arrival, &f.Price, &f.AvailableSeats); err != nil {
			return nil, fmt.Errorf("failed to scan flight row: %v", err)
		}
		flights = append(flights, f)
	}

	// Store the result in cache
	data, _ := json.Marshal(flights)
	if err := r.Cache.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		log.Printf("Redis cache set failed: %v", err)
	}

	log.Println("Flights fetched from Db and cached")
	return flights, nil
}

// AddFlight inserts a new flight into the database.
func (r *Repository) AddFlight(f Flight) error {
	query := `
		INSERT INTO flights (airline, source, destination, departure, arrival, price, available_seats)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.Exec(query, f.Airline, f.Source, f.Destination, f.Departure, f.Arrival, f.Price, f.AvailableSeats)
	if err != nil {
		return fmt.Errorf("failed to insert flight: %v", err)
	}

	// Invalidate cache after insert
	ctx := context.Background()
	r.Cache.Del(ctx, "flights:all")

	log.Println("Redis cache invalidated after flight insert")

	return nil
}
