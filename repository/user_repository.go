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

func (r *UserRepository) GetAll() ([]User, error) {
	rows, err := r.DB.Query("SELECT id, email, full_name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
