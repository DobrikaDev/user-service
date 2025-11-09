package sql

import (
	"DobrikaDev/user-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (s *SqlStorage) GetBalance(ctx context.Context, maxID string) (*domain.Balance, error) {
	query, args := sq.Select(
		"b.id",
		"b.user_id",
		"b.balance",
	).
		From("balances b").
		Join("users u ON u.max_id = b.user_id").
		Where(sq.Eq{"u.max_id": maxID}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var balance domain.Balance
	err := s.trf.Transaction(ctx).GetContext(ctx, &balance, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBalanceNotFound
		}
		s.logger.Error("failed to get balance", zap.Error(err), zap.String("max_id", maxID))
		return nil, ErrBalanceInternal
	}

	return &balance, nil
}

func (s *SqlStorage) GetBalanceOperations(ctx context.Context, maxID string, limit int, offset int) ([]*domain.BalanceOperation, int32, error) {
	sb := sq.Select(
		"bo.id",
		"bo.balance_id",
		"bo.amount",
		"bo.type",
		"bo.description",
		"bo.created_at",
	).
		From("balance_operations bo").
		Join("balances b ON b.id = bo.balance_id").
		Where(sq.Eq{"b.user_id": maxID}).
		OrderBy("bo.created_at DESC").
		PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		sb = sb.Limit(uint64(limit))
	}
	if offset > 0 {
		sb = sb.Offset(uint64(offset))
	}

	query, args := sb.MustSql()

	operations := make([]*domain.BalanceOperation, 0, limit)
	err := s.trf.Transaction(ctx).SelectContext(ctx, &operations, query, args...)
	if err != nil {
		s.logger.Error("failed to get balance operations", zap.Error(err), zap.String("max_id", maxID))
		return nil, 0, ErrBalanceInternal
	}

	countQuery, countArgs := sq.Select("COUNT(*)").
		From("balance_operations bo").
		Join("balances b ON b.id = bo.balance_id").
		Where(sq.Eq{"b.user_id": maxID}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var total int32
	err = s.trf.Transaction(ctx).GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		s.logger.Error("failed to count balance operations", zap.Error(err), zap.String("max_id", maxID))
		return nil, 0, ErrBalanceInternal
	}

	return operations, total, nil
}

func (s *SqlStorage) CreateBalanceOperation(ctx context.Context, operation *domain.BalanceOperation) (*domain.BalanceOperation, error) {
	if operation == nil {
		return nil, ErrBalanceInvalid
	}

	err := s.TransactionManager.Do(ctx, func(txCtx context.Context) error {
		db := s.trf.Transaction(txCtx)

		var balance domain.Balance
		err := db.GetContext(txCtx, &balance, "SELECT id, user_id, balance FROM balances WHERE id = $1", operation.BalanceID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrBalanceNotFound
			}
			s.logger.Error("failed to get balance for operation", zap.Error(err), zap.String("balance_id", operation.BalanceID))
			return ErrBalanceInternal
		}

		switch operation.Type {
		case domain.BalanceOperationTypeWithdraw:
			if balance.Balance < operation.Amount {
				return ErrBalanceNotEnough
			}
			balance.Balance -= operation.Amount
		case domain.BalanceOperationTypeDeposit:
			balance.Balance += operation.Amount
		default:
			return ErrBalanceInvalid
		}

		opID := operation.ID
		if opID == "" {
			opID = uuid.NewString()
		}

		now := time.Now().UTC()
		_, err = db.ExecContext(txCtx,
			"INSERT INTO balance_operations (id, balance_id, amount, type, description, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
			opID,
			operation.BalanceID,
			operation.Amount,
			operation.Type,
			operation.Description,
			now,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrForeignKeyViolation {
				return ErrBalanceNotFound
			}
			s.logger.Error("failed to insert balance operation", zap.Error(err), zap.String("balance_id", operation.BalanceID))
			return ErrBalanceInternal
		}

		_, err = db.ExecContext(txCtx,
			"UPDATE balances SET balance = $1, updated_at = $2 WHERE id = $3",
			balance.Balance,
			now,
			balance.ID,
		)
		if err != nil {
			s.logger.Error("failed to update balance", zap.Error(err), zap.String("balance_id", balance.ID))
			return ErrBalanceInternal
		}

		operation.ID = opID
		operation.CreatedAt = now

		if operation.Type == domain.BalanceOperationTypeDeposit {
			var totalDeposits int
			err = db.GetContext(
				txCtx,
				&totalDeposits,
				"SELECT COALESCE(SUM(amount), 0) FROM balance_operations WHERE balance_id = $1 AND type = $2",
				operation.BalanceID,
				domain.BalanceOperationTypeDeposit,
			)
			if err != nil {
				s.logger.Error("failed to sum deposit operations", zap.Error(err), zap.String("balance_id", operation.BalanceID))
				return ErrBalanceInternal
			}

			var groupID int
			err = db.GetContext(
				txCtx,
				&groupID,
				`SELECT id
				 FROM reputation_groups
				 WHERE reputation_need <= $1
				 ORDER BY reputation_need DESC
				 LIMIT 1`,
				totalDeposits,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					groupID = 1
				} else {
					s.logger.Error("failed to select reputation group", zap.Error(err), zap.Int("total_deposits", totalDeposits))
					return ErrBalanceInternal
				}
			}

			if groupID > 0 {
				var currentGroupID int
				err = db.GetContext(
					txCtx,
					&currentGroupID,
					"SELECT reputation_group_id FROM users WHERE max_id = $1",
					balance.UserID,
				)
				if err != nil {
					s.logger.Error("failed to get current reputation group", zap.Error(err), zap.String("user_id", balance.UserID))
					return ErrBalanceInternal
				}

				if currentGroupID != groupID {
					_, err = db.ExecContext(
						txCtx,
						"UPDATE users SET reputation_group_id = $1, updated_at = $2 WHERE max_id = $3",
						groupID,
						now,
						balance.UserID,
					)
					if err != nil {
						s.logger.Error("failed to update reputation group", zap.Error(err), zap.String("user_id", balance.UserID))
						return ErrBalanceInternal
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return operation, nil
}
