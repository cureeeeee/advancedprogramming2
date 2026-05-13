FROM golang:1.26-alpine

WORKDIR /app/notification-service
RUN apk add --no-cache git

COPY notification-service/go.mod notification-service/go.sum ./
COPY notification-service .

RUN go mod download
RUN go build -o /app/bin/notification-service ./cmd

CMD ["/app/bin/notification-service"]
