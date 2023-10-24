package api

import (
	_ "auth/api/docs"
	"auth/pkg/helper"

	"auth/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	r.POST("/post", helper.AuthMiddleWare, h.CreatePost)
	r.GET("/post/:id", helper.AuthMiddleWare, h.GetPost)
	r.GET("/post", helper.AuthMiddleWare, h.GetAllPost)
	r.PUT("/post/:id", helper.AuthMiddleWare, h.UpdatePost)
	r.DELETE("/post", helper.AuthMiddleWare, h.DeletePost)

	r.GET("/my/post", helper.AuthMiddleWare, h.GetAllMyPost)

	r.POST("/like", helper.AuthMiddleWare, h.CreateLike)
	r.GET("/like-count/:post_id", h.GetLike)
	r.DELETE("/like", helper.AuthMiddleWare, h.DeleteLike)

	// file uploading
	r.POST("/file/upload", helper.AuthMiddleWare, h.CreateFile)
	r.POST("/files/upload", helper.AuthMiddleWare, h.CreateFiles)

	// post comment section
	r.POST("/comment", helper.AuthMiddleWare, h.CreateComment)
	r.GET("/comment/:id", helper.AuthMiddleWare, h.GetComment)
	r.GET("/comment", h.GetPostComments)
	r.PUT("/comment", helper.AuthMiddleWare, h.UpdateComment)
	r.DELETE("/comment/:id", helper.AuthMiddleWare, h.DeleteComment)

	// comment likes
	r.POST("/comment-like", helper.AuthMiddleWare, h.CreateCommentLike)
	r.GET("/comment-like/:comment_id", h.GetCommentLikes)
	r.DELETE("/comment-like", helper.AuthMiddleWare, h.DeleteCommentLike)

	// Serve Swagger API documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
