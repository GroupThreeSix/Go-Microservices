package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product-service/config"
	"net"

	// "product-service/proto"
	"time"

	inventory_pb "product-service/proto/inventory"
	order_product_pb "product-service/proto/orderproduct"
	product_grpc "product-service/grpc"
	"product-service/model"
	// "product-service/proto/orderproduct"

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

var products []model.Product
var inventoryClient inventory_pb.InventoryServiceClient

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
	inventoryClient = inventory_pb.NewInventoryServiceClient(conn)

	router := mux.NewRouter()

	// Add health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Sample data
	products = append(products, model.Product{ID: "1", Name: "Laptop", Price: 999.99})
	products = append(products, model.Product{ID: "2", Name: "Mouse", Price: 29.99})

	// Start gRPC server
	grpcAddr := fmt.Sprintf("%s:%s", cfg.GrpcHost, cfg.GrpcPort)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ser := product_grpc.NewServer(inventoryClient ,products)
	grpcServer := grpc.NewServer()
	order_product_pb.RegisterOrderProductServiceServer(grpcServer, ser)
	go func() {
		log.Printf("Starting gRPC server on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// HTTP Endpoint
	router.HandleFunc("/products", GetProducts).Methods("GET")
    router.HandleFunc("/products/{id}", GetProduct).Methods("GET")
    router.HandleFunc("/products", CreateProduct).Methods("POST")
    router.HandleFunc("/products/{id}", UpdateProduct).Methods("PUT")
    router.HandleFunc("/products/{id}", DeleteProduct).Methods("DELETE")

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

		resp, err := inventoryClient.CheckStock(ctx, &inventory_pb.StockRequest{ProductId: product.ID})
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			resp, err := inventoryClient.CheckStock(ctx, &inventory_pb.StockRequest{ProductId: item.ID})
			if err != nil {
				log.Printf("Error checking stock for product %s: %v", params["id"], err)
			}

			enrichedProduct := Product{
				ID:       item.ID,
				Name:     item.Name,
				Price:    item.Price,
				InStock:  resp.InStock,
				Quantity: resp.Quantity,
			}

			json.NewEncoder(w).Encode(enrichedProduct)
			return
		}
	}
	
	json.NewEncoder(w).Encode(&Product{})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var product model.Product
    json.NewDecoder(r.Body).Decode(&product)

    // Add to inventory
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    _, err := inventoryClient.AddStock(ctx, &inventory_pb.AddStockRequest{
        ProductId: product.ID,
        Quantity: product.Quantity,
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    products = append(products, product)
    json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    var updatedProduct model.Product
    json.NewDecoder(r.Body).Decode(&updatedProduct)

    for i, item := range products {
        if item.ID == params["id"] {
            // Update inventory
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()
            
            _, err := inventoryClient.UpdateStock(ctx, &inventory_pb.UpdateStockRequest{
                ProductId: params["id"],
                Quantity: updatedProduct.Quantity,
            })
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            products[i] = updatedProduct
            json.NewEncoder(w).Encode(updatedProduct)
            return
        }
    }
    http.Error(w, "Product not found", http.StatusNotFound)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)

    for i, item := range products {
        if item.ID == params["id"] {
            // Delete from inventory
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()
            
            _, err := inventoryClient.DeleteStock(ctx, &inventory_pb.StockRequest{
                ProductId: params["id"],
            })
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            products = append(products[:i], products[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.Error(w, "Product not found", http.StatusNotFound)
}