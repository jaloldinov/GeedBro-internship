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

type likeRepo struct {
	db *pgxpool.Pool
}

func NewLikeRepo(db *pgxpool.Pool) *likeRepo {
	return &likeRepo{
		db: db,
	}
}

func (b *likeRepo) AddLike(c context.Context, req *models.CreateLike) error {
	// Check if the user has already liked the post
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
		SELECT COUNT(id) FROM post_likes WHERE post_id = $1 AND user_id = $2
	`

	count := 0
	err := b.db.QueryRow(c, query, req.PostId, userInfo.User_id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing like: %w", err)
	}

	if count > 0 {
		// User has already liked the post, do not add another like
		return nil
	}

	// User has not liked the post, add the new like
	insertQuery := `
		INSERT INTO "post_likes" ("id", "post_id", "user_id", "created_at")
		VALUES ($1, $2, $3, NOW())
	`

	_, err = b.db.Exec(c, insertQuery, uuid.NewString(), req.PostId, userInfo.User_id)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}

	return nil
}

func (b *likeRepo) DeleteLike(c context.Context, req *models.DeleteLike) (resp string, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "post_likes" 
		SET 
			"deleted_at" = NOW()
		WHERE 
		"deleted_at" IS NULL AND
		 "user_id" = $1 AND "post_id" = $2
`
	result, err := b.db.Exec(
		context.Background(),
		query,
		userInfo.User_id,
		req.PostId,
	)
	if err != nil {
		return "", fmt.Errorf("failed to delete like: %w", err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("like not found")
	}

	return "like removed", nil
}

func (b *likeRepo) GetLikesCount(c context.Context, req string) (int, error) {
	query := `
		SELECT COUNT(id) AS like_count
		FROM post_likes
		WHERE deleted_at IS NULL AND post_id = $1
	`

	count := 0
	err := b.db.QueryRow(c, query, req).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	fmt.Println(count)

	return count, nil
}
