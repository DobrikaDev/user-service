package sql

import (
	"DobrikaDev/user-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
)

type GetUsersResponse struct {
	Users []*domain.User `json:"users"`
	Total int            `json:"total"`
}

type userRow struct {
	MaxID             string    `db:"max_id"`
	Name              string    `db:"name"`
	Geolocation       string    `db:"geolocation"`
	Age               int       `db:"age"`
	Sex               string    `db:"sex"`
	About             string    `db:"about"`
	Role              string    `db:"role"`
	Status            string    `db:"status"`
	ReputationGroupID int       `db:"reputation_group_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	ReputationName    string    `db:"rg_name"`
	ReputationDesc    string    `db:"rg_description"`
	ReputationCoeff   float64   `db:"rg_coefficient"`
	ReputationNeed    int       `db:"rg_reputation_need"`
}

func (r *userRow) toDomain() *domain.User {
	return &domain.User{
		MaxID:             r.MaxID,
		Name:              r.Name,
		Geolocation:       r.Geolocation,
		Age:               r.Age,
		Sex:               domain.Sex(r.Sex),
		About:             r.About,
		Role:              domain.UserRole(r.Role),
		Status:            domain.UserStatus(r.Status),
		ReputationGroupID: r.ReputationGroupID,
		ReputationGroup: &domain.ReputationGroup{
			ID:             r.ReputationGroupID,
			Name:           r.ReputationName,
			Description:    r.ReputationDesc,
			Coefficient:    r.ReputationCoeff,
			ReputationNeed: r.ReputationNeed,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
func (s *SqlStorage) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.MaxID == "" {
		s.logger.Warn("max_id is required for user creation")
		return nil, ErrUserInvalid
	}

	if user.ReputationGroupID == 0 {
		user.ReputationGroupID = 1
	}

	now := time.Now().UTC()

	var created domain.User
	err := s.TransactionManager.Do(ctx, func(txCtx context.Context) error {
		tx := s.trf.Transaction(txCtx)

		ib := sq.Insert("users").
			Columns(
				"max_id",
				"name",
				"geolocation",
				"age",
				"sex",
				"about",
				"role",
				"status",
				"reputation_group_id",
				"created_at",
				"updated_at",
			).
			Values(
				user.MaxID,
				user.Name,
				user.Geolocation,
				user.Age,
				user.Sex,
				user.About,
				user.Role,
				user.Status,
				user.ReputationGroupID,
				now,
				now,
			).
			Suffix("RETURNING max_id, name, geolocation, age, sex, about, role, status, reputation_group_id, created_at, updated_at").
			PlaceholderFormat(sq.Dollar)

		q, args := ib.MustSql()

		if err := tx.QueryRowContext(txCtx, q, args...).Scan(
			&created.MaxID,
			&created.Name,
			&created.Geolocation,
			&created.Age,
			&created.Sex,
			&created.About,
			&created.Role,
			&created.Status,
			&created.ReputationGroupID,
			&created.CreatedAt,
			&created.UpdatedAt,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case pgErrUniqueViolation:
					return ErrUserAlreadyExists
				case pgErrForeignKeyViolation:
					return ErrReputationGroupNotFound
				}
			}

			s.logger.Error("failed to create user", zap.Error(err))
			return ErrUserInternal
		}

		balanceID := uuid.NewString()
		var group domain.ReputationGroup
		if err := tx.QueryRowContext(
			txCtx,
			"SELECT id, name, description, coefficient, reputation_need FROM reputation_groups WHERE id = $1",
			created.ReputationGroupID,
		).Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.Coefficient,
			&group.ReputationNeed,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrReputationGroupNotFound
			}
			s.logger.Error("failed to fetch reputation group for user", zap.Error(err), zap.Int("reputation_group_id", created.ReputationGroupID))
			return ErrUserInternal
		}

		if _, err := tx.ExecContext(
			txCtx,
			"INSERT INTO balances (id, user_id, balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
			balanceID,
			created.MaxID,
			0,
			now,
			now,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case pgErrUniqueViolation:
					return ErrBalanceAlreadyExists
				case pgErrForeignKeyViolation:
					return ErrUserInvalid
				}
			}

			s.logger.Error("failed to create balance for user", zap.Error(err), zap.String("max_id", created.MaxID))
			return ErrUserInternal
		}

		created.ReputationGroup = &group

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &created, nil
}

type ListUsersOpts func(sq.SelectBuilder) sq.SelectBuilder

func ListUsersWithLimit(limit int) ListUsersOpts {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if limit > 0 {
			return sb.Limit(uint64(limit))
		}
		return sb
	}
}

func ListUsersWithOffset(offset int) ListUsersOpts {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if offset > 0 {
			return sb.Offset(uint64(offset))
		}
		return sb
	}
}

func ListUsersWithMaxID(maxID string) ListUsersOpts {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		return sb.Where(sq.Eq{"u.max_id": maxID})
	}
}

func ListUsersWithMaxIDs(maxIDs []string) ListUsersOpts {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if len(maxIDs) == 0 {
			return sb
		}
		return sb.Where(sq.Eq{"u.max_id": maxIDs})
	}
}

func ListUsersWithStatuses(statuses []domain.UserStatus) ListUsersOpts {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if len(statuses) == 0 {
			return sb
		}

		values := make([]string, 0, len(statuses))
		for _, status := range statuses {
			values = append(values, string(status))
		}

		return sb.Where(sq.Eq{"u.status": values})
	}
}

func ListUsersWithRoles(roles []domain.UserRole) ListUsersOpts {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if len(roles) == 0 {
			return sb
		}

		values := make([]string, 0, len(roles))
		for _, role := range roles {
			values = append(values, string(role))
		}

		return sb.Where(sq.Eq{"u.role": values})
	}
}

func (s *SqlStorage) GetUsers(ctx context.Context, opts ...ListUsersOpts) (*GetUsersResponse, error) {
	sb := sq.Select(
		"u.max_id",
		"u.name",
		"u.geolocation",
		"u.age",
		"u.sex",
		"u.about",
		"u.role",
		"u.status",
		"u.reputation_group_id",
		"u.created_at",
		"u.updated_at",
		"rg.name AS rg_name",
		"rg.description AS rg_description",
		"rg.coefficient AS rg_coefficient",
		"rg.reputation_need AS rg_reputation_need",
	).
		From("users u").
		Join("reputation_groups rg ON rg.id = u.reputation_group_id").
		OrderBy("u.created_at DESC").
		PlaceholderFormat(sq.Dollar)

	for _, opt := range opts {
		sb = opt(sb)
	}

	q, args := sb.MustSql()

	rows := make([]userRow, 0, 10)
	err := s.trf.Transaction(ctx).SelectContext(ctx, &rows, q, args...)
	if err != nil {
		s.logger.Error("failed to get users", zap.Error(err))
		return nil, ErrUserInternal
	}

	users := make([]*domain.User, 0, len(rows))
	for i := range rows {
		users = append(users, rows[i].toDomain())
	}

	total, err := s.CountUsers(ctx)
	if err != nil {
		s.logger.Error("failed to get total users", zap.Error(err))
		return nil, ErrUserInternal
	}

	return &GetUsersResponse{Users: users, Total: total}, nil
}

func (s *SqlStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	ub := sq.Update("users").
		Set("name", user.Name).
		Set("geolocation", user.Geolocation).
		Set("age", user.Age).
		Set("sex", user.Sex).
		Set("about", user.About).
		Set("role", user.Role).
		Set("status", user.Status).
		Set("reputation_group_id", user.ReputationGroupID).
		Set("updated_at", sq.Expr("NOW() AT TIME ZONE 'UTC'")).
		Where(sq.Eq{"max_id": user.MaxID}).
		PlaceholderFormat(sq.Dollar)

	q, args := ub.MustSql()

	result, err := s.trf.Transaction(ctx).ExecContext(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrUniqueViolation:
				return ErrUserAlreadyExists
			case pgErrForeignKeyViolation:
				return ErrReputationGroupNotFound
			}
		}
		s.logger.Error("failed to update user", zap.Error(err))
		return ErrUserInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("failed to get rows affected", zap.Error(err))
		return ErrUserInternal
	}

	if rowsAffected == 0 {
		s.logger.Error("user not found", zap.String("max_id", user.MaxID))
		return ErrUserNotFound
	}

	return nil
}

func (s *SqlStorage) DeleteUser(ctx context.Context, maxID string) error {
	db := s.trf.Transaction(ctx)

	q, args := sq.Delete("users").
		Where(sq.Eq{"max_id": maxID}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	result, err := db.ExecContext(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrForeignKeyViolation:
				return ErrUserInvalid
			}
		}

		s.logger.Error("failed to delete user", zap.Error(err))
		return ErrUserInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("failed to get rows affected on delete", zap.Error(err))
		return ErrUserInternal
	}

	if rowsAffected == 0 {
		s.logger.Warn("user not found during delete", zap.String("max_id", maxID))
		return ErrUserNotFound
	}

	return nil
}

func (s *SqlStorage) CountUsers(ctx context.Context) (int, error) {
	sb := sq.Select("COUNT(*)").
		From("users").
		PlaceholderFormat(sq.Dollar)

	q, args := sb.MustSql()

	var total int
	err := s.trf.Transaction(ctx).QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		s.logger.Error("failed to get total users", zap.Error(err))
		return 0, ErrUserInternal
	}

	return total, nil
}
