package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"l0/cmd/internal/cache"
	"l0/cmd/internal/db"
	"l0/cmd/internal/handlers"
	"l0/cmd/internal/kafka"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Order Service v1.0...")

	// Инициализация подключения к PostgreSQL
	dsn := getEnv("DB_DSN", "host=postgres port=5432 user=order_user password=password dbname=ordersdb sslmode=disable")
	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer dbConn.Close()

	// Загрузка кэша из БД
	if err := cache.LoadFromDB(dbConn); err != nil {
		log.Fatalf("Cache initialization error: %v", err)
	}
	log.Println("✅ Cache loaded successfully")

	// Инициализация Kafka Consumer
	consumerCfg := kafka.ConsumerConfig{
		DB:     dbConn,
		Broker: getEnv("KAFKA_BROKER", "kafka:29092"),
		Topic:  getEnv("KAFKA_TOPIC", "orders_topic"),
	}
	go kafka.StartConsumer(consumerCfg)
	log.Printf("✅ Kafka consumer started (broker: %s, topic: %s)", consumerCfg.Broker, consumerCfg.Topic)

	// Настройка HTTP сервера
	r := mux.NewRouter()
	r.HandleFunc("/order/{id}", handlers.GetOrderHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 HTTP server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("🛑 Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("✅ Server exited properly")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
