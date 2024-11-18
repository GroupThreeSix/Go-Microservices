package client

import (
	"context"
	order_product_pb "order-service/proto/orderproduct"

	"google.golang.org/grpc"
)

type ProductClient struct {
	client order_product_pb.OrderProductServiceClient
}

func NewProductClient(conn *grpc.ClientConn) *ProductClient {
	return &ProductClient{
		client: order_product_pb.NewOrderProductServiceClient(conn),
	}
}

func (c *ProductClient) ValidateProducts(ctx context.Context, productIDs []string) (*order_product_pb.ValidateProductsResponse, error) {
	return c.client.ValidateProducts(ctx, &order_product_pb.ValidateProductsRequest{
		ProductIds: productIDs,
	})
}

func (c *ProductClient) UpdateStock(ctx context.Context, items []*order_product_pb.OrderItem) error {
	_, err := c.client.UpdateProductStock(ctx, &order_product_pb.UpdateStockRequest{
		Items: items,
	})
	return err
}