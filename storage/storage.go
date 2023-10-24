package storage

import (
	"auth/api/response"
	"auth/models"
	"context"
)

type StorageI interface {
	User() UsersI
	Post() PostsI
	Like() LikesI
	File() UserFileUploadI
	Comment() PostCommentsI
	CommentLike() CommentLikeI
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
	GetAllMyActivePost(context.Context, *models.GetAllMyPostRequest) (*models.GetAllPost, error)
}

type LikesI interface {
	AddLike(context.Context, *models.CreateLike) error
	DeleteLike(context.Context, *models.DeleteLike) (string, error)
	GetLikesCount(context.Context, string) (int, error)
}

type UserFileUploadI interface {
	CreateFile(context.Context, *models.CreateFile) (*models.CreateFileResponse, *response.ErrorResp)
	CreateFiles(context.Context, *models.CreateFiles) (*models.CreateFilesResponse, *response.ErrorResp)
}

type PostCommentsI interface {
	CreateComment(context.Context, *models.CreateComment) (string, error)
	GetComment(context.Context, *models.IdRequest) (*models.Comment, error)
	GetPostComments(context.Context, *models.GetAllPostComments) (*models.GetAllCommentResponse, error)
	UpdateComment(context.Context, *models.UpdateComment) (string, error)
	DeleteComment(context.Context, *models.DeleteComment) (string, error)
}

type CommentLikeI interface {
	AddLike(context.Context, *models.CreateCommentLike) error
	DeleteLike(context.Context, *models.DeleteCommentLike) (string, error)
	GetLikesCount(context.Context, string) (int, error)
}
