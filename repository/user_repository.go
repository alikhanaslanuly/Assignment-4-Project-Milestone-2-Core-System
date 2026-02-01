package repository

import "database/sql"

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4)",
		user.ID, user.Email, user.Password, user.FullName,
	)
	return err
}
