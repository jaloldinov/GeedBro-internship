package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateLike(c *gin.Context) {
	var like models.CreateLike
	err := c.ShouldBindJSON(&like)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	err = h.storage.Like().AddLike(c, &like)
	if err != nil {
		fmt.Println("error Like Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "like added"})
}

func (h *Handler) GetLike(c *gin.Context) {
	postId := c.Param("post_id")
	fmt.Println(postId)

	resp, err := h.storage.Like().GetLikesCount(c, postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, logger.Error(err))
		fmt.Println("error Like Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteLike(c *gin.Context) {
	var like models.DeleteLike
	err := c.ShouldBind(&like)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.Like().DeleteLike(c, &models.DeleteLike{PostId: like.PostId})
	if err != nil {
		h.log.Error("error deleting like:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "res": resp})
}
