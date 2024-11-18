package grpc

import (
	inventory_pb "inventory-service/proto/inventory"
	"inventory-service/model"
	"context"
	"errors"
)

type Server struct {
	inventory_pb.UnimplementedInventoryServiceServer
    ProductInventory model.ProductInventory
}

func NewServer(productInventory model.ProductInventory) *Server {
	return &Server{
		ProductInventory: productInventory,
	}
}

func (s *Server) CheckStock(ctx context.Context, req *inventory_pb.StockRequest) (*inventory_pb.StockResponse, error) {
    quantity := s.ProductInventory.Inventory[req.ProductId]
    return &inventory_pb.StockResponse{
        ProductId: req.ProductId,
        Quantity:  quantity,
        InStock:   quantity > 0,
    }, nil
}

func (s *Server) UpdateStock(ctx context.Context, req *inventory_pb.UpdateStockRequest) (*inventory_pb.StockResponse, error) {
    if _, exists := s.ProductInventory.Inventory[req.ProductId]; !exists {
        return nil, errors.New("product not found")
    }
    
    s.ProductInventory.Inventory[req.ProductId] = req.Quantity
    return &inventory_pb.StockResponse{
        ProductId: req.ProductId,
        Quantity:  req.Quantity,
        InStock:   req.Quantity > 0,
    }, nil
}

func (s *Server) AddStock(ctx context.Context, req *inventory_pb.AddStockRequest) (*inventory_pb.StockResponse, error) {
    s.ProductInventory.Inventory[req.ProductId] = req.Quantity
    return &inventory_pb.StockResponse{
        ProductId: req.ProductId,
        Quantity:  req.Quantity,
        InStock:   req.Quantity > 0,
    }, nil
}

func (s *Server) DeleteStock(ctx context.Context, req *inventory_pb.StockRequest) (*inventory_pb.DeleteResponse, error) {
    if _, exists := s.ProductInventory.Inventory[req.ProductId]; !exists {
        return &inventory_pb.DeleteResponse{
            Success: false,
            Message: "product not found",
        }, nil
    }
    
    delete(s.ProductInventory.Inventory, req.ProductId)
    return &inventory_pb.DeleteResponse{
        Success: true,
        Message: "stock deleted successfully",
    }, nil
}