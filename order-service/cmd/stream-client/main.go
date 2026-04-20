package main

import (
	"context"
	"flag"
	"io"
	"log"

	orderv1 "github.com/cureeeeee/ap2-contracts-generated/gen/go/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("addr", "localhost:50052", "order tracking gRPC address")
	orderID := flag.String("order", "", "order id")
	flag.Parse()

	if *orderID == "" {
		log.Fatal("-order is required")
	}

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial order tracking server: %v", err)
	}
	defer conn.Close()

	client := orderv1.NewOrderTrackingServiceClient(conn)
	stream, err := client.SubscribeToOrderUpdates(context.Background(), &orderv1.SubscribeToOrderUpdatesRequest{OrderId: *orderID})
	if err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			log.Println("stream closed")
			return
		}
		if err != nil {
			log.Fatalf("stream recv: %v", err)
		}
		log.Printf("order=%s status=%s updated_at=%s", update.GetOrderId(), update.GetStatus(), update.GetUpdatedAt().AsTime())
	}
}
