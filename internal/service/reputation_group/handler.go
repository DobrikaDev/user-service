package reputationgroup

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/internal/storage/sql"
	"context"
	"errors"

	"go.uber.org/zap"
)

func (s *ReputationGroupService) GetReputationGroups(ctx context.Context) ([]*domain.ReputationGroup, error) {
	groups, err := s.storage.GetReputationGroups(ctx)
	if err != nil {
		s.logger.Error("failed to get reputation groups", zap.Error(err))
		return nil, err
	}
	return groups, nil
}

func (s *ReputationGroupService) GetReputationGroupByID(ctx context.Context, id int) (*domain.ReputationGroup, error) {
	group, err := s.storage.GetReputationGroupByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get reputation group by id", zap.Error(err), zap.Int("id", id))
		if errors.Is(err, sql.ErrReputationGroupNotFound) {
			return nil, ErrReputationGroupNotFound
		}
		return nil, err
	}
	return group, nil
}
