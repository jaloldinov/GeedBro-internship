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

	r.POST("/auth/login", h.Login)
	r.POST("/auth/sign-up", h.SignUp)

	r.POST("/user", helper.AuthMiddleWare, h.CreateUser)
	r.GET("/user/:id", helper.AuthMiddleWare, h.GetUser)
	r.GET("/user", helper.AuthMiddleWare, h.GetAllUser)
	r.PUT("/user/:id", helper.AuthMiddleWare, h.UpdateUser)
	r.DELETE("/user/:id", helper.AuthMiddleWare, h.DeleteUser)

	r.GET("/deleted-users", helper.AuthMiddleWare, h.GetAllDeletedUser)

	r.POST("/post", helper.AuthMiddleWare, h.CreatePost)
	r.GET("/post/:id", helper.AuthMiddleWare, h.GetPost)
	r.GET("/post", helper.AuthMiddleWare, h.GetAllPost)
	r.PUT("/post/:id", helper.AuthMiddleWare, h.UpdatePost)
	r.DELETE("/post", helper.AuthMiddleWare, h.DeletePost)

	r.GET("/deleted-posts", helper.AuthMiddleWare, h.GetAllDeletedPost)

	// Serve Swagger API documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
