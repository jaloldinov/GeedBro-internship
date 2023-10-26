package api

import (
	"auth/pkg/helper"

	"auth/api/handler"

	"github.com/gin-gonic/gin"
)

func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	// authentication sign up and login
	r.POST("/auth/login", h.Login)
	r.POST("/auth/sign-up", h.SignUp)

	// user routes
	r.POST("/user", helper.AuthMiddleWare, h.CreateUser)
	r.GET("/user/:id", helper.AuthMiddleWare, h.GetUser)
	r.GET("/user", helper.AuthMiddleWare, h.GetAllUser)
	r.PUT("/user/:id", helper.AuthMiddleWare, h.UpdateUser)
	r.DELETE("/user/:id", helper.AuthMiddleWare, h.DeleteUser)

	// delted users and posts
	r.GET("/deleted-users", helper.AuthMiddleWare, h.GetAllDeletedUser)
	r.GET("/deleted-posts", helper.AuthMiddleWare, h.GetAllDeletedPost)

	// posts
	r.POST("/post", helper.AuthMiddleWare, h.CreatePost)
	r.GET("/post/:post_id", helper.AuthMiddleWare, h.GetPost)
	r.GET("/posts/all", h.GetAllPost)
	r.PUT("/post/:post_id", helper.AuthMiddleWare, h.UpdatePost)
	r.DELETE("/post/:post_id", helper.AuthMiddleWare, h.DeletePost)

	r.GET("/my/posts", helper.AuthMiddleWare, h.GetAllMyPost)

	// post_likes
	r.POST("/like", helper.AuthMiddleWare, h.CreateLike)
	r.GET("/like-count/:post_id", h.GetLike)
	r.DELETE("/like", helper.AuthMiddleWare, h.DeleteLike)

	// post comment section
	r.POST("/comment/:post_id", helper.AuthMiddleWare, h.CreateComment)
	r.GET("/my/comments", helper.AuthMiddleWare, h.GetMyComments)
	r.GET("/post/comment/by/post/:post_id", h.GetPostComments)
	r.PUT("/comment", helper.AuthMiddleWare, h.UpdateComment)
	r.DELETE("/comment/:id", helper.AuthMiddleWare, h.DeleteComment)
	r.DELETE("/my/comment/delete/:id", helper.AuthMiddleWare, h.DeleteMyPostComment)

	// comment likes
	r.POST("/comment-like", helper.AuthMiddleWare, h.CreateCommentLike)
	r.GET("/comment-like/:comment_id", h.GetCommentLikes)
	r.DELETE("/comment-like", helper.AuthMiddleWare, h.DeleteCommentLike)

	return r
}
