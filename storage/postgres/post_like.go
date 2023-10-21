package postgres

import (
	"auth/models"
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
	query := `
		SELECT COUNT(*) FROM post_likes WHERE post_id = $1 AND user_id = $2
	`

	count := 0
	err := b.db.QueryRow(c, query, req.PostId, req.UserId).Scan(&count)
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

	_, err = b.db.Exec(c, insertQuery, uuid.NewString(), req.PostId, req.UserId)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}

	return nil
}

func (b *likeRepo) DeleteLike(c context.Context, req *models.DeleteLike) (resp string, err error) {
	query := `
	 	UPDATE "post_likes" 
		SET 
			"deleted_at" = NOW()
		WHERE 
		"deleted_at" IS  NULL AND
		 "user_id" = $1 AND "post_id" = $2
`
	result, err := b.db.Exec(
		context.Background(),
		query,
		req.UserId,
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
		SELECT COUNT(*) AS like_count
		FROM post_likes
		WHERE post_id = $1
	`

	fmt.Println(req)
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