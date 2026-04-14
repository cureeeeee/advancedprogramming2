package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	orderv1 "github.com/youruser/ap2-contracts-generated/gen/go/order/v1"
	paymentclient "github.com/youruser/order-service/internal/client/payment"
	"github.com/youruser/order-service/internal/config"
	grpcdelivery "github.com/youruser/order-service/internal/delivery/grpc"
	httpdelivery "github.com/youruser/order-service/internal/delivery/http"
	"github.com/youruser/order-service/internal/pubsub"
	"github.com/youruser/order-service/internal/repository/memory"
	"github.com/youruser/order-service/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	repo := memory.NewOrderRepository()
	notifier := pubsub.NewOrderNotifier()

	paymentCli, err := paymentclient.NewClient(cfg.PaymentGRPCAddr)
	if err != nil {
		log.Fatalf("create payment gRPC client: %v", err)
	}
	defer func() {
		_ = paymentCli.Close()
	}()

	uc := usecase.NewOrderUseCase(repo, paymentCli, notifier)

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("listen order gRPC: %v", err)
		}

		grpcServer := grpc.NewServer()
		orderv1.RegisterOrderTrackingServiceServer(grpcServer, grpcdelivery.NewTrackingServer(uc, notifier))
		log.Printf("order gRPC server is running on %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve order gRPC: %v", err)
		}
	}()

	router := gin.Default()
	httpdelivery.NewHandler(uc).RegisterRoutes(router)
	log.Printf("order HTTP server is running on :%s", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("run HTTP server: %v", err)
	}
}
