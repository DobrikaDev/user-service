package balance

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/internal/storage/sql"
	"context"

	"errors"

	"go.uber.org/zap"
)

func (s *BalanceService) GetBalance(ctx context.Context, maxID string) (*domain.Balance, error) {
	balance, err := s.storage.GetBalance(ctx, maxID)
	if err != nil {
		s.logger.Error("failed to get balance", zap.Error(err), zap.String("max_id", maxID))
		if errors.Is(err, sql.ErrBalanceNotFound) {
			return nil, ErrBalanceNotFound
		}
		return nil, err
	}

	return &domain.Balance{
		ID:      balance.ID,
		UserID:  balance.UserID,
		Balance: balance.Balance,
	}, nil
}

func (s *BalanceService) GetBalanceOperations(ctx context.Context, maxID string, limit int, offset int) ([]*domain.BalanceOperation, int32, error) {
	operations, total, err := s.storage.GetBalanceOperations(ctx, maxID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get balance operations", zap.Error(err), zap.String("max_id", maxID), zap.Int("limit", limit), zap.Int("offset", offset))
		return nil, 0, err
	}
	return operations, total, nil
}

func (s *BalanceService) CreateOperation(ctx context.Context, maxID string, operation *domain.BalanceOperation) (*domain.BalanceOperation, error) {
	if operation.Amount <= 0 {
		return nil, ErrBalanceInvalid
	}
	if operation.Type == "" {
		return nil, ErrBalanceInvalid
	}
	balance, err := s.storage.GetBalance(ctx, maxID)
	if err != nil {
		if errors.Is(err, sql.ErrBalanceNotFound) {
			return nil, ErrBalanceNotFound
		}
		s.logger.Error("failed to get balance by user id", zap.Error(err), zap.String("max_id", maxID))
		return nil, ErrBalanceInternal
	}
	operation.BalanceID = balance.ID
	created, err := s.storage.CreateBalanceOperation(ctx, operation)
	if err != nil {
		s.logger.Error("failed to create operation", zap.Error(err), zap.Any("operation", operation))
		if errors.Is(err, sql.ErrBalanceNotFound) {
			return nil, ErrBalanceNotFound
		}
		if errors.Is(err, sql.ErrBalanceInvalid) {
			return nil, ErrBalanceInvalid
		}
		if errors.Is(err, sql.ErrBalanceNotEnough) {
			return nil, ErrBalanceNotEnough
		}
		return nil, ErrBalanceInternal
	}
	return &domain.BalanceOperation{
		ID:          created.ID,
		BalanceID:   created.BalanceID,
		Amount:      created.Amount,
		Type:        created.Type,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
	}, nil
}
