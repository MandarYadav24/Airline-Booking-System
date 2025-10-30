package flight

// Flight represents the structure of a flight record.
type Flight struct {
	ID             int     `json:"id"`
	Airline        string  `json:"airline"`
	Source         string  `json:"source"`
	Destination    string  `json:"destination"`
	Departure      string  `json:"departure"`
	Arrival        string  `json:"arrival"`
	Price          float64 `json:"price"`
	AvailableSeats int     `json:"available_seats"`
}
