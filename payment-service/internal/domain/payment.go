package domain

import "time"

type Payment struct {
	OrderID       string
	Amount        float64
	Currency      string
	PaymentMethod string
	RequestedAt   time.Time
}

type PaymentResult struct {
	Success       bool
	TransactionID string
	Message       string
	ProcessedAt   time.Time
}
