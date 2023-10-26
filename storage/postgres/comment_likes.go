package postgres

import (
	"auth/models"
	"auth/pkg/helper"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type commentLikeRepo struct {
	db *pgxpool.Pool
}

func NewCommentLikeRepo(db *pgxpool.Pool) *commentLikeRepo {
	return &commentLikeRepo{
		db: db,
	}
}

func (b *commentLikeRepo) AddLike(c context.Context, req *models.CreateCommentLike) error {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	// Check if the user has already liked the post
	query := `
		SELECT COUNT(id) FROM comment_likes WHERE comment_id = $1 AND user_id = $2
	`

	count := 0
	err := b.db.QueryRow(c, query, req.CommentId, userInfo.User_id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing like: %w", err)
	}

	if count > 0 {
		// User has already liked the post, do not add another like
		return nil
	}

	// User has not liked the post, add the new like
	insertQuery := `
		INSERT INTO "comment_likes" ("id", "comment_id", "user_id", "created_at")
		VALUES ($1, $2, $3, NOW())
	`

	_, err = b.db.Exec(c, insertQuery, uuid.NewString(), req.CommentId, userInfo.User_id)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}

	return nil
}

func (b *commentLikeRepo) DeleteLike(c context.Context, req *models.DeleteCommentLike) (resp string, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
		UPDATE "comment_likes" 
		SET 
			"deleted_at" = NOW()
		WHERE "deleted_at" IS NULL AND
			"user_id" = $1 AND "comment_id" = $2
		RETURNING "id"
	`

	var deletedID string
	err = b.db.QueryRow(c, query, userInfo.User_id, req.CommentId).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("you can't modify the like")
		}
		return "", fmt.Errorf("failed to delete like: %w", err)
	}

	return "like removed", nil
}

func (b *commentLikeRepo) GetLikesCount(c context.Context, req string) (int, error) {
	query := `
		SELECT COUNT(id) AS like_count
		FROM comment_likes
		WHERE "deleted_at" IS NULL AND  comment_id = $1
	`

	count := 0
	err := b.db.QueryRow(c, query, req).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	return count, nil
}
