package users

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

// GetUserByID fetches a user by their ID.
func (r *Repository) GetUserByID(
	ctx context.Context,
	userID string,
) (*User, error) {
	var user User

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE id = ?
		LIMIT 1
	`, datastore.GetColumns(User{}))

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}

	roles, err := r.GetRolesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}
	user.Roles = roles

	profileURL, err := r.GetProfilePictureURLByUserID(ctx, userID)
	if err == nil {
		user.ProfilePicture = structs.StringToNullableString(profileURL)
	}

	return &user, nil
}

func (r *Repository) CheckUserWhitelist(
	ctx context.Context,
	email string,
) ([]int, error) {
	var roleIDs []int
	query := `SELECT role_id FROM whitelists WHERE email = ?`
	err := r.db.SelectContext(ctx, &roleIDs, query, email)
	return roleIDs, err
}

func (r *Repository) GetRoleByID(
	ctx context.Context,
	roleID int,
) (*Role, error) {
	var role Role
	query := fmt.Sprintf(`
		SELECT %s
		FROM roles
		WHERE id = ?
	`, datastore.GetColumns(Role{}))

	err := r.db.GetContext(ctx, &role, query, roleID)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *Repository) GetRolesByUserID(
	ctx context.Context,
	userID string,
) ([]Role, error) {
	var roles []Role
	query := `
		SELECT r.id, r.name
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	err := r.db.SelectContext(ctx, &roles, query, userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context,
	email, authType string,
) (*User, error) {
	var user User

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE email = ? AND auth_type = ?
		LIMIT 1
	`, datastore.GetColumns(User{}))

	err := r.db.GetContext(ctx, &user, query, email, authType)
	if err != nil {
		return nil, err
	}

	roles, err := r.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}
	user.Roles = roles

	profileURL, err := r.GetProfilePictureURLByUserID(ctx, user.ID)
	if err == nil {
		user.ProfilePicture = structs.StringToNullableString(profileURL)
	}

	return &user, nil
}

// CreateUser inserts or updates a user.
func (r *Repository) CreateUser(
	ctx context.Context,
	tx datastore.DB,
	user User,
) error {
	exclude := []string{"updated_at"}
	cols, vals := datastore.GetInsertStatement(User{}, exclude)
	onDuplicateKeyStmt := datastore.GetOnDuplicateKeyUpdateStatement(
		User{},
		exclude,
	)
	query := fmt.Sprintf(`
			INSERT INTO users (id, %s)
			VALUES (:id, %s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, onDuplicateKeyStmt)

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	// Insert roles if any
	for _, role := range user.Roles {
		assignment := RoleAssignment{
			UserID: user.ID,
			RoleID: role.ID,
			Reason: structs.StringToNullableString("Initial account creation"),
		}
		if err := r.AssignRole(ctx, tx, assignment); err != nil {
			return fmt.Errorf("failed to assign initial role: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetProfilePictureURLByUserID(
	ctx context.Context,
	userID string,
) (string, error) {
	var fileURL string

	query := `
		SELECT f.file_url
		FROM profile_pictures pp
		JOIN files f ON f.id = pp.file_id
		WHERE pp.user_id = ?
		ORDER BY f.created_at DESC
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &fileURL, query, userID)
	if err != nil {
		return "", err
	}

	return fileURL, nil
}
 
func (r *Repository) GetStudentCORURLByUserID(
	ctx context.Context,
	userID string,
) (string, error) {
	var fileURL string

	query := `
		SELECT f.file_url
		FROM student_cors sc
		JOIN files f ON f.id = sc.file_id
		WHERE sc.student_id = ?
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &fileURL, query, userID)
	if err != nil {
		return "", err
	}

	return fileURL, nil
}

func (r *Repository) PostProfilePicture(
	ctx context.Context,
	tx datastore.DB,
	userID string,
	fileID string,
) error {
	deleteQuery := `DELETE FROM profile_pictures WHERE user_id = ?`
	if _, err := tx.ExecContext(ctx, deleteQuery, userID); err != nil {
		return fmt.Errorf("failed to delete old profile picture: %w", err)
	}

	insertQuery := `
		INSERT INTO profile_pictures (user_id, file_id)
		VALUES (?, ?)
	`

	_, err := tx.ExecContext(ctx, insertQuery, userID, fileID)
	if err != nil {
		return fmt.Errorf("failed to post profile picture: %w", err)
	}

	return nil
}

func (r *Repository) BlockUser(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	query := `UPDATE users SET is_active = 0 WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) UnblockUser(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	query := `UPDATE users SET is_active = 1 WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) GetUserIDsByRole(
	ctx context.Context,
	roleID int,
) ([]string, error) {
	var userIDs []string
	query := `SELECT user_id FROM user_roles WHERE role_id = ?`
	err := r.db.SelectContext(ctx, &userIDs, query, roleID)
	return userIDs, err
}

func (r *Repository) ListUsers(
	ctx context.Context,
	params ListUsersRequest,
) ([]User, int, error) {
	var users []User
	var total int

	baseQuery := `FROM users u`
	whereClause := ` WHERE 1=1`
	args := []interface{}{}

	if params.RoleID > 0 {
		baseQuery += ` JOIN user_roles ur ON ur.user_id = u.id`
		whereClause += ` AND ur.role_id = ?`
		args = append(args, params.RoleID)
	}

	if params.Search != "" {
		whereClause += ` AND (u.email LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ?)`
		like := "%" + params.Search + "%"
		args = append(args, like, like, like)
	}

	if params.Active != nil {
		whereClause += ` AND u.is_active = ?`
		activeVal := 0
		if *params.Active {
			activeVal = 1
		}
		args = append(args, activeVal)
	}

	// Count total
	countQuery := `SELECT COUNT(DISTINCT u.id) ` + baseQuery + whereClause
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated users
	selectQuery := fmt.Sprintf(
		`SELECT DISTINCT %s `,
		datastore.GetPrefixColumns(User{}, "u")) +
		baseQuery + whereClause + " ORDER BY u.created_at DESC LIMIT ? OFFSET ?"

	limit := params.PageSize
	offset := (params.Page - 1) * params.PageSize
	args = append(args, limit, offset)

	err = r.db.SelectContext(ctx, &users, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	if len(users) == 0 {
		return []User{}, 0, nil
	}

	// Fetch roles for all users in one go to avoid N+1
	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	query, roleArgs, err := sqlx.In(`
		SELECT ur.user_id, r.id, r.name
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id IN (?)
	`, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build role query: %w", err)
	}

	type UserRoleRow struct {
		UserID string `db:"user_id"`
		ID     int    `db:"id"`
		Name   string `db:"name"`
	}
	var roleRows []UserRoleRow
	err = r.db.SelectContext(ctx, &roleRows, query, roleArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch roles for list: %w", err)
	}

	userRolesMap := make(map[string][]Role)
	for _, row := range roleRows {
		userRolesMap[row.UserID] = append(userRolesMap[row.UserID], Role{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	for i := range users {
		users[i].Roles = userRolesMap[users[i].ID]
	}

	// Fetch profile pictures in one go
	picQuery, picArgs, err := sqlx.In(`
		SELECT pp.user_id, f.file_url
		FROM profile_pictures pp
		JOIN files f ON f.id = pp.file_id
		WHERE pp.user_id IN (?)
	`, userIDs)
	if err == nil {
		picQuery = r.db.Rebind(picQuery)
		var pics []struct {
			UserID  string `db:"user_id"`
			FileURL string `db:"file_url"`
		}
		if err := r.db.SelectContext(ctx, &pics, picQuery, picArgs...); err == nil {
			picMap := make(map[string]string)
			for _, p := range pics {
				picMap[p.UserID] = p.FileURL
			}
			for i := range users {
				if url, ok := picMap[users[i].ID]; ok {
					users[i].ProfilePicture = structs.StringToNullableString(url)
				}
			}
		}
	}

	return users, total, nil
}

func (r *Repository) GetRoleDistribution(
	ctx context.Context,
) ([]RoleDistributionDTO, error) {
	query := `
		SELECT r.name as role_name, COUNT(ur.user_id) as count
		FROM roles r
		LEFT JOIN user_roles ur ON ur.role_id = r.id
		GROUP BY r.name
	`

	var distribution []RoleDistributionDTO
	err := r.db.SelectContext(ctx, &distribution, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get role distribution: %w", err)
	}

	return distribution, nil
}

func (r *Repository) AssignRole(
	ctx context.Context,
	tx datastore.DB,
	assignment RoleAssignment,
) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_by, reason, reference_id)
		VALUES (:user_id, :role_id, :assigned_by, :reason, :reference_id)
		ON DUPLICATE KEY UPDATE
			assigned_by = VALUES(assigned_by),
			reason = VALUES(reason),
			reference_id = VALUES(reference_id)
	`

	_, err := tx.NamedExecContext(ctx, query, assignment)
	return err
}

func (r *Repository) RemoveRoles(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	query := `DELETE FROM user_roles WHERE user_id = ?`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) AddUserToWhitelist(
	ctx context.Context,
	tx datastore.DB,
	email string,
	roleID int,
) error {
	query := `INSERT INTO whitelists (email, role_id) VALUES (?, ?)`
	_, err := tx.ExecContext(ctx, query, email, roleID)
	return err
}

func (r *Repository) RemoveUserFromWhitelist(
	ctx context.Context,
	tx datastore.DB,
	email string,
) error {
	query := `DELETE FROM whitelists WHERE email = ?`
	_, err := tx.ExecContext(ctx, query, email)
	return err
}
