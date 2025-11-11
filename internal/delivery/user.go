package delivery

import (
	"context"
	"strings"

	"DobrikaDev/user-service/internal/domain"
	userpb "DobrikaDev/user-service/internal/generated/proto/user"
	balance "DobrikaDev/user-service/internal/service/balance"
	reputationgroup "DobrikaDev/user-service/internal/service/reputation_group"
	"DobrikaDev/user-service/internal/service/user"

	"github.com/dr3dnought/gospadi"
	"go.uber.org/zap"
)

func (s *Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user, err := s.userService.CreateUser(ctx, convertUserToDomain(req.User))
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err), zap.Any("user", req.User))
		return &userpb.CreateUserResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Debug("user created", zap.Any("user", user))

	return &userpb.CreateUserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *Server) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	filter := user.GetUsersFilter{
		MaxID:  req.MaxId,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	if req.Status != userpb.Status_STATUS_UNSPECIFIED {
		filter.Statuses = []domain.UserStatus{convertStatusToDomain(req.Status)}
	}

	if req.Role != userpb.Role_ROLE_UNSPECIFIED {
		filter.Roles = []domain.UserRole{convertRoleToDomain(req.Role)}
	}

	users, err := s.userService.GetUsers(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get users", zap.Error(err), zap.Any("filter", req))
		return &userpb.GetUsersResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Debug("users fetched", zap.Any("users", users))

	return &userpb.GetUsersResponse{
		Users: gospadi.Map(users.Users, convertUserToProto),
		Total: int32(users.Total),
	}, nil
}

func (s *Server) GetUserByMaxID(ctx context.Context, req *userpb.GetUserByMaxIDRequest) (*userpb.GetUserByMaxIDResponse, error) {
	user, err := s.userService.GetUserByMaxID(ctx, req.MaxId)
	if err != nil {
		s.logger.Error("failed to get user by max id", zap.Error(err), zap.String("max_id", req.MaxId))
		return &userpb.GetUserByMaxIDResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Debug("user by max id fetched", zap.Any("user", user))

	return &userpb.GetUserByMaxIDResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	if req.GetUser() == nil || strings.TrimSpace(req.GetUser().GetMaxId()) == "" {
		return &userpb.UpdateUserResponse{
			Error: &userpb.Error{
				Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max_id is required",
			},
		}, nil
	}

	existing, err := s.userService.GetUserByMaxID(ctx, req.GetUser().GetMaxId())
	if err != nil {
		s.logger.Error("failed to fetch user before update", zap.Error(err), zap.String("max_id", req.GetUser().GetMaxId()))
		return &userpb.UpdateUserResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	merged := mergeUser(existing, req.GetUser())

	err = s.userService.UpdateUser(ctx, merged)
	if err != nil {
		s.logger.Error("failed to update user", zap.Error(err), zap.Any("user", req.User))
		return &userpb.UpdateUserResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Debug("user updated", zap.Any("user", merged))

	return &userpb.UpdateUserResponse{
		User: convertUserToProto(merged),
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {

	err := s.userService.DeleteUser(ctx, req.MaxId)
	if err != nil {
		s.logger.Error("failed to delete user", zap.Error(err), zap.String("max_id", req.MaxId))
		return &userpb.DeleteUserResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Debug("user deleted", zap.String("max_id", req.MaxId))

	return &userpb.DeleteUserResponse{}, nil
}

func convertErrorToProto(err error) *userpb.Error {
	switch err {
	case user.ErrUserNotFound:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case user.ErrUserAlreadyExists:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
			Message: err.Error(),
		}
	case user.ErrUserInvalid:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
			Message: err.Error(),
		}
	case user.ErrUserInternal:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	case reputationgroup.ErrReputationGroupNotFound:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case reputationgroup.ErrReputationGroupAlreadyExists:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
			Message: err.Error(),
		}
	case reputationgroup.ErrReputationGroupInvalid:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_VALIDATION,
			Message: err.Error(),
		}
	case reputationgroup.ErrReputationGroupInternal:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	case balance.ErrBalanceNotFound:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case balance.ErrBalanceInternal:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	case balance.ErrBalanceNotEnough:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_NOT_ENOUGH,
			Message: err.Error(),
		}
	default:
		return &userpb.Error{
			Code:    userpb.ErrorCode_ERROR_CODE_UNSPECIFIED,
			Message: err.Error(),
		}
	}
}

func convertUserToProto(user *domain.User) *userpb.User {
	var reputationGroup *userpb.ReputationGroup
	if user.ReputationGroup != nil {
		reputationGroup = &userpb.ReputationGroup{
			Id:             int32(user.ReputationGroup.ID),
			Name:           user.ReputationGroup.Name,
			Description:    user.ReputationGroup.Description,
			Coefficient:    user.ReputationGroup.Coefficient,
			ReputationNeed: int32(user.ReputationGroup.ReputationNeed),
		}
	}

	return &userpb.User{
		MaxId:           user.MaxID,
		Name:            user.Name,
		Geolocation:     user.Geolocation,
		Age:             int32(user.Age),
		Sex:             convertSexToProto(user.Sex),
		About:           user.About,
		Role:            convertRoleToProto(user.Role),
		Status:          convertStatusToProto(user.Status),
		ReputationGroup: reputationGroup,
	}
}

func convertUserToDomain(user *userpb.User) *domain.User {
	domainUser := &domain.User{
		MaxID:       user.MaxId,
		Name:        user.Name,
		Geolocation: user.Geolocation,
		Age:         int(user.Age),
		Sex:         convertSexToDomain(user.Sex),
		About:       user.About,
		Role:        convertRoleToDomain(user.Role),
		Status:      convertStatusToDomain(user.Status),
	}

	if user.ReputationGroup != nil {
		domainUser.ReputationGroupID = int(user.ReputationGroup.Id)
		domainUser.ReputationGroup = &domain.ReputationGroup{
			ID:             int(user.ReputationGroup.Id),
			Name:           user.ReputationGroup.Name,
			Description:    user.ReputationGroup.Description,
			Coefficient:    user.ReputationGroup.Coefficient,
			ReputationNeed: int(user.ReputationGroup.ReputationNeed),
		}
	}

	return domainUser
}

func mergeUser(existing *domain.User, incoming *userpb.User) *domain.User {
	merged := *existing

	if incoming == nil {
		return &merged
	}

	if name := strings.TrimSpace(incoming.GetName()); name != "" {
		merged.Name = name
	}
	if geo := strings.TrimSpace(incoming.GetGeolocation()); geo != "" {
		merged.Geolocation = geo
	}
	if incoming.GetAge() > 0 {
		merged.Age = int(incoming.GetAge())
	}
	if incoming.GetSex() != userpb.Sex_SEX_UNSPECIFIED {
		merged.Sex = convertSexToDomain(incoming.GetSex())
	}
	if about := strings.TrimSpace(incoming.GetAbout()); about != "" {
		merged.About = about
	}
	if incoming.GetRole() != userpb.Role_ROLE_UNSPECIFIED {
		merged.Role = convertRoleToDomain(incoming.GetRole())
	}
	if incoming.GetStatus() != userpb.Status_STATUS_UNSPECIFIED {
		merged.Status = convertStatusToDomain(incoming.GetStatus())
	}
	if rg := incoming.GetReputationGroup(); rg != nil && rg.GetId() > 0 {
		merged.ReputationGroupID = int(rg.GetId())
		merged.ReputationGroup = &domain.ReputationGroup{
			ID:             int(rg.GetId()),
			Name:           rg.GetName(),
			Description:    rg.GetDescription(),
			Coefficient:    rg.GetCoefficient(),
			ReputationNeed: int(rg.GetReputationNeed()),
		}
	}

	return &merged
}

func convertStatusToDomain(status userpb.Status) domain.UserStatus {
	switch status {
	case userpb.Status_STATUS_ACTIVE:
		return domain.UserStatusActive
	case userpb.Status_STATUS_INACTIVE:
		return domain.UserStatusInactive
	default:
		return domain.UserStatusInactive
	}
}

func convertSexToDomain(sex userpb.Sex) domain.Sex {
	switch sex {
	case userpb.Sex_SEX_MALE:
		return domain.SexMale
	case userpb.Sex_SEX_FEMALE:
		return domain.SexFemale
	default:
		return domain.SexUnknown
	}
}

func convertRoleToDomain(role userpb.Role) domain.UserRole {
	switch role {
	case userpb.Role_ROLE_USER:
		return domain.UserRoleUser
	case userpb.Role_ROLE_ADMIN:
		return domain.UserRoleAdmin
	default:
		return domain.UserRoleUser
	}
}

func convertSexToProto(sex domain.Sex) userpb.Sex {
	switch sex {
	case domain.SexMale:
		return userpb.Sex_SEX_MALE
	case domain.SexFemale:
		return userpb.Sex_SEX_FEMALE
	default:
		return userpb.Sex_SEX_UNSPECIFIED
	}
}

func convertRoleToProto(role domain.UserRole) userpb.Role {
	switch role {
	case domain.UserRoleUser:
		return userpb.Role_ROLE_USER
	case domain.UserRoleAdmin:
		return userpb.Role_ROLE_ADMIN
	default:
		return userpb.Role_ROLE_USER
	}
}

func convertStatusToProto(status domain.UserStatus) userpb.Status {
	switch status {
	case domain.UserStatusActive:
		return userpb.Status_STATUS_ACTIVE
	case domain.UserStatusInactive:
		return userpb.Status_STATUS_INACTIVE
	default:
		return userpb.Status_STATUS_UNSPECIFIED
	}
}
