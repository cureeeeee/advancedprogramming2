package domain

import "time"

type Order struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StatusUpdate struct {
	OrderID   string
	Status    string
	UpdatedAt time.Time
}
