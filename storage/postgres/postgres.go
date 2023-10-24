package postgres

import (
	"auth/config"
	"auth/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db           *pgxpool.Pool
	users        *userRepo
	posts        *postRepo
	likes        *likeRepo
	files        *fileRepo
	comments     *commentRepo
	commentLikes *commentLikeRepo
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = cfg.PostgresMaxConnections

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (b *store) User() storage.UsersI {
	if b.users == nil {
		b.users = NewUserRepo(b.db)
	}
	return b.users
}

func (b *store) Post() storage.PostsI {
	if b.posts == nil {
		b.posts = NewPostRepo(b.db)
	}
	return b.posts
}

func (b *store) Like() storage.LikesI {
	if b.likes == nil {
		b.likes = NewLikeRepo(b.db)
	}
	return b.likes
}

func (b *store) File() storage.UserFileUploadI {
	if b.files == nil {
		b.files = NewFileRepo(b.db)
	}
	return b.files
}

func (b *store) Comment() storage.PostCommentsI {
	if b.comments == nil {
		b.comments = NewCommentRepo(b.db)
	}
	return b.comments
}

func (b *store) CommentLike() storage.CommentLikeI {
	if b.commentLikes == nil {
		b.commentLikes = NewCommentLikeRepo(b.db)
	}
	return b.commentLikes
}
