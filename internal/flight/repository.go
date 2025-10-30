package flight

import (
	//"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{DB: db}
}

// GetAllFlights fetches all flights from the database.
func (r *Repository) GetAllFlights() ([]Flight, error) {
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
	return nil
}
