package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) GetById(ctx context.Context, id int) (string, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Select("status").
		From("transactions").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("can't create query: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	status := ""
	err = t.db.QueryRowxContext(dbCtx, sql, args...).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("can't do db query: %w", err)
	}

	return status, nil
}

// func (t *Transaction) GetByUserEmail(ctx context.Context, email string) ([]*internal.Transactions, error)
