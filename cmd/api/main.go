package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benokwulu-lgtm/triax-habit-tracker/internal/auth"
	"github.com/benokwulu-lgtm/triax-habit-tracker/internal/habit"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func buildConnStr() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	dbname := getEnv("DB_NAME", "triax_habit_tracker")
	sslmode := getEnv("DB_SSLMODE", "disable")
	password := os.Getenv("DB_PASSWORD") // no fallback — may legitimately be empty (Mac trust auth)

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)
	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}
	return connStr
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	db, err := sql.Open("postgres", buildConnStr())
	if err != nil {
		log.Fatal("failed to open database connection:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	fmt.Println("Connected to PostgreSQL successfully")

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	repo := habit.NewRepository(db)
	service := habit.NewService(repo)
	handler := habit.NewHandler(service)

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/register", authHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)
	http.HandleFunc("/api/habits", authHandler.RequireAuth(handler.ListCreate))
	http.HandleFunc("/api/habits/", authHandler.RequireAuth(handler.Detail))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
