package main

import (
	"database/sql"
	"encoding/json"
	"eventify/availability_service"
	"eventify/booking_service"
	"fmt"
	"log"
	"net/http"
	"time"

	"eventify/payment_service"
	"eventify/repository"
	"eventify/user_service"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=09022107 dbname=eventify sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка конфигурации БД:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("База недоступна:", err)
	}
	fmt.Println("EVENTIFY backend успешно подключен")

	userRepo := &repository.UserRepository{DB: db}
	paymentRepo := &repository.PaymentRepository{DB: db}

	userService := &user_service.Service{Repo: userRepo}
	paymentService := &payment_service.Service{Repo: paymentRepo}

	userHandler := &UserHandler{Service: userService}

	http.HandleFunc("/register", userHandler.Register)

	http.HandleFunc("/payments/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			BookingID string  `json:"booking_id"`
			Amount    float64 `json:"amount"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := paymentService.CreatePayment(req.BookingID, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Payment processing started"})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "alive", "database": "connected"})
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	avaSrv := &availability_service.Service{DB: db}

	bookSrv := &booking_service.Service{
		DB:           db,
		Availability: avaSrv,
	}

	http.HandleFunc("/bookings/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			CustomerID     string  `json:"customer_id"`
			ProfessionalID string  `json:"professional_id"`
			EventDate      string  `json:"event_date"`
			EventTime      string  `json:"event_time"`
			EventLocation  string  `json:"event_location"`
			Price          float64 `json:"price"`
			StartTime      string  `json:"start_time"`
			EndTime        string  `json:"end_time"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		eventDate, _ := time.Parse("2006-01-02", req.EventDate)
		start, _ := time.Parse(time.RFC3339, req.StartTime)
		end, _ := time.Parse(time.RFC3339, req.EndTime)

		id, err := bookSrv.Create(booking_service.CreateBookingRequest{
			CustomerID:     req.CustomerID,
			ProfessionalID: req.ProfessionalID,
			EventDate:      eventDate,
			EventTime:      req.EventTime,
			Location:       req.EventLocation,
			Price:          req.Price,
			StartTime:      start,
			EndTime:        end,
		})

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"booking_id": id,
			"status":     "pending",
		})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
