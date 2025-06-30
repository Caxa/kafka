package handlers

import (
	"encoding/json"
	"l0/cmd/internal/cache"
	"l0/cmd/internal/db"
	"net/http"

	"github.com/gorilla/mux"
)

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if order, found := cache.Get(id); found {
		json.NewEncoder(w).Encode(order)
		return
	}

	order, err := db.GetOrderByID(db.Conn, id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	cache.Set(order)
	json.NewEncoder(w).Encode(order)
}
