package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"l0/cmd/internal/cache"
	"l0/cmd/internal/db"
	"l0/cmd/internal/models"

	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	DB     *sql.DB
	Broker string
	Topic  string
}

func StartConsumer(config ConsumerConfig) {
	// Настройка Kafka Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{config.Broker},
		Topic:       config.Topic,
		GroupID:     "order-group",
		StartOffset: kafka.LastOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
	})
	defer r.Close()

	log.Printf("Kafka consumer started on topic: %s (broker: %s)", config.Topic, config.Broker)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			time.Sleep(2 * time.Second) // Задержка при ошибках
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		if err := db.InsertOrder(config.DB, order); err != nil {
			log.Printf("Failed to save order: %v", err)
			continue
		}

		cache.Set(order)
		log.Printf("Order %s saved and cached", order.OrderUID)
	}
}
