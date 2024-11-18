package main

import (
	"context"
	"fmt"
	"inventory-service/config"
	"inventory-service/proto"
	"log"
	"net"

	"google.golang.org/grpc"
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

	log.Printf("Inventory service is running on port %s", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
