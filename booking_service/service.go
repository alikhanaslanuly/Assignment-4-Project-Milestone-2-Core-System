package booking_service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AvailabilityChecker interface {
	IsAvailable(id string, start, end time.Time) (bool, error)
}

type Service struct {
	DB           *sql.DB
	Availability AvailabilityChecker
}

type CreateBookingRequest struct {
	CustomerID     string
	ProfessionalID string
	EventDate      time.Time
	EventTime      string
	Location       string
	Price          float64
	StartTime      time.Time
	EndTime        time.Time
}

func (s *Service) Create(req CreateBookingRequest) (string, error) {
	ok, err := s.Availability.IsAvailable(req.ProfessionalID, req.StartTime, req.EndTime)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("professional is not available")
	}

	id := uuid.New().String()

	_, err = s.DB.Exec(
		`INSERT INTO bookings (id, customer_id, event_date, event_time, event_location, total_price)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		id, req.CustomerID, req.EventDate, req.EventTime, req.Location, req.Price,
	)
	if err != nil {
		return "", err
	}

	_, _ = s.DB.Exec(
		`INSERT INTO availability (id, professional_id, start_time, end_time)
		 VALUES ($1,$2,$3,$4)`,
		uuid.New().String(), req.ProfessionalID, req.StartTime, req.EndTime,
	)

	return id, nil
}
