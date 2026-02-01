package repository

import "database/sql"

type PaymentRepository struct {
	DB *sql.DB
}

func (r *PaymentRepository) Create(p Payment) error {
	_, err := r.DB.Exec(
		"INSERT INTO payments (id, booking_id, amount, status) VALUES ($1, $2, $3, $4)",
		p.ID, p.BookingID, p.Amount, p.Status,
	)
	return err
}

func (r *PaymentRepository) UpdateStatus(id, status string) error {
	_, err := r.DB.Exec(
		"UPDATE payments SET status=$1 WHERE id=$2",
		status, id,
	)
	return err
}
