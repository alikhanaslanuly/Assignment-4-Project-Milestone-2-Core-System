package professional_service

import (
	"database/sql"
)

type Service struct {
	DB *sql.DB
}

type Professional struct {
	ID     string
	Name   string
	Job    string
	Bio    string
	Price  float64
	Rating float64
}

func (s *Service) List() ([]Professional, error) {
	rows, err := s.DB.Query(
		"SELECT id, full_name, profession, bio, price_per_hour, rating_avg FROM professionals",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Professional
	for rows.Next() {
		var p Professional
		rows.Scan(&p.ID, &p.Name, &p.Job, &p.Bio, &p.Price, &p.Rating)
		res = append(res, p)
	}

	return res, nil
}
