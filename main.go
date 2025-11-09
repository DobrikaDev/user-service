package main

import (
	"DobrikaDev/user-service/di"
	userpb "DobrikaDev/user-service/internal/generated/proto/user"
	"DobrikaDev/user-service/utils/config"
	"DobrikaDev/user-service/utils/logger"
	"context"
	"os"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoadConfigFromFile("deployments/config.yaml")
	logger, _ := logger.NewLogger()
	defer logger.Sync()

	container := di.NewContainer(ctx, cfg, logger)

	userpb.RegisterUserServiceServer(
		container.GetGRPCServer(),
		container.GetRpcServer(),
	)

	logger.Info("Starting application with port", zap.String("port", cfg.Port))

	err := container.GetGRPCServer().Serve(*container.GetNetListener())
	if err != nil {
		logger.Error("Error while serving grpcServer:", zap.Error(err))
		os.Exit(1)
	}
}
