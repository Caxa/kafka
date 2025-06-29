package main

import (
	"log"
	"net/http"
	"time"

	"l0/cmd/internal/cache"
	"l0/cmd/internal/db"
	"l0/cmd/internal/handlers"
	"l0/cmd/internal/kafka"

	"github.com/gorilla/mux"
)

func main() {
	// Настройка логгера
	log.Println("Starting Order Service...")

	// Подключение к БД
	dsn := "host=localhost port=5432 user=postgres password=1234 dbname=kafka sslmode=disable"

	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	// Загрузка кэша из БД
	if err := cache.LoadFromDB(dbConn); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}
	log.Println("Cache initialized")

	// Запуск Kafka consumer
	go kafka.StartConsumer(dbConn)

	// Настройка HTTP сервера
	r := mux.NewRouter()
	r.HandleFunc("/order/{id}", handlers.GetOrderHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates")))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("HTTP server listening on :8081")
	log.Fatal(srv.ListenAndServe())
}
