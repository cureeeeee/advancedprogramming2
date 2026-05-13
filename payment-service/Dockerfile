FROM golang:1.26-alpine

WORKDIR /app/payment-service
RUN apk add --no-cache git

COPY payment-service/go.mod payment-service/go.sum ./
COPY contracts-generated ../contracts-generated
COPY payment-service .

RUN go mod download
RUN go build -o /app/bin/payment-service ./cmd

CMD ["/app/bin/payment-service"]
