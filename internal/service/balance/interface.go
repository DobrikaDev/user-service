package balance

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/utils/config"
	"context"

	"go.uber.org/zap"
)

type storage interface {
	GetBalance(ctx context.Context, maxID string) (*domain.Balance, error)
	GetBalanceOperations(ctx context.Context, maxID string, limit int, offset int) ([]*domain.BalanceOperation, int32, error)
	CreateBalanceOperation(ctx context.Context, operation *domain.BalanceOperation) (*domain.BalanceOperation, error)
}

type BalanceService struct {
	storage storage
	cfg     *config.Config
	logger  *zap.Logger
}

func NewBalanceService(storage storage, cfg *config.Config, logger *zap.Logger) *BalanceService {
	return &BalanceService{storage: storage, cfg: cfg, logger: logger}
}
