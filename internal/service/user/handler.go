package user

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/internal/storage/sql"
	"context"
	"errors"

	"go.uber.org/zap"
)

type GetUsersResponse struct {
	Users []*domain.User `json:"users"`
	Total int            `json:"total"`
}

type GetUsersFilter struct {
	MaxIDs   []string
	MaxID    string
	Statuses []domain.UserStatus
	Roles    []domain.UserRole
	Limit    int
	Offset   int
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	created, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err), zap.Any("user", user))
		if errors.Is(err, sql.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		if errors.Is(err, sql.ErrReputationGroupNotFound) || errors.Is(err, sql.ErrUserInvalid) {
			return nil, ErrUserInvalid
		}
		return nil, ErrUserInternal
	}

	return created, nil
}

func (s *UserService) GetUsers(ctx context.Context, filter GetUsersFilter) (*GetUsersResponse, error) {
	opts := make([]sql.ListUsersOpts, 0, 6)

	if filter.Limit > 0 {
		opts = append(opts, sql.ListUsersWithLimit(filter.Limit))
	}
	if filter.Offset > 0 {
		opts = append(opts, sql.ListUsersWithOffset(filter.Offset))
	}
	if len(filter.MaxIDs) > 0 {
		opts = append(opts, sql.ListUsersWithMaxIDs(filter.MaxIDs))
	}
	if filter.MaxID != "" {
		opts = append(opts, sql.ListUsersWithMaxID(filter.MaxID))
	}
	if len(filter.Statuses) > 0 {
		opts = append(opts, sql.ListUsersWithStatuses(filter.Statuses))
	}
	if len(filter.Roles) > 0 {
		opts = append(opts, sql.ListUsersWithRoles(filter.Roles))
	}

	response, err := s.storage.GetUsers(ctx, opts...)
	if err != nil {
		s.logger.Error("failed to get users", zap.Error(err), zap.Any("filter", filter))
		return nil, ErrUserInternal
	}

	return &GetUsersResponse{Users: response.Users, Total: response.Total}, nil
}

func (s *UserService) GetUserByMaxID(ctx context.Context, maxID string) (*domain.User, error) {
	response, err := s.storage.GetUsers(ctx, sql.ListUsersWithMaxID(maxID))
	if err != nil {
		s.logger.Error("failed to get user by max id", zap.Error(err), zap.String("max_id", maxID))
		if errors.Is(err, sql.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrUserInternal
	}

	if len(response.Users) == 0 {
		s.logger.Warn("user not found by max id", zap.String("max_id", maxID))
		return nil, ErrUserNotFound
	}

	return response.Users[0], nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	err := s.storage.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to update user", zap.Error(err), zap.Any("user", user))
		if errors.Is(err, sql.ErrUserAlreadyExists) {
			return ErrUserAlreadyExists
		}
		if errors.Is(err, sql.ErrReputationGroupNotFound) {
			return ErrUserInvalid
		}
		return ErrUserInternal
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, maxID string) error {
	err := s.storage.DeleteUser(ctx, maxID)
	if err != nil {
		s.logger.Error("failed to delete user", zap.Error(err), zap.String("max_id", maxID))
		if errors.Is(err, sql.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return ErrUserInternal
	}

	return nil
}
