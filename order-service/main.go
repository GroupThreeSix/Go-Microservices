package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
	"order-service/config"
	"fmt"
)

type Order struct {
    ID         string   `json:"id"`
    ProductIDs []string `json:"product_ids"`
    Total      float64  `json:"total"`
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

    // Routes
    router.HandleFunc("/orders", GetOrders).Methods("GET")
    router.HandleFunc("/orders/{id}", GetOrder).Methods("GET")

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