package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"l0/cmd/internal/cache"
	"l0/cmd/internal/db"
	"l0/cmd/internal/handlers"
	"l0/cmd/internal/kafka"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Order Service...")

	// Подключение к БД
	dsn := getEnv("DB_DSN", "")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}
	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	// Инициализация кэша
	if err := cache.LoadFromDB(dbConn); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}
	log.Println("Cache initialized")

	// Запуск Kafka consumer
	kafkaConfig := kafka.ConsumerConfig{
		DB:     dbConn,
		Broker: getEnv("KAFKA_BROKER", ""),
		Topic:  getEnv("KAFKA_TOPIC", ""),
	}
	if kafkaConfig.Broker == "" || kafkaConfig.Topic == "" {
		log.Fatal("KAFKA_BROKER and KAFKA_TOPIC environment variables are required")
	}
	go kafka.StartConsumer(kafkaConfig)

	// Настройка HTTP сервера
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	r.HandleFunc("/order/{id}", handlers.GetOrderHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates")))

	httpPort := getEnv("HTTP_PORT", ":8081")
	srv := &http.Server{
		Handler:      r,
		Addr:         httpPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("HTTP server listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
