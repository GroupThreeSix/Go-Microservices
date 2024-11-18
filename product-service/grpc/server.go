package grpc

import (
    // "context"
    order_product_pb "product-service/proto/orderproduct"
	"product-service/model"
	"context"
	inventory_product_pb "product-service/proto/inventory"
)

type Server struct {
    order_product_pb.UnimplementedOrderProductServiceServer
	inventoryClient  inventory_product_pb.InventoryServiceClient
    products []model.Product
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	InStock  bool    `json:"in_stock"`
	Quantity int32   `json:"quantity"`
}


func NewServer(inventoryClient inventory_product_pb.InventoryServiceClient ,products []model.Product) *Server {
    return &Server{
		inventoryClient: inventoryClient,
        products: products,
    }
}

func (s *Server) ValidateProducts(ctx context.Context, req *order_product_pb.ValidateProductsRequest) (*order_product_pb.ValidateProductsResponse, error) {
    var validProducts []*order_product_pb.ProductInfo
    
    for _, id := range req.ProductIds {
        for _, product := range s.products {
            if product.ID == id {
                validProducts = append(validProducts, &order_product_pb.ProductInfo{
                    Id:       product.ID,
                    Name:     product.Name,
                    Price:    product.Price,
                    InStock:  product.InStock,
                    Quantity: product.Quantity,
                })
                break
            }
        }
    }
    
    if len(validProducts) != len(req.ProductIds) {
        return &order_product_pb.ValidateProductsResponse{
            Valid: false,
            Error: "some products not found",
        }, nil
    }
    
    return &order_product_pb.ValidateProductsResponse{
        Valid:    true,
        Products: validProducts,
    }, nil
}

func (s *Server) UpdateProductStock(ctx context.Context, req *order_product_pb.UpdateStockRequest) (*order_product_pb.UpdateStockResponse, error) {
    // Update inventory through gRPC client
    for _, item := range req.Items {
        // Call inventory service to update stock
        _, err := s.inventoryClient.UpdateStock(ctx, &inventory_product_pb.UpdateStockRequest{
            ProductId: item.ProductId,
            Quantity: item.Quantity,
        })
        if err != nil {
            return &order_product_pb.UpdateStockResponse{
                Success: false,
                Error:   err.Error(),
            }, nil
        }
    }
    
    return &order_product_pb.UpdateStockResponse{
        Success: true,
    }, nil
}