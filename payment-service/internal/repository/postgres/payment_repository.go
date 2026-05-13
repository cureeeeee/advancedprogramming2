package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cureeeeee/payment-service/internal/domain"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) (*PaymentRepository, error) {
	repo := &PaymentRepository{db: db}
	if err := repo.ensureSchema(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PaymentRepository) ensureSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS payments (
		transaction_id TEXT PRIMARY KEY,
		order_id TEXT NOT NULL,
		amount DOUBLE PRECISION NOT NULL,
		currency TEXT NOT NULL,
		payment_method TEXT NOT NULL,
		requested_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *PaymentRepository) Save(transactionID string, payment domain.Payment) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec(
		`INSERT INTO payments (transaction_id, order_id, amount, currency, payment_method, requested_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		transactionID,
		payment.OrderID,
		payment.Amount,
		payment.Currency,
		payment.PaymentMethod,
		payment.RequestedAt,
		payment.RequestedAt,
	)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
