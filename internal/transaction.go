package internal

import (
	"context"
	"database/sql"
	"time"
)

type Transactions struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"userId" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Amount    float64   `json:"amount" db:"amount"`
	Сurrency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdateAt  time.Time `json:"updated_at" db:"updated_at"`
	Status    string    `json:"status" db:"status"`

	DeletedAt sql.NullTime `db:"deleted_at"`
}

type TransactionStorage interface {
	// Add(ctx context.Context, userId int, email string, amount float64, currency string) error
	// Cancel(ctx context.Context, id int) error
	// Update(ctx context.Context, id int, status string) error
	GetById(ctx context.Context, id int) (string, error)
	// GetByUserId(ctx context.Context, userId int) ([]*Transactions, error)
	// GetByUserEmail(ctx context.Context, email string) ([]*Transactions, error)
}
