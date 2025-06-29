package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"l0/cmd/internal/cache"
	"l0/cmd/internal/db"
	"l0/cmd/internal/models"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(dbConn *sql.DB) {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "orders_topic"
		log.Println("KAFKA_TOPIC not set, defaulting to 'orders_topic'")
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		log.Println("KAFKA_BROKER not set, defaulting to 'localhost:9092'")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       topic,
		GroupID:     "order-group",
		StartOffset: kafka.LastOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
	})
	defer r.Close()

	log.Println("Kafka consumer started on topic:", topic)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message from Kafka:", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		if err := db.InsertOrder(dbConn, order); err != nil {
			log.Println("Failed to save order:", err)
			continue
		}

		cache.Set(order)
		log.Printf("Order %s saved and cached", order.OrderUID)
	}
}
