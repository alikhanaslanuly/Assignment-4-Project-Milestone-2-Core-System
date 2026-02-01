package repository

import "time"

type User struct {
	ID        string
	Email     string
	Password  string
	FullName  string
	CreatedAt time.Time
}

type Payment struct {
	ID        string
	BookingID string
	Amount    float64
	Status    string
	CreatedAt time.Time
}
