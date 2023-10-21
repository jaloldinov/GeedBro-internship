package storage

import (
	"auth/models"
	"context"
)

type StorageI interface {
	User() UsersI
	Post() PostsI
	Like() LikesI
}

type UsersI interface {
	CreateUser(context.Context, *models.CreateUser) (string, error)
	GetUser(context.Context, *models.IdRequest) (*models.User, error)
	GetAllActiveUser(context.Context, *models.GetAllUserRequest) (*models.GetAllUser, error)
	UpdateUser(context.Context, *models.UpdateUser) (string, error)
	DeleteUser(context.Context, *models.IdRequest) (string, error)

	GetAllDeletedUser(context.Context, *models.GetAllUserRequest) (*models.GetAllUser, error)
	GetByUsername(context.Context, *models.LoginRequest) (*models.LoginDataRespond, error)
}

type PostsI interface {
	CreatePost(context.Context, *models.CreatePost) (string, error)
	GetPost(context.Context, *models.IdRequest) (*models.Post, error)
	GetAllActivePost(context.Context, *models.GetAllPostRequest) (*models.GetAllPost, error)
	UpdatePost(context.Context, *models.UpdatePost) (string, error)
	DeletePost(context.Context, *models.DeletePost) (string, error)

	GetAllDeletedPost(context.Context, *models.GetAllPostRequest) (*models.GetAllPost, error)
}

type LikesI interface {
	AddLike(context.Context, *models.CreateLike) error
	DeleteLike(context.Context, *models.DeleteLike) (string, error)
	GetLikesCount(context.Context, string) (int, error)
}
