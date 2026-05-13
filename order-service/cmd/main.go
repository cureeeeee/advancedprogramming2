package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderv1 "github.com/cureeeeee/ap2-contracts-generated/gen/go/order/v1"
	paymentclient "github.com/cureeeeee/order-service/internal/client/payment"
	"github.com/cureeeeee/order-service/internal/config"
	grpcdelivery "github.com/cureeeeee/order-service/internal/delivery/grpc"
	httpdelivery "github.com/cureeeeee/order-service/internal/delivery/http"
	"github.com/cureeeeee/order-service/internal/pubsub"
	"github.com/cureeeeee/order-service/internal/repository/memory"
	"github.com/cureeeeee/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	grpcServer := grpc.NewServer()
	orderv1.RegisterOrderTrackingServiceServer(grpcServer, grpcdelivery.NewTrackingServer(uc, notifier))

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen order gRPC: %v", err)
	}

	go func() {
		log.Printf("order gRPC server is running on %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("order gRPC server stopped: %v", err)
		}
	}()

	router := gin.Default()
	httpdelivery.NewHandler(uc).RegisterRoutes(router)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("order HTTP server is running on :%s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("run HTTP server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("order service shutdown requested")

	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown HTTP server: %v", err)
	}
}
