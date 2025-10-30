package booking

// Booking represents a flight booking record
type Booking struct {
	ID         int     `json:"id"`
	FlightID   int     `json:"flight_id"`
	Passenger  string  `json:"passenger"`
	Seats      int     `json:"seats"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}
