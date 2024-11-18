package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product-service/config"
	"product-service/proto"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	InStock  bool    `json:"in_stock"`
	Quantity int32   `json:"quantity"`
}

var products []Product
var inventoryClient proto.InventoryServiceClient

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Set up gRPC connection to inventory service
	inventoryAddr := fmt.Sprintf("%s:%s",
		cfg.InventoryServiceHost,
		cfg.InventoryServicePort,
	)
	conn, err := grpc.NewClient(inventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	inventoryClient = proto.NewInventoryServiceClient(conn)

	router := mux.NewRouter()

	// Add health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Sample data
	products = append(products, Product{ID: "1", Name: "Laptop", Price: 999.99})
	products = append(products, Product{ID: "2", Name: "Mouse", Price: 29.99})

	router.HandleFunc("/products", GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", GetProduct).Methods("GET")

	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Product service is running on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "product-service",
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get inventory information for each product
	enrichedProducts := make([]Product, len(products))
	for i, product := range products {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := inventoryClient.CheckStock(ctx, &proto.StockRequest{ProductId: product.ID})
		if err != nil {
			log.Printf("Error checking stock for product %s: %v", product.ID, err)
			continue
		}

		enrichedProducts[i] = Product{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			InStock:  resp.InStock,
			Quantity: resp.Quantity,
		}
	}

	json.NewEncoder(w).Encode(enrichedProducts)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range products {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Product{})
}
