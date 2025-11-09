package sqlxtrm

import (
	"context"
	"database/sql"
	"time"

	"DobrikaDev/user-service/internal/storage/deps"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	trmSettings "github.com/avito-tech/go-transaction-manager/trm/settings"
	"github.com/jmoiron/sqlx"
)

type SqlxTransactionManager struct {
	manager trm.Manager
}

func NewSqlxTransactionManager(db *sqlx.DB) (*SqlxTransactionManager, error) {
	trmManager, err := manager.New(trmsqlx.NewDefaultFactory(db))
	if err != nil {
		return nil, err
	}

	return &SqlxTransactionManager{
		manager: trmManager,
	}, nil
}

func (stm *SqlxTransactionManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return stm.manager.Do(context.WithoutCancel(ctx), fn)
}

func (stm *SqlxTransactionManager) DoWithCancel(ctx context.Context, fn func(context.Context) error) error {
	return stm.manager.Do(ctx, fn)
}

func (stm *SqlxTransactionManager) DoWithTimeout(
	ctx context.Context,
	timeout time.Duration,
	fn func(context.Context) error,
) error {
	return stm.manager.DoWithSettings(
		ctx,
		trmSettings.Must(trmSettings.WithTimeout(timeout)),
		fn,
	)
}

type SqlxTransactionFactory struct {
	*sqlx.DB
}

func NewSqlxTransactionFactory(db *sqlx.DB) *SqlxTransactionFactory {
	return &SqlxTransactionFactory{db}
}

func (t *SqlxTransactionFactory) Transaction(ctx context.Context) deps.Transaction {
	return trmsqlx.DefaultCtxGetter.DefaultTrOrDB(context.WithoutCancel(ctx), t)
}

func (t *SqlxTransactionFactory) TransactionWithCancel(ctx context.Context) deps.Transaction {
	return trmsqlx.DefaultCtxGetter.DefaultTrOrDB(ctx, t)
}

func (t *SqlxTransactionFactory) GetDB() *sql.DB {
	return t.DB.DB
}
