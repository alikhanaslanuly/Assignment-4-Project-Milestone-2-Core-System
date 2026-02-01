package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
