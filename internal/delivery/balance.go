package delivery

import (
	"DobrikaDev/user-service/internal/domain"
	userpb "DobrikaDev/user-service/internal/generated/proto/user"
	"context"

	"github.com/dr3dnought/gospadi"
	"go.uber.org/zap"
)

func (s *Server) CreateOperation(ctx context.Context, req *userpb.CreateOperationRequest) (*userpb.CreateOperationResponse, error) {
	if req.Amount <= 0 {
		return &userpb.CreateOperationResponse{
			Error: &userpb.Error{
				Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "amount is required",
			},
		}, nil
	}
	if req.Type == userpb.BalanceOperationType_BALANCE_OPERATION_TYPE_UNSPECIFIED {
		return &userpb.CreateOperationResponse{
			Error: &userpb.Error{
				Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "type is required",
			},
		}, nil
	}

	operation, err := s.balanceService.CreateOperation(ctx, req.MaxId, &domain.BalanceOperation{
		Amount:      int(req.Amount),
		Type:        convertBalanceOperationTypeToDomain(req.Type),
		Description: req.Description,
	})
	if err != nil {
		s.logger.Error("failed to create balance operation", zap.Error(err), zap.String("max_id", req.MaxId))
		return &userpb.CreateOperationResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	return &userpb.CreateOperationResponse{
		Operation: &userpb.BalanceOperation{
			Id:          operation.ID,
			BalanceId:   operation.BalanceID,
			Amount:      int32(operation.Amount),
			Type:        convertBalanceOperationTypeToProto(operation.Type),
			Description: operation.Description,
			CreatedAt:   int32(operation.CreatedAt.Unix()),
		},
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *userpb.GetBalanceRequest) (*userpb.GetBalanceResponse, error) {
	if req.MaxId == "" {
		return &userpb.GetBalanceResponse{
			Error: &userpb.Error{
				Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max_id is required",
			},
		}, nil
	}
	balance, err := s.balanceService.GetBalance(ctx, req.MaxId)
	if err != nil {
		s.logger.Error("failed to get balance", zap.Error(err), zap.String("max_id", req.MaxId))
		return &userpb.GetBalanceResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	return &userpb.GetBalanceResponse{
		Balance: int32(balance.Balance),
	}, nil
}

func (s *Server) GetBalanceOperations(ctx context.Context, req *userpb.GetBalanceOperationsRequest) (*userpb.GetBalanceOperationsResponse, error) {
	if req.MaxId == "" {
		return &userpb.GetBalanceOperationsResponse{
			Error: &userpb.Error{
				Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max_id is required",
			},
		}, nil
	}
	operations, total, err := s.balanceService.GetBalanceOperations(ctx, req.MaxId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("failed to get balance operations", zap.Error(err), zap.String("max_id", req.MaxId), zap.Int("limit", int(req.Limit)), zap.Int("offset", int(req.Offset)))
		return &userpb.GetBalanceOperationsResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	return &userpb.GetBalanceOperationsResponse{
		Operations: convertBalanceOperationsToProto(operations),
		Total:      total,
	}, nil
}

func convertBalanceOperationTypeToDomain(t userpb.BalanceOperationType) domain.BalanceOperationType {
	switch t {
	case userpb.BalanceOperationType_BALANCE_OPERATION_TYPE_DEPOSIT:
		return domain.BalanceOperationTypeDeposit
	case userpb.BalanceOperationType_BALANCE_OPERATION_TYPE_WITHDRAW:
		return domain.BalanceOperationTypeWithdraw
	default:
		return domain.BalanceOperationTypeWithdraw
	}
}

func convertBalanceOperationsToProto(operations []*domain.BalanceOperation) []*userpb.BalanceOperation {
	return gospadi.Map(operations, func(operation *domain.BalanceOperation) *userpb.BalanceOperation {
		return &userpb.BalanceOperation{
			Id:          operation.ID,
			BalanceId:   operation.BalanceID,
			Amount:      int32(operation.Amount),
			Type:        convertBalanceOperationTypeToProto(operation.Type),
			Description: operation.Description,
			CreatedAt:   int32(operation.CreatedAt.Unix()),
		}
	})
}

func convertBalanceOperationTypeToProto(t domain.BalanceOperationType) userpb.BalanceOperationType {
	switch t {
	case domain.BalanceOperationTypeDeposit:
		return userpb.BalanceOperationType_BALANCE_OPERATION_TYPE_DEPOSIT
	case domain.BalanceOperationTypeWithdraw:
		return userpb.BalanceOperationType_BALANCE_OPERATION_TYPE_WITHDRAW
	default:
		return userpb.BalanceOperationType_BALANCE_OPERATION_TYPE_UNSPECIFIED
	}
}
