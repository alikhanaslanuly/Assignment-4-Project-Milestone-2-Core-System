package availability_service

import (
	"database/sql"
	"time"
)

type Service struct {
	DB *sql.DB
}

func (s *Service) IsAvailable(profID string, start, end time.Time) (bool, error) {
	row := s.DB.QueryRow(
		`SELECT COUNT(*) FROM availability
		 WHERE professional_id=$1
		   AND $2 < end_time
		   AND $3 > start_time`,
		profID, start, end,
	)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}
