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

func (b *commentRepo) GetMyComments(c context.Context) (resp *models.GetAllCommentResponse, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
    SELECT
        pc."id",
        pc."post_id",
        pc."comment",
        (SELECT COUNT(*)
            FROM "comment_likes" cl
            WHERE cl."comment_id" = pc."id"
        ) AS "likes_count",
        pc."created_at"
    FROM "post_comments" pc
    WHERE
        pc."deleted_at" IS NULL
        AND pc."created_by" = $1
	GROUP BY pc.id ORDER BY created_at desc
`
	countQuery := `SELECT COUNT(1) AS counts FROM "post_comments" WHERE "deleted_at" IS NULL AND "created_by" = $1`

	rows, err := b.db.Query(c, query, userInfo.User_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var created_at sql.NullTime

		comment := models.Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostId,
			&comment.Comment,
			&comment.LikeCount,
			&created_at,
		)
		if err != nil {
			return nil, err
		}

		comment.CreatedAt = created_at.Time.Format(time.RFC3339)

		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	count := 0
	err = b.db.QueryRow(c, countQuery, userInfo.User_id).Scan(&count)
	if err != nil {
		return nil, err
	}

	response := &models.GetAllCommentResponse{
		Comments: comments,
		Count:    count,
	}

	return response, nil
}

// get all post comments
func (b *commentRepo) GetPostComments(c context.Context, req *models.GetAllPostComments) (*models.GetAllCommentResponse, error) {

	filter := fmt.Sprintf(`  WHERE "deleted_at" IS NULL AND "post_id" = '%s'`, *req.PostId)

	query := `
		SELECT 
				"id", 
				"post_id", 
				"comment", 
				(SELECT COUNT(*) 
				FROM "comment_likes"
				WHERE "deleted_at" IS  NULL
				AND "comment_id" = post_comments."id"
			) AS "likes_count",
				"created_at"
		FROM "post_comments"
	`

	countQuery := fmt.Sprintf(`SELECT count(*) FROM "post_comments" WHERE "deleted_at" IS NULL AND "post_id" = '%s'`, *req.PostId)

	if *req.Page != 0 && *req.Limit != 0 {
		offset := (*req.Page - 1) * (*req.Limit)
		filter += fmt.Sprintf(" ORDER BY created_at desc LIMIT %d OFFSET %d ", *req.Limit, offset)
	}

	query += filter

	rows, err := b.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)

	for rows.Next() {
		var created_at sql.NullTime

		comment := models.Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostId,
			&comment.Comment,
			&comment.LikeCount,
			&created_at,
		)
		if err != nil {
			return nil, err
		}

		comment.CreatedAt = created_at.Time.Format(time.RFC3339)

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

func (b *commentRepo) DeleteMyPostComment(c context.Context, req *models.DeleteComment) (resp string, err error) {
	userInfo := c.Value("user_info").(helper.TokenInfo)

	postIdQuery := `SELECT post_id FROM post_comments WHERE id = $1`

	var post_id string
	err = b.db.QueryRow(c, postIdQuery, req.Id).Scan(
		&post_id,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("comment not found")
		}
		return "", fmt.Errorf("failed to get post id: %w", err)
	}

	postOwnerQuery := `SELECT created_by FROM post WHERE id = $1`
	var post_id_owner string
	err = b.db.QueryRow(c, postOwnerQuery, post_id).Scan(
		&post_id_owner,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("post not found")
		}
		return "", fmt.Errorf("failed to get post id: %w", err)
	}

	trueR := userInfo.User_id == post_id_owner

	query := fmt.Sprintf(`UPDATE "post_comments"
							SET
								"deleted_at" = NOW(),
								"deleted_by" = $1

							WHERE
								"deleted_at" IS  NULL AND %v AND
								"id" = $2`, trueR)

	result, err := b.db.Exec(
		context.Background(),
		query,
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
