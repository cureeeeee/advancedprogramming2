package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	paymentv1 "github.com/cureeeeee/ap2-contracts-generated/gen/go/payment/v1"
	"github.com/cureeeeee/payment-service/internal/config"
	grpcdelivery "github.com/cureeeeee/payment-service/internal/delivery/grpc"
	"github.com/cureeeeee/payment-service/internal/messaging"
	"github.com/cureeeeee/payment-service/internal/repository/postgres"
	"github.com/cureeeeee/payment-service/internal/usecase"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	repo, err := postgres.NewPaymentRepository(db)
	if err != nil {
		log.Fatalf("initialize payment repository: %v", err)
	}

	publisher, err := messaging.NewRabbitMQPublisher(cfg.RabbitMQURL, cfg.PaymentQueue)
	if err != nil {
		log.Fatalf("create rabbitmq publisher: %v", err)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("close rabbitmq publisher: %v", err)
		}
	}()

	uc := usecase.NewPaymentUseCase(repo, publisher)
	server := grpcdelivery.NewPaymentServer(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcdelivery.LoggingInterceptor()))
	paymentv1.RegisterPaymentServiceServer(grpcServer, server)

	go func() {
		log.Printf("payment gRPC server is running on %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("grpc server stopped: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("payment service shutdown requested")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdownCtx.Err(); err != nil {
		log.Printf("shutdown context error: %v", err)
	}

	grpcServer.GracefulStop()
}
