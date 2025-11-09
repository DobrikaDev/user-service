package delivery

import (
	"DobrikaDev/user-service/internal/domain"
	"DobrikaDev/user-service/internal/generated/proto/user"
	"context"

	"github.com/dr3dnought/gospadi"
	"go.uber.org/zap"
)

func (s *Server) GetReputationGroups(ctx context.Context, req *user.GetReputationGroupsRequest) (*user.GetReputationGroupsResponse, error) {
	reputationGroups, err := s.reputationGroupService.GetReputationGroups(ctx)
	if err != nil {
		s.logger.Error("failed to get reputation groups", zap.Error(err))
		return &user.GetReputationGroupsResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	return &user.GetReputationGroupsResponse{
		ReputationGroups: convertReputationGroupsToProto(reputationGroups),
	}, nil
}

func (s *Server) GetReputationGroupByID(ctx context.Context, req *user.GetReputationGroupByIDRequest) (*user.GetReputationGroupByIDResponse, error) {
	if req.Id == 0 {
		return &user.GetReputationGroupByIDResponse{
			Error: &user.Error{
				Code:    user.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "id is required",
			},
		}, nil
	}
	reputationGroup, err := s.reputationGroupService.GetReputationGroupByID(ctx, int(req.Id))
	if err != nil {
		s.logger.Error("failed to get reputation group by id", zap.Error(err), zap.Int("id", int(req.Id)))
		return &user.GetReputationGroupByIDResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	return &user.GetReputationGroupByIDResponse{
		ReputationGroup: convertReputationGroupToProto(reputationGroup),
	}, nil
}

func convertReputationGroupToProto(reputationGroup *domain.ReputationGroup) *user.ReputationGroup {
	if reputationGroup == nil {
		return nil
	}
	return &user.ReputationGroup{
		Id:             int32(reputationGroup.ID),
		Name:           reputationGroup.Name,
		Description:    reputationGroup.Description,
		Coefficient:    reputationGroup.Coefficient,
		ReputationNeed: int32(reputationGroup.ReputationNeed),
	}
}
func convertReputationGroupsToProto(reputationGroups []*domain.ReputationGroup) []*user.ReputationGroup {
	if len(reputationGroups) == 0 {
		return nil
	}
	return gospadi.Map(reputationGroups, convertReputationGroupToProto)
}
