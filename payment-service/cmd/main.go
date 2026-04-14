package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	paymentv1 "github.com/youruser/ap2-contracts-generated/gen/go/payment/v1"
	"github.com/youruser/payment-service/internal/config"
	grpcdelivery "github.com/youruser/payment-service/internal/delivery/grpc"
	"github.com/youruser/payment-service/internal/repository/memory"
	"github.com/youruser/payment-service/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	repo := memory.NewPaymentRepository()
	uc := usecase.NewPaymentUseCase(repo)
	server := grpcdelivery.NewPaymentServer(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcdelivery.LoggingInterceptor()))
	paymentv1.RegisterPaymentServiceServer(grpcServer, server)

	log.Printf("payment gRPC server is running on %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("grpc server stopped: %v", err)
		os.Exit(1)
	}
}
