package payment_service

import (
	"eventify/repository"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Repo *repository.PaymentRepository
}

func (s *Service) CreatePayment(bookingID string, amount float64) error {
	payment := repository.Payment{
		ID:        uuid.New().String(),
		BookingID: bookingID,
		Amount:    amount,
		Status:    "pending",
	}

	err := s.Repo.Create(payment)
	if err != nil {
		return err
	}

	go s.processPayment(payment.ID)

	return nil
}

func (s *Service) processPayment(paymentID string) {
	time.Sleep(2 * time.Second)
	s.Repo.UpdateStatus(paymentID, "paid")
}
