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

type postRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) *postRepo {
	return &postRepo{
		db: db,
	}
}

func (b *postRepo) CreatePost(c context.Context, req *models.CreatePost) (string, error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)
	id := uuid.NewString()

	query := `
		INSERT INTO "post"(
			"id",
			"description", 
			"photos", 
			"created_by",
			"created_at"
			)
			
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := b.db.Exec(context.Background(), query,
		id,
		req.Description,
		req.Photos,
		userInfo.User_id,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create post: %w", err)
	}

	return id, nil
}

func (b *postRepo) GetPost(c context.Context, req *models.IdRequest) (resp *models.Post, err error) {

	var (
		created_at sql.NullTime
		created_by sql.NullString
		updated_at sql.NullTime
		updated_by sql.NullString
		deleted_at sql.NullTime
		deleted_by sql.NullString
	)

	query := `
			SELECT 
				"id", 
				"description", 
				"photos", 
				"created_at",
				"created_by",
				"updated_at",
				"updated_by",
				"deleted_at",
				"deleted_by"
			FROM "post" 
				WHERE 
				"deleted_at" IS NULL AND
			    "id"=$1`

	post := models.Post{}
	err = b.db.QueryRow(context.Background(), query, req.Id).Scan(
		&post.ID,
		&post.Description,
		&post.Photos,
		&created_at,
		&created_by,
		&updated_at,
		&updated_by,
		&deleted_at,
		&deleted_by,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	post.CreatedAt = created_at.Time.Format(time.RFC3339)
	post.CreatedBy = created_by.String
	if updated_at.Valid {
		post.UpdatedAt = updated_at.Time.Format(time.RFC3339)
	}

	if updated_by.Valid {
		post.UpdatedBy = updated_by.String
	}

	if deleted_at.Valid {
		post.DeletedAt = deleted_at.Time.Format(time.RFC3339)
	}

	if deleted_by.Valid {
		post.DeletedBy = deleted_by.String
	}

	return &post, nil
}

func (b *postRepo) GetAllActivePost(c context.Context, req *models.GetAllPostRequest) (*models.GetAllPost, error) {
	filter := ` WHERE deleted_at IS NULL `

	query := `
		SELECT 
			"id", 
			"description", 
			"photos", 
			"created_at",
			"created_by",
			"updated_at",
			"updated_by",
			"deleted_at",
			"deleted_by"
		FROM "post"
	`

	countQuery := `SELECT count(*) FROM post WHERE deleted_at IS NULL `

	if *req.Search != "" {
		filter += fmt.Sprintf(` AND description ILIKE  '%s' `, "%"+*req.Search+"%")
		countQuery += fmt.Sprintf(` AND description ILIKE '%s'`, *req.Search)
	}

	if *req.Page != 0 && *req.Limit != 0 {
		offset := (*req.Page - 1) * (*req.Limit)
		filter += fmt.Sprintf(" LIMIT %d OFFSET %d", *req.Limit, offset)
	}

	query += filter

	rows, err := b.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)

	for rows.Next() {
		var (
			created_at sql.NullTime
			created_by sql.NullString
			updated_at sql.NullTime
			updated_by sql.NullString
			deleted_at sql.NullTime
			deleted_by sql.NullString
		)
		post := models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Description,
			&post.Photos,
			&created_at,
			&created_by,
			&updated_at,
			&updated_by,
			&deleted_at,
			&deleted_by,
		)
		if err != nil {
			return nil, err
		}

		post.CreatedAt = created_at.Time.Format(time.RFC3339)
		post.CreatedBy = created_by.String
		if updated_at.Valid {
			post.UpdatedAt = updated_at.Time.Format(time.RFC3339)
		}

		if updated_by.Valid {
			post.UpdatedBy = updated_by.String
		}

		if deleted_at.Valid {
			post.DeletedAt = deleted_at.Time.Format(time.RFC3339)
		}

		if deleted_by.Valid {
			post.DeletedBy = deleted_by.String
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	count := 0
	err = b.db.QueryRow(c, countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}

	response := &models.GetAllPost{
		Posts: posts,
		Count: count,
	}

	return response, nil
}

func (b *postRepo) GetAllMyActivePost(c context.Context, req *models.GetAllMyPostRequest) (*models.GetAllPost, error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	filter := fmt.Sprintf(` WHERE deleted_at IS NULL AND created_by = '%s'`, userInfo.User_id)

	query := `
		SELECT 
			"id", 
			"description", 
			"photos", 
			"created_at",
			"created_by",
			"updated_at",
			"updated_by",
			"deleted_at",
			"deleted_by"
		FROM "post"
	`

	countQuery := fmt.Sprintf(`SELECT count(*) FROM post WHERE deleted_at IS NULL AND created_by = '%s'`, userInfo.User_id)

	if *req.Search != "" {
		filter += fmt.Sprintf(` AND description ILIKE  '%s' `, "%"+*req.Search+"%")
		countQuery += fmt.Sprintf(` AND description ILIKE '%s'`, "%"+*req.Search+"%")
	}

	if *req.Page != 0 && *req.Limit != 0 {
		offset := (*req.Page - 1) * (*req.Limit)
		filter += fmt.Sprintf(" LIMIT %d OFFSET %d", *req.Limit, offset)
	}

	query += filter

	rows, err := b.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)

	for rows.Next() {
		var (
			created_at sql.NullTime
			created_by sql.NullString
			updated_at sql.NullTime
			updated_by sql.NullString
			deleted_at sql.NullTime
			deleted_by sql.NullString
		)
		post := models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Description,
			&post.Photos,
			&created_at,
			&created_by,
			&updated_at,
			&updated_by,
			&deleted_at,
			&deleted_by,
		)
		if err != nil {
			return nil, err
		}

		post.CreatedAt = created_at.Time.Format(time.RFC3339)
		post.CreatedBy = created_by.String
		if updated_at.Valid {
			post.UpdatedAt = updated_at.Time.Format(time.RFC3339)
		}

		if updated_by.Valid {
			post.UpdatedBy = updated_by.String
		}

		if deleted_at.Valid {
			post.DeletedAt = deleted_at.Time.Format(time.RFC3339)
		}

		if deleted_by.Valid {
			post.DeletedBy = deleted_by.String
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	count := 0
	err = b.db.QueryRow(c, countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}

	response := &models.GetAllPost{
		Posts: posts,
		Count: count,
	}

	return response, nil
}

func (b *postRepo) UpdatePost(c context.Context, req *models.UpdatePost) (string, error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
			UPDATE "post" 
				SET 
				"description" = $1,
				"photos" = $2,
				"updated_at" = NOW(),
				"updated_by" = $3
				WHERE "id" = $4 AND "created_by" = $5 AND "deleted_at" IS NULL`

	result, err := b.db.Exec(
		context.Background(),
		query,
		req.Description,
		req.Photos,
		userInfo.User_id,
		req.ID,
		userInfo.User_id,
	)

	if err != nil {
		return "", fmt.Errorf("failed to update post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("post with ID %s not found", req.ID)
	}

	return req.ID, nil
}

func (b *postRepo) DeletePost(c context.Context, req *models.DeletePost) (resp string, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "post" 
		SET 
			"deleted_at" = NOW(),
			"deleted_by" = $1 

		WHERE 
			"deleted_at" IS  NULL AND
			"created_by" = $2 AND
			"id" = $3
`
	result, err := b.db.Exec(
		context.Background(),
		query,
		userInfo.User_id,
		req.Id,
		userInfo.User_id,
	)
	if err != nil {
		return "", fmt.Errorf("failed to update post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("post not found")
	}

	return req.Id, nil
}

func (b *postRepo) GetAllDeletedPost(c context.Context, req *models.GetAllPostRequest) (*models.GetAllPost, error) {
	filter := ` WHERE deleted_at IS NOT NULL`

	query := `
		SELECT 
			"id", 
			"description", 
			"photos", 
			"created_at",
			"created_by",
			"updated_at",
			"updated_by",
			"deleted_at",
			"deleted_by"
		FROM "post"
	`

	countQuery := `SELECT count(*) FROM post WHERE deleted_at IS NOT NULL `

	if *req.Search != "" {
		filter += fmt.Sprintf(` AND description ILIKE  '%s' `, "%"+*req.Search+"%")
		countQuery += fmt.Sprintf(` AND description ILIKE '%s'`, *req.Search)
	}

	if *req.Page != 0 && *req.Limit != 0 {
		offset := (*req.Page - 1) * (*req.Limit)
		filter += fmt.Sprintf(" LIMIT %d OFFSET %d", *req.Limit, offset)
	}

	query += filter

	rows, err := b.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)

	for rows.Next() {
		var (
			created_at sql.NullTime
			created_by sql.NullString
			updated_at sql.NullTime
			updated_by sql.NullString
			deleted_at sql.NullTime
			deleted_by sql.NullString
		)
		post := models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Description,
			&post.Photos,
			&created_at,
			&created_by,
			&updated_at,
			&updated_by,
			&deleted_at,
			&deleted_by,
		)
		if err != nil {
			return nil, err
		}

		post.CreatedAt = created_at.Time.Format(time.RFC3339)
		post.CreatedBy = created_by.String
		if updated_at.Valid {
			post.UpdatedAt = updated_at.Time.Format(time.RFC3339)
		}

		if updated_by.Valid {
			post.UpdatedBy = updated_by.String
		}

		if deleted_at.Valid {
			post.DeletedAt = deleted_at.Time.Format(time.RFC3339)
		}

		if deleted_by.Valid {
			post.DeletedBy = deleted_by.String
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	count := 0
	err = b.db.QueryRow(c, countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}

	response := &models.GetAllPost{
		Posts: posts,
		Count: count,
	}

	return response, nil
}
