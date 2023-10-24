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

type commentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *commentRepo {
	return &commentRepo{
		db: db,
	}
}

// creates comment on one post
func (r *commentRepo) CreateComment(ctx context.Context, req *models.CreateComment) (string, error) {
	userInfo := ctx.Value("user_info").(helper.TokenInfo)
	id := uuid.NewString()

	query := `
		INSERT INTO "post_comments" (
			"id",
			"post_id",
			"comment",
			"created_by",
			"created_at"
		)
	VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.PostId,
		req.Comment,
		userInfo.User_id,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create comment: %v", err)
	}

	return id, nil
}

// get one comment based on given comment_id
func (b *commentRepo) GetComment(c context.Context, req *models.IdRequest) (resp *models.Comment, err error) {

	var (
		created_at sql.NullTime
		created_by sql.NullString
		updated_at sql.NullTime
		updated_by sql.NullString
	)

	query := `
			SELECT 
				"id", 
				"post_id", 
				"comment", 
				"created_at",
				"created_by",
				"updated_at",
				"updated_by"
			FROM "post_comments" 
				WHERE 
				"deleted_at" IS NULL AND
			    "id"=$1`

	comment := models.Comment{}
	err = b.db.QueryRow(context.Background(), query, req.Id).Scan(
		&comment.ID,
		&comment.PostId,
		&comment.Comment,
		&created_at,
		&created_by,
		&updated_at,
		&updated_by,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	comment.CreatedAt = created_at.Time.Format(time.RFC3339)
	comment.CreatedBy = created_by.String
	if updated_at.Valid {
		comment.UpdatedAt = updated_at.Time.Format(time.RFC3339)
	}

	if updated_by.Valid {
		comment.UpdatedBy = updated_by.String
	}

	return &comment, nil
}

// get all post comments
func (b *commentRepo) GetPostComments(c context.Context, req *models.GetAllPostComments) (*models.GetAllCommentResponse, error) {

	filter := fmt.Sprintf(`  WHERE "deleted_at" IS NULL AND "post_id" = '%s'`, *req.PostId)

	query := `
		SELECT 
				"id", 
				"post_id", 
				"comment", 
				"created_at",
				"created_by",
				"updated_at",
				"updated_by"
		FROM "post_comments"
	`

	countQuery := fmt.Sprintf(`SELECT count(*) FROM "post_comments" WHERE "deleted_at" IS NULL AND "post_id" = '%s'`, *req.PostId)

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

	comments := make([]models.Comment, 0)

	for rows.Next() {
		var (
			created_at sql.NullTime
			created_by sql.NullString
			updated_at sql.NullTime
			updated_by sql.NullString
		)
		comment := models.Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostId,
			&comment.Comment,
			&created_at,
			&created_by,
			&updated_at,
			&updated_by,
		)
		if err != nil {
			return nil, err
		}

		comment.CreatedAt = created_at.Time.Format(time.RFC3339)
		comment.CreatedBy = created_by.String
		if updated_at.Valid {
			comment.UpdatedAt = updated_at.Time.Format(time.RFC3339)
		}

		if updated_by.Valid {
			comment.UpdatedBy = updated_by.String
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	count := 0
	err = b.db.QueryRow(c, countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}

	response := &models.GetAllCommentResponse{
		Comments: comments,
		Count:    count,
	}

	return response, nil
}

func (b *commentRepo) UpdateComment(c context.Context, req *models.UpdateComment) (string, error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
			UPDATE "post_comments" 
				SET 
				"comment" = $1,
				"updated_at" = NOW(),
				"updated_by" = $2
				WHERE "id" = $3 AND "created_by" = $4 AND "deleted_at" IS NULL`

	result, err := b.db.Exec(
		context.Background(),
		query,
		req.Comment,
		userInfo.User_id,
		req.ID,
		userInfo.User_id,
	)

	if err != nil {
		return "", fmt.Errorf("you can't edit this comment!: %w", err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("post with ID %s not found", req.ID)
	}

	return "updated", nil
}

func (b *commentRepo) DeleteComment(c context.Context, req *models.DeleteComment) (resp string, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "post_comments" 
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
		userInfo.User_id,
		req.Id,
	)
	if err != nil {
		return "", fmt.Errorf("failed to delete(update) comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("comment not found")
	}

	return "deleted", nil
}
