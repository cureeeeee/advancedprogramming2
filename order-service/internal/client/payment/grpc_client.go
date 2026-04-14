package payment

import (
	"context"
	"fmt"
	"time"

	paymentv1 "github.com/youruser/ap2-contracts-generated/gen/go/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn   *grpc.ClientConn
	client paymentv1.PaymentServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect payment service: %w", err)
	}

	return &Client{
		conn:   conn,
		client: paymentv1.NewPaymentServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ProcessPayment(ctx context.Context, orderID string, amount float64, currency string) (string, error) {
	resp, err := c.client.ProcessPayment(ctx, &paymentv1.PaymentRequest{
		OrderId:       orderID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: "card",
		RequestedAt:   timestamppb.New(time.Now().UTC()),
	})
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", fmt.Errorf("payment failed: %s", resp.GetMessage())
	}
	return resp.GetTransactionId(), nil
}
