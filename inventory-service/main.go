package main

import (
	"context"
	"fmt"
	"inventory-service/config"
	"inventory-service/proto"
	"log"
	"net"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type server struct {
	proto.UnimplementedInventoryServiceServer
}

// In-memory inventory data
var inventory = map[string]int32{
	"1": 100, // Product ID 1 has 100 items
	"2": 50,  // Product ID 2 has 50 items
}

func (s *server) CheckStock(ctx context.Context, req *proto.StockRequest) (*proto.StockResponse, error) {
    quantity := inventory[req.ProductId]
    return &proto.StockResponse{
        ProductId: req.ProductId,
        Quantity:  quantity,
        InStock:   quantity > 0,
    }, nil
}

func (s *server) UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.StockResponse, error) {
    if _, exists := inventory[req.ProductId]; !exists {
        return nil, errors.New("product not found")
    }
    
    inventory[req.ProductId] = req.Quantity
    return &proto.StockResponse{
        ProductId: req.ProductId,
        Quantity:  req.Quantity,
        InStock:   req.Quantity > 0,
    }, nil
}

func (s *server) AddStock(ctx context.Context, req *proto.AddStockRequest) (*proto.StockResponse, error) {
    inventory[req.ProductId] = req.Quantity
    return &proto.StockResponse{
        ProductId: req.ProductId,
        Quantity:  req.Quantity,
        InStock:   req.Quantity > 0,
    }, nil
}

func (s *server) DeleteStock(ctx context.Context, req *proto.StockRequest) (*proto.DeleteResponse, error) {
    if _, exists := inventory[req.ProductId]; !exists {
        return &proto.DeleteResponse{
            Success: false,
            Message: "product not found",
        }, nil
    }
    
    delete(inventory, req.ProductId)
    return &proto.DeleteResponse{
        Success: true,
        Message: "stock deleted successfully",
    }, nil
}

func main() {
	//Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	//Set up gRPC
	grpcAddr := fmt.Sprintf("%s:%s", cfg.GrpcHost, cfg.GrpcPort)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterInventoryServiceServer(s, &server{})

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	log.Printf("Inventory service is running on port %s", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
