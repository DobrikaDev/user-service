package reputationgroup

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/utils/config"
	"context"

	"go.uber.org/zap"
)

type storage interface {
	GetReputationGroups(ctx context.Context) ([]*domain.ReputationGroup, error)
	GetReputationGroupByID(ctx context.Context, id int) (*domain.ReputationGroup, error)
}

type ReputationGroupService struct {
	storage storage
	cfg     *config.Config
	logger  *zap.Logger
}

func NewReputationGroupService(storage storage, cfg *config.Config, logger *zap.Logger) *ReputationGroupService {
	return &ReputationGroupService{storage: storage, cfg: cfg, logger: logger}
}
