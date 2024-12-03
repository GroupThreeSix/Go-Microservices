package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-service/client"
	"order-service/config"
	"order-service/model"
	order_product_pb "order-service/proto/orderproduct"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var orders []model.Order
var productClient *client.ProductClient

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// Setup gRPC connection to product service
	productAddr := fmt.Sprintf("%s:%s",
		cfg.ProductServiceHost,
		cfg.ProductServicePort,
	)
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()

	productClient = client.NewProductClient(productConn)

	// Initialize router
	router := mux.NewRouter()

	// Sample data
	orders = append(orders, model.Order{ID: "1", ProductIDs: []string{"1", "2"}, Total: 1029.98})

	// Add health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")
	// Routes
	router.HandleFunc("/orders", GetOrders).Methods("GET")
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
	json.NewEncoder(w).Encode(&model.Order{})
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "order-service",
	})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Validate products through gRPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := productClient.ValidateProducts(ctx, order.ProductIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !resp.Valid {
		http.Error(w, resp.Error, http.StatusBadRequest)
		return
	}

	// Calculate total from validated products
	var total float64
	for _, product := range resp.Products {
		total += product.Price
	}
	order.Total = total
	order.Status = "pending"

	// Update product stock
	var orderItems []*order_product_pb.OrderItem
	for _, productID := range order.ProductIDs {
		orderItems = append(orderItems, &order_product_pb.OrderItem{
			ProductId: productID,
			Quantity: 1,
		})
	}

	if err := productClient.UpdateStock(ctx, orderItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedOrder model.Order
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
