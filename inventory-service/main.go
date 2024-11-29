package main

import (
	"context"
	"errors"
	"fmt"
	"inventory-service/config"
	"inventory-service/model"

	// "inventory-service/proto"
	inventory_pb "inventory-service/proto/inventory"
	"log"
	"net"
    inventory_grpc "inventory-service/grpc"

	"google.golang.org/grpc"
	// "google.golang.org/grpc/health"
	// "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	inventory_pb.UnimplementedInventoryServiceServer
    productInfo model.ProductInventory
}

// func NewServer(inventory model.Inventory) *server {
    
// }

// In-memory inventory data
// var inventory = map[string]int32{
// 	"1": 100, // Product ID 1 has 100 items
// 	"2": 50,  // Product ID 2 has 50 items
// }

func (s *Server) CheckStock(ctx context.Context, req *inventory_pb.StockRequest) (*inventory_pb.StockResponse, error) {
	quantity := s.productInfo.Inventory[req.ProductId]
	return &inventory_pb.StockResponse{
		ProductId: req.ProductId,
		Quantity:  quantity,
		InStock:   quantity > 0,
	}, nil
}

func (s *Server) UpdateStock(ctx context.Context, req *inventory_pb.UpdateStockRequest) (*inventory_pb.StockResponse, error) {
	if _, exists := s.productInfo.Inventory[req.ProductId]; !exists {
		return nil, errors.New("product not found")
	}

	s.productInfo.Inventory[req.ProductId] = req.Quantity
	return &inventory_pb.StockResponse{
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
		InStock:   req.Quantity > 0,
	}, nil
}

func (s *Server) AddStock(ctx context.Context, req *inventory_pb.AddStockRequest) (*inventory_pb.StockResponse, error) {
	s.productInfo.Inventory[req.ProductId] = req.Quantity
	return &inventory_pb.StockResponse{
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
		InStock:   req.Quantity > 0,
	}, nil
}

var productInfo model.ProductInventory

func (s *Server) DeleteStock(ctx context.Context, req *inventory_pb.StockRequest) (*inventory_pb.DeleteResponse, error) {
	if _, exists := s.productInfo.Inventory[req.ProductId]; !exists {
		return &inventory_pb.DeleteResponse{
			Success: false,
			Message: "product not found",
		}, nil
	}


// func (s *server) CheckStock(ctx context.Context, req *inventorypb.StockRequest) (*inventorypb.StockResponse, error) {
//     quantity := inventory[req.ProductId]
//     return &inventorypb.StockResponse{
//         ProductId: req.ProductId,
//         Quantity:  quantity,
//         InStock:   quantity > 0,
//     }, nil
// }

// func (s *server) UpdateStock(ctx context.Context, req *inventorypb.UpdateStockRequest) (*inventorypb.StockResponse, error) {
//     if _, exists := inventory[req.ProductId]; !exists {
//         return nil, errors.New("product not found")
//     }
    
//     inventory[req.ProductId] = req.Quantity
//     return &inventorypb.StockResponse{
//         ProductId: req.ProductId,
//         Quantity:  req.Quantity,
//         InStock:   req.Quantity > 0,
//     }, nil
// }

// func (s *server) AddStock(ctx context.Context, req *inventorypb.AddStockRequest) (*inventorypb.StockResponse, error) {
//     inventory[req.ProductId] = req.Quantity
//     return &inventorypb.StockResponse{
//         ProductId: req.ProductId,
//         Quantity:  req.Quantity,
//         InStock:   req.Quantity > 0,
//     }, nil
// }

// func (s *server) DeleteStock(ctx context.Context, req *inventorypb.StockRequest) (*inventorypb.DeleteResponse, error) {
//     if _, exists := inventory[req.ProductId]; !exists {
//         return &inventorypb.DeleteResponse{
//             Success: false,
//             Message: "product not found",
//         }, nil
//     }
    
//     delete(inventory, req.ProductId)
//     return &inventorypb.DeleteResponse{
//         Success: true,
//         Message: "stock deleted successfully",
//     }, nil
// }

	delete(s.productInfo.Inventory, req.ProductId)
	return &inventory_pb.DeleteResponse{
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

	productInfo = model.ProductInventory{
		Inventory: map[string]int32{"1": 100, "2": 50},
	}

	//Set up gRPC
	grpcAddr := fmt.Sprintf("%s:%s", cfg.GrpcHost, cfg.GrpcPort)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

    server := inventory_grpc.NewServer(productInfo)
	grpcServer := grpc.NewServer()
	inventory_pb.RegisterInventoryServiceServer(grpcServer, server)

	// Register health service
	// healthServer := health.NewServer()
	// grpc_health_v1.RegisterHealthServer(s, healthServer)

	log.Printf("Inventory service is running on port %s", cfg.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
