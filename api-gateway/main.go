package main

import (
	"api-gateway/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

var cfg config.Config
var err error

func main() {
	// Load configuration
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// Initialize router
	router := mux.NewRouter()

	// Add health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Routes
	router.PathPrefix("/products").HandlerFunc(handleProduct)
	router.PathPrefix("/orders").HandlerFunc(handleOrder)

	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("API Gateway is running on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "api-gateway",
	})
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	// Using the service name defined in docker-compose
	productServiceURL, _ := url.Parse(cfg.ProductServiceURL)
	proxy := httputil.NewSingleHostReverseProxy(productServiceURL)
	proxy.ServeHTTP(w, r)
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	// Using the service name defined in docker-compose
	orderServiceURL, _ := url.Parse(cfg.OrderServiceURL)
	proxy := httputil.NewSingleHostReverseProxy(orderServiceURL)
	proxy.ServeHTTP(w, r)
}
