package sql

import (
	"DobrikaDev/user-service/internal/domain"
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (s *SqlStorage) GetReputationGroups(ctx context.Context) ([]*domain.ReputationGroup, error) {
	sb := sq.Select("id", "name", "description", "coefficient", "reputation_need").
		From("reputation_groups").
		OrderBy("reputation_need ASC").
		PlaceholderFormat(sq.Dollar)

	query, args := sb.MustSql()

	groups := make([]*domain.ReputationGroup, 0, 4)
	if err := s.trf.Transaction(ctx).SelectContext(ctx, &groups, query, args...); err != nil {
		s.logger.Error("failed to get reputation groups", zap.Error(err))
		return nil, ErrReputationGroupInternal
	}

	return groups, nil
}

func (s *SqlStorage) GetReputationGroupByID(ctx context.Context, id int) (*domain.ReputationGroup, error) {
	sb := sq.Select("id", "name", "description", "coefficient", "reputation_need").
		From("reputation_groups").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args := sb.MustSql()

	var group domain.ReputationGroup
	if err := s.trf.Transaction(ctx).GetContext(ctx, &group, query, args...); err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("reputation group not found", zap.Int("id", id))
			return nil, ErrReputationGroupNotFound
		}

		s.logger.Error("failed to get reputation group", zap.Error(err), zap.Int("id", id))
		return nil, ErrReputationGroupInternal
	}

	return &group, nil
}
