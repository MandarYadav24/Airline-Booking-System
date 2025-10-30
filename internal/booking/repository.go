package booking

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) AddBooking(b Booking) error {
	query := `
		INSERT INTO bookings (flight_id, passenger, seats, total_price, status)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, b.FlightID, b.Passenger, b.Seats, b.TotalPrice, b.Status)
	if err != nil {
		return fmt.Errorf("failed to insert booking: %w", err)
	}
	return nil
}

func (r *Repository) GetAllBookings() ([]Booking, error) {
	rows, err := r.DB.Query(`SELECT id, flight_id, passenger, seats, total_price, status FROM bookings`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings: %w", err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.FlightID, &b.Passenger, &b.Seats, &b.TotalPrice, &b.Status); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}
