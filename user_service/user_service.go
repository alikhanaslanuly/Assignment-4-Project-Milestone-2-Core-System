package user_service

import (
	"eventify/repository"

	"github.com/google/uuid"
)

type Service struct {
	Repo *repository.UserRepository
}

func (s *Service) Register(email, password, name string) error {
	user := repository.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: password,
		FullName: name,
	}

	return s.Repo.Create(user)
}
