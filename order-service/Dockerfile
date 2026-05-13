FROM golang:1.26-alpine

WORKDIR /app/order-service
RUN apk add --no-cache git

COPY order-service/go.mod order-service/go.sum ./
COPY contracts-generated ../contracts-generated
COPY order-service .

RUN go mod download
RUN go build -o /app/bin/order-service ./cmd

CMD ["/app/bin/order-service"]
