package review_service

import (
	"database/sql"

	"github.com/google/uuid"
)

type Service struct {
	DB *sql.DB
}

type ReviewRequest struct {
	ProfessionalID string
	CustomerID     string
	Rating         int
	Comment        string
}

func (s *Service) AddReview(r ReviewRequest) error {
	_, err := s.DB.Exec(
		`INSERT INTO reviews (id, professional_id, customer_id, rating, comment)
		 VALUES ($1,$2,$3,$4,$5)`,
		uuid.New().String(), r.ProfessionalID, r.CustomerID, r.Rating, r.Comment,
	)
	return err
}
