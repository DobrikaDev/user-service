package user

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/internal/storage/sql"
	"DobrikaDev/user-service/utils/config"
	"context"

	"go.uber.org/zap"
)

type storage interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUsers(ctx context.Context, opts ...sql.ListUsersOpts) (*sql.GetUsersResponse, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, maxID string) error
}

type UserService struct {
	storage storage
	cfg     *config.Config
	logger  *zap.Logger
}

func NewUserService(storage storage, cfg *config.Config, logger *zap.Logger) *UserService {
	return &UserService{storage: storage, cfg: cfg, logger: logger}
}
