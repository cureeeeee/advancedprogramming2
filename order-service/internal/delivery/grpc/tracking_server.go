package grpc

import (
	"errors"
	"strings"

	orderv1 "github.com/youruser/ap2-contracts-generated/gen/go/order/v1"
	"github.com/youruser/order-service/internal/pubsub"
	"github.com/youruser/order-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TrackingServer struct {
	orderv1.UnimplementedOrderTrackingServiceServer
	uc       *usecase.OrderUseCase
	notifier *pubsub.OrderNotifier
}

func NewTrackingServer(uc *usecase.OrderUseCase, notifier *pubsub.OrderNotifier) *TrackingServer {
	return &TrackingServer{uc: uc, notifier: notifier}
}

func (s *TrackingServer) SubscribeToOrderUpdates(req *orderv1.OrderRequest, stream orderv1.OrderTrackingService_SubscribeToOrderUpdatesServer) error {
	orderID := strings.TrimSpace(req.GetOrderId())
	if orderID == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := s.uc.GetOrder(stream.Context(), orderID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return status.Error(codes.NotFound, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	if err := stream.Send(&orderv1.OrderStatusUpdate{
		OrderId:   order.ID,
		Status:    order.Status,
		UpdatedAt: timestamppb.New(order.UpdatedAt),
	}); err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	updates := s.notifier.Subscribe(orderID)
	defer s.notifier.Unsubscribe(orderID, updates)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if err := stream.Send(&orderv1.OrderStatusUpdate{
				OrderId:   update.OrderID,
				Status:    update.Status,
				UpdatedAt: timestamppb.New(update.UpdatedAt),
			}); err != nil {
				return status.Error(codes.Unavailable, err.Error())
			}
		}
	}
}
