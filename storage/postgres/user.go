package postgres

import (
	"auth/models"
	"auth/pkg/helper"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (b *userRepo) CreateUser(c context.Context, req *models.CreateUser) (string, error) {
	id := uuid.NewString()

	query := `
		INSERT INTO "users"(
			"id",
			"username", 
			"password", 
			"created_at")
		VALUES ($1, $2, $3, NOW())
	`
	_, err := b.db.Exec(context.Background(), query,
		id,
		req.Username,
		req.Password,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (b *userRepo) GetUser(c context.Context, req *models.IdRequest) (resp *models.User, err error) {

	var (
		created_at sql.NullTime
		updated_at sql.NullTime
		deleted_at sql.NullTime
	)
	query := `
			SELECT 
				"id", 
				"username", 
				"password", 
				"is_active", 
				"created_at",
				"updated_at",
				"deleted_at"
			FROM "users" 
				WHERE 
				"is_active" = true AND
			    "id"=$1`

	user := models.User{}
	err = b.db.QueryRow(context.Background(), query, req.Id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Is_active,
		&created_at,
		&updated_at,
		&deleted_at,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	user.CreatedAt = created_at.Time.Format(time.RFC3339)
	if updated_at.Valid {
		user.UpdatedAt = updated_at.Time.Format(time.RFC3339)
	}

	if deleted_at.Valid {
		user.DeletedAt = deleted_at.Time.Format(time.RFC3339)
	}

	return &user, nil
}

func (b *userRepo) GetAllActiveUser(c context.Context, req *models.GetAllUserRequest) (*models.GetAllUser, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllUser{}

	resp.Users = make([]models.User, 0)

	filter := " WHERE is_active = true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"username", 
				"password", 
				"is_active", 
				"created_at",
				"updated_at",
				"deleted_at"
			FROM "users"
		`
	if req.Search != "" {
		filter += ` AND "username" ILIKE '%' || :username || '%' `
		params["username"] = req.Search
	}

	offset := (req.Page - 1) * req.Limit
	params["limit"] = req.Limit
	params["offset"] = offset

	query = query + filter + " ORDER BY created_at DESC OFFSET :offset LIMIT :limit "
	rquery, pArr := helper.ReplaceQueryParams(query, params)

	rows, err := b.db.Query(context.Background(), rquery, pArr...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			created_at sql.NullTime
			updated_at sql.NullTime
			deleted_at sql.NullTime
		)
		user := models.User{}

		err := rows.Scan(
			&resp.Count,
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Is_active,
			&created_at,
			&updated_at,
			&deleted_at,
		)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = updated_at.Time.Format(time.RFC3339)
		if updated_at.Valid {
			user.UpdatedAt = updated_at.Time.Format(time.RFC3339)
		}
		if deleted_at.Valid {
			user.DeletedAt = deleted_at.Time.Format(time.RFC3339)
		}

		resp.Users = append(resp.Users, user)
	}
	return resp, nil
}

func (b *userRepo) UpdateUser(c context.Context, req *models.UpdateUser) (string, error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
			UPDATE users 
				SET 
				"username" = $1,
				"password" = $2,
				"updated_at" = NOW() 
				WHERE 
				"is_active" = true AND
				"id" = $3`

	result, err := b.db.Exec(
		context.Background(),
		query,
		req.Username,
		req.Password,
		userInfo.User_id,
	)

	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user with ID %s not found", req.ID)
	}

	return req.ID, nil
}

func (b *userRepo) DeleteUser(c context.Context, req *models.IdRequest) (resp string, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "users" 
		SET 
			"is_active" = $1,
			"deleted_at" = NOW() 
		WHERE 
			"is_active" = true AND
			"id" = $2
`
	result, err := b.db.Exec(
		context.Background(),
		query,
		false,
		userInfo.User_id,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user with ID %s not found", req.Id)

	}

	return req.Id, nil
}

func (b *userRepo) GetAllDeletedUser(c context.Context, req *models.GetAllUserRequest) (*models.GetAllUser, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllUser{}

	resp.Users = make([]models.User, 0)

	filter := " WHERE is_active != true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"username", 
				"password", 
				"is_active", 
				"created_at",
				"updated_at",
				"deleted_at"
			FROM "users"
		`
	if req.Search != "" {
		filter += ` AND "username" ILIKE '%' || :username || '%' `
		params["username"] = req.Search
	}

	offset := (req.Page - 1) * req.Limit
	params["limit"] = req.Limit
	params["offset"] = offset

	query = query + filter + " ORDER BY created_at DESC OFFSET :offset LIMIT :limit "
	rquery, pArr := helper.ReplaceQueryParams(query, params)

	rows, err := b.db.Query(context.Background(), rquery, pArr...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			created_at sql.NullTime
			updated_at sql.NullTime
			deleted_at sql.NullTime
		)
		user := models.User{}

		err := rows.Scan(
			&resp.Count,
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Is_active,
			&created_at,
			&updated_at,
			&deleted_at,
		)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = created_at.Time.Format(time.RFC3339)
		if updated_at.Valid {
			user.UpdatedAt = updated_at.Time.Format(time.RFC3339)
		}
		if deleted_at.Valid {
			user.DeletedAt = deleted_at.Time.Format(time.RFC3339)
		}

		resp.Users = append(resp.Users, user)
	}
	return resp, nil
}

func (b *userRepo) GetByUsername(c context.Context, req *models.LoginRequest) (resp *models.LoginDataRespond, err error) {

	query := `
			SELECT 
			"id",
				"username", 
				"password"
			FROM "users" 
				WHERE "username"=$1`

	user := models.LoginDataRespond{}
	err = b.db.QueryRow(context.Background(), query, req.Username).Scan(
		&user.User_id,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
