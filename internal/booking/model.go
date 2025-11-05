package booking

// Booking represents a flight booking record
type Booking struct {
	ID         int     `db:"id" json:"id"`
	FlightID   int     `db:"flight_id" json:"flight_id"`
	Passenger  string  `db:"passenger" json:"passenger"`
	Seats      int     `db:"seats" json:"seats"`
	TotalPrice float64 `db:"total_price" json:"total_price"`
	Status     string  `db:"status" json:"status"`
}
