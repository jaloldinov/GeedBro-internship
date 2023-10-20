package api

import (
	_ "auth/api/docs"

	"auth/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("/user", h.CreateUser)
	r.GET("/user/:id", h.GetUser)
	r.GET("/user", h.GetAllUser)
	r.PUT("/user/:id", h.UpdateUser)
	r.DELETE("/user/:id", h.DeleteUser)

	r.GET("/deleted-users", h.GetAllDeletedUser)

	r.POST("/post", h.CreatePost)
	r.GET("/post/:id", h.GetPost)
	r.GET("/post", h.GetAllPost)
	r.PUT("/post/:id", h.UpdatePost)
	r.DELETE("/post", h.DeletePost)

	r.GET("/deleted-posts", h.GetAllDeletedPost)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
