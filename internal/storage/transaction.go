package storage

import (
	"context"
	"fmt"
	"time"

	tr "github.com/xRepZ/payments/internal"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) AddTransaction(ctx context.Context, userId int, email string, amount float64, currency string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := psql.Insert("transactions").
		Columns("user_id", "email", "amount", "currency", "status").Values(userId, email, amount, currency, "new").
		ToSql()
	if err != nil {
		return fmt.Errorf("can't create query: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = t.db.QueryRowxContext(dbCtx, sql, args...).Err()
	if err != nil {
		return fmt.Errorf("can't do db query: %w", err)
	}
	fmt.Println("ok")
	return nil
}

func (t *Transaction) UpdateById(ctx context.Context, id int, status string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Update("transactions").Where("id = ?", id).Set("status", status).ToSql()

	if err != nil {
		return fmt.Errorf("can't create query: %w", err)
	}
	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	err = t.db.QueryRowxContext(dbCtx, sql, args...).Err()
	if err != nil {
		return fmt.Errorf("can't do db query: %w", err)
	}

	return nil
}

func (t *Transaction) CancelById(ctx context.Context, id int) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Delete("*").
		From("transactions").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't create query: %w", err)
	}
	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = t.db.QueryRowxContext(dbCtx, sql, args...).Err()
	if err != nil {
		return fmt.Errorf("can't do db query: %w", err)
	}
	return nil

}

func (t *Transaction) GetByMail(ctx context.Context, email string) ([]*tr.Transactions, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Select("*").
		From("transactions").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("can't create query: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	rows, err := t.db.QueryxContext(dbCtx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Error then trying to get rows: %w", err)
	}

	defer rows.Close()
	allUserT := make([]*tr.Transactions, 0)
	for rows.Next() {
		row := &tr.Transactions{}
		err = rows.StructScan(&row)
		if err != nil {
			return nil, fmt.Errorf("Error when trying to wtire rows: %w", err)
		}
		allUserT = append(allUserT, row)
	}

	return allUserT, nil
}

func (t *Transaction) GetByUserId(ctx context.Context, userId int) ([]*tr.Transactions, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Select("*").
		From("transactions").
		Where("user_id = ?", userId).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("can't create query: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := t.db.QueryxContext(dbCtx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Error then trying to get rows: %w", err)
	}

	defer rows.Close()
	allUserT := make([]*tr.Transactions, 0)
	for rows.Next() {
		row := &tr.Transactions{}
		err = rows.StructScan(&row)
		if err != nil {
			return nil, fmt.Errorf("Error when trying to wtire rows: %w", err)
		}
		allUserT = append(allUserT, row)
	}

	return allUserT, nil

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
