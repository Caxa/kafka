package cache

import (
	"database/sql"
	"l0/cmd/internal/db"
	"l0/cmd/internal/models"
	"sync"
)

var (
	mu    sync.RWMutex
	store = make(map[string]models.Order)
)

func Get(id string) (models.Order, bool) {
	mu.RLock()
	defer mu.RUnlock()
	order, found := store[id]
	return order, found
}

func Set(order models.Order) {
	mu.Lock()
	defer mu.Unlock()
	store[order.OrderUID] = order
}

func LoadFromDB(conn *sql.DB) error {
	orders, err := db.GetAllOrders(conn)
	if err != nil {
		return err
	}
	for _, order := range orders {
		Set(order)
	}
	return nil
}
