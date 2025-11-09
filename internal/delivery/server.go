package delivery

import (
	userpb "DobrikaDev/user-service/internal/generated/proto/user"
	"DobrikaDev/user-service/internal/service/balance"
	reputationgroup "DobrikaDev/user-service/internal/service/reputation_group"
	"DobrikaDev/user-service/internal/service/user"
	"DobrikaDev/user-service/utils/config"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	userService            *user.UserService
	reputationGroupService *reputationgroup.ReputationGroupService
	balanceService         *balance.BalanceService
	userpb.UnimplementedUserServiceServer

	cfg    *config.Config
	logger *zap.Logger
}

func NewServer(ctx context.Context, userService *user.UserService, reputationGroupService *reputationgroup.ReputationGroupService, balanceService *balance.BalanceService, cfg *config.Config, logger *zap.Logger) *Server {
	server := &Server{userService: userService, reputationGroupService: reputationGroupService, balanceService: balanceService, cfg: cfg, logger: logger}
	return server
}

func (s *Server) Register(grpcServer *grpc.Server) {
	userpb.RegisterUserServiceServer(grpcServer, s)
}
