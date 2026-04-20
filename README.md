
# Advanced Programming 2 - Assignment 2

This repository contains the gRPC migration of the Order/Payment microservices with a contract-first workflow.

## Repositories

- Proto repository: https://github.com/cureeeeee/advancedprogramming2/tree/main/contracts-proto
- Generated contracts repository: https://github.com/cureeeeee/advancedprogramming2/tree/main/contracts-generated
- Services repository: https://github.com/cureeeeee/advancedprogramming2

If your real repository names are different, update the links above before submission.

## Architecture

- External clients call Order Service via REST (Gin).
- Order Service calls Payment Service via gRPC (unary RPC).
- Order Service also exposes gRPC server-side streaming for order tracking.
- Contracts are defined in the proto repository and generated into a separate shared repository.

See the diagram: [ARCHITECTURE.md](ARCHITECTURE.md)

## Project Structure

- contracts-proto: protobuf contracts, Buf config, and GitHub Actions for remote generation.
- contracts-generated: generated Go contracts used by services.
- order-service: REST API + gRPC tracking server + gRPC payment client.
- payment-service: gRPC payment server + interceptor logging.

## Environment Variables

Order Service:

- ORDER_HTTP_PORT (example: 8080)
- ORDER_GRPC_PORT (example: 50052)
- PAYMENT_GRPC_ADDRESS (example: localhost:50051)

Payment Service:

- PAYMENT_GRPC_PORT (example: 50051)

Use provided examples:

- order-service/.env.example
- payment-service/.env.example

## Run Instructions

1. Start Payment Service:

```bash
cd payment-service
go mod tidy
go run ./cmd
```

2. Start Order Service:

```bash
cd order-service
go mod tidy
go run ./cmd
```

3. Create order via REST:

```bash
curl -X POST http://localhost:8080/orders \
	-H "Content-Type: application/json" \
	-d '{"amount":100,"currency":"KZT"}'
```

4. Subscribe to order updates via gRPC stream:

```bash
cd order-service
go run ./cmd/stream-client -addr localhost:50052 -order <ORDER_ID>
```

5. Trigger a real status update (updates repository state and pushes into stream):

```bash
curl -X PUT http://localhost:8080/orders/<ORDER_ID>/status \
	-H "Content-Type: application/json" \
	-d '{"status":"SHIPPED"}'
```

## Contract-First Workflow

Proto contracts are stored in contracts-proto and generated Go code is synced into contracts-generated.

Local check:

```bash
cd contracts-proto
buf lint
buf generate
```

Remote generation is configured in:

- contracts-proto/.github/workflows/remote-generate.yml

## Bonus: gRPC Interceptor

Payment Service includes a unary interceptor that logs method name and duration for each incoming request.

## Evidence Checklist

Capture screenshots for:

1. Successful payment gRPC call from Order Service to Payment Service.
2. Stream client receiving immediate update after status change.
3. Payment interceptor logs in server console.
4. GitHub Actions run that syncs generated contracts.

Add these under an `evidence/` folder and include them in the ZIP submission.

