package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo           *Repository
	sessionService *sessions.Service
}

// NewService creates a new users service.
func NewService(
	repo *Repository,
	sessionService *sessions.Service,
) *Service {
	return &Service{
		repo:           repo,
		sessionService: sessionService,
	}
}

// GetUserByID retrieves a user by their ID and returns a DTO.
func (s *Service) GetUserByID(
	ctx context.Context,
	userID string,
) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := s.mapToResponse(user)
	profilePicture, err := s.repo.GetProfilePictureURLByUserID(ctx, userID)
	if err == nil {
		resp.ProfilePicture = profilePicture
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	return resp, nil
}

// GetUserByEmail retrieves a user by their email and auth type.
func (s *Service) GetUserByEmail(
	ctx context.Context,
	email string,
	authType string,
) (*UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email, authType)
	if err != nil {
		return nil, err
	}

	resp := s.mapToResponse(user)
	profilePicture, err := s.repo.GetProfilePictureURLByUserID(ctx, user.ID)
	if err == nil {
		resp.ProfilePicture = profilePicture
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	return resp, nil
}

func (s *Service) GetUserIDsByRole(
	ctx context.Context,
	roleID int,
) ([]string, error) {
	return s.repo.GetUserIDsByRole(ctx, roleID)
}

func (s *Service) GetRoleDistribution(
	ctx context.Context,
) ([]RoleDistributionDTO, error) {
	return s.repo.GetRoleDistribution(ctx)
}

func (s *Service) GetRolesByUserID(
	ctx context.Context,
	userID string,
) ([]Role, error) {
	return s.repo.GetRolesByUserID(ctx, userID)
}

func (s *Service) GetRoleByID(
	ctx context.Context,
	roleID int,
) (*Role, error) {
	return s.repo.GetRoleByID(ctx, roleID)
}

func (s *Service) CheckUserWhitelist(
	ctx context.Context,
	email string,
) ([]int, error) {
	return s.repo.CheckUserWhitelist(ctx, email)
}

func (s *Service) ListUsers(
	ctx context.Context,
	params ListUsersRequest,
) (*ListUsersResponse, error) {
	users, total, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	dtos := make([]UserResponse, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, *s.mapToResponse(&u))
	}

	return &ListUsersResponse{
		Users: dtos,
		Meta:  structs.CalculateMetadata(total, params.Page, params.PageSize),
	}, nil
}

func (s *Service) PostProfilePicture(
	ctx context.Context,
	userID string,
	fileID string,
) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.PostProfilePicture(ctx, tx, userID, fileID)
		},
	)
}

func (s *Service) BlockUser(ctx context.Context, userID string) error {
	err := s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.BlockUser(ctx, tx, userID)
		},
	)
	if err != nil {
		return err
	}

	// Revoke sessions
	_ = s.sessionService.RevokeAllUserSessions(ctx, userID)

	return nil
}

func (s *Service) UnblockUser(ctx context.Context, userID string) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.UnblockUser(ctx, tx, userID)
		},
	)
}

func (s *Service) UpdateUserRoles(
	ctx context.Context,
	req UpdateRolesRequest,
	adminID string,
) error {
	return s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		// Remove current roles
		if err := s.repo.RemoveRoles(ctx, tx, req.UserID); err != nil {
			return fmt.Errorf("failed to remove old roles: %w", err)
		}

		// Assign new roles
		for _, roleID := range req.RoleIDs {
			assignment := RoleAssignment{
				UserID:      req.UserID,
				RoleID:      roleID,
				AssignedBy:  structs.StringToNullableString(adminID),
				Reason:      structs.StringToNullableString(req.Reason),
				ReferenceID: structs.StringToNullableString(req.ReferenceID),
			}
			if err := s.repo.AssignRole(ctx, tx, assignment); err != nil {
				return fmt.Errorf("failed to assign role %d: %w", roleID, err)
			}
		}

		return nil
	})
}

func (s *Service) AddUserToWhitelist(
	ctx context.Context,
	req AddUserToWhitelistRequest,
) error {
	return s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		for _, roleID := range req.RoleIDs {
			if err := s.repo.AddUserToWhitelist(
				ctx, tx, req.Email, roleID,
			); err != nil {
				return fmt.Errorf(
					"failed to add user to whitelist: %w", err,
				)
			}
		}
		return nil
	})
}

func (s *Service) RemoveUserFromWhitelist(
	ctx context.Context,
	req RemoveUserFromWhitelistRequest,
) error {
	return s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		if err := s.repo.RemoveUserFromWhitelist(ctx, tx, req.Email); err != nil {
			return fmt.Errorf("failed to remove user from whitelist: %w", err)
		}
		return nil
	})
}
func (s *Service) mapToResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Roles:      user.Roles,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		SuffixName: user.SuffixName,
		Email:      user.Email,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Time.String(),
		UpdatedAt:  user.UpdatedAt.Time.String(),
	}
}
