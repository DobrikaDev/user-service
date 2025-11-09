package deps

import (
	"context"
	"database/sql"
	"time"

	"github.com/avito-tech/go-transaction-manager/sqlx"
)

type Transaction interface {
	sqlx.Tr
}

type TransactionFactory interface {
	Transaction(ctx context.Context) Transaction
	TransactionWithCancel(ctx context.Context) Transaction
	GetDB() *sql.DB
}

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
	DoWithCancel(ctx context.Context, fn func(ctx context.Context) error) error
	DoWithTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error
}

type TrmStub struct{}

func NewTrmStub() *TrmStub {
	return &TrmStub{}
}

func (*TrmStub) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (*TrmStub) DoWithCancel(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (*TrmStub) DoWithTimeout(
	ctx context.Context,
	_ time.Duration,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}
