package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-service/config"

	"github.com/gorilla/mux"
)

type Order struct {
	ID         string   `json:"id"`
	ProductIDs []string `json:"product_ids"`
	Total      float64  `json:"total"`
	Status     string   `json:"status"`
}

var orders []Order

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// Initialize router
	router := mux.NewRouter()

	// Sample data
	orders = append(orders, Order{ID: "1", ProductIDs: []string{"1", "2"}, Total: 1029.98})

	// Add health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")
	// Routes
    router.HandleFunc("/orders/{id}", GetOrder).Methods("GET")
    router.HandleFunc("/orders", CreateOrder).Methods("POST")
    router.HandleFunc("/orders/{id}", UpdateOrder).Methods("PUT")
    router.HandleFunc("/orders/{id}", DeleteOrder).Methods("DELETE")

	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Order service is running on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range orders {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Order{})
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "order-service",
	})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	order.Status = "pending"
	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedOrder Order
	json.NewDecoder(r.Body).Decode(&updatedOrder)

	for i, item := range orders {
		if item.ID == params["id"] {
			orders[i] = updatedOrder
			json.NewEncoder(w).Encode(updatedOrder)
			return
		}
	}
	http.Error(w, "Order not found", http.StatusNotFound)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, item := range orders {
		if item.ID == params["id"] {
			orders = append(orders[:i], orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Order not found", http.StatusNotFound)
}
