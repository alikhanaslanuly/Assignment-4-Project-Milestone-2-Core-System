package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"eventify/config"
	"eventify/middleware"
	"eventify/payment_service"
	"eventify/repository"
	"eventify/user_service"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Database Connection
	db, err := sql.Open("postgres", cfg.DBConnStr)
	if err != nil {
		log.Fatal("Ошибка конфигурации БД:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("База недоступна:", err)
	}
	fmt.Println("EVENTIFY backend успешно подключен")

	// 3. Initialize Repositories and Services
	userRepo := &repository.UserRepository{DB: db}
	paymentRepo := &repository.PaymentRepository{DB: db}

	userService := &user_service.Service{Repo: userRepo}
	paymentService := &payment_service.Service{Repo: paymentRepo}

	userHandler := &UserHandler{Service: userService}

	// 4. Setup Routes
	mux := http.NewServeMux()

	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/users", userHandler.ListUsers)

	mux.HandleFunc("/payments/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			BookingID string  `json:"booking_id"`
			Amount    float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		err := paymentService.CreatePayment(req.BookingID, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Payment processing started"})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "alive", "database": "connected"})
	})

	// 5. Apply Middleware
	// Wrap the entire mux with Logging and Auth Mock middleware
	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.AuthMockMiddleware(handler)

	fmt.Printf("Сервер запущен на http://localhost%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, handler))
}
