package grpc

import (
	"context"
	"strings"
	"time"

	paymentv1 "github.com/youruser/ap2-contracts-generated/gen/go/payment/v1"
	"github.com/youruser/payment-service/internal/domain"
	"github.com/youruser/payment-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentServer(uc *usecase.PaymentUseCase) *PaymentServer {
	return &PaymentServer{uc: uc}
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	if strings.TrimSpace(req.GetOrderId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.GetAmount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than 0")
	}
	if strings.TrimSpace(req.GetCurrency()) == "" {
		return nil, status.Error(codes.InvalidArgument, "currency is required")
	}

	requestedAt := time.Now().UTC()
	if req.GetRequestedAt() != nil {
		requestedAt = req.GetRequestedAt().AsTime()
	}

	result, err := s.uc.ProcessPayment(ctx, domain.Payment{
		OrderID:       req.GetOrderId(),
		Amount:        req.GetAmount(),
		Currency:      req.GetCurrency(),
		PaymentMethod: req.GetPaymentMethod(),
		RequestedAt:   requestedAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "process payment: %v", err)
	}

	return &paymentv1.PaymentResponse{
		Success:       result.Success,
		TransactionId: result.TransactionID,
		Message:       result.Message,
		ProcessedAt:   timestamppb.New(result.ProcessedAt),
	}, nil
}
