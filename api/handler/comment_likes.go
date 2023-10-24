package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCommentLike(c *gin.Context) {

	var like models.CreateCommentLike

	err := c.ShouldBind(&like)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	err = h.storage.CommentLike().AddLike(c, &like)
	if err != nil {
		fmt.Println("error Like Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "like added"})
}

func (h *Handler) GetCommentLikes(c *gin.Context) {

	comment_id := c.Param("comment_id")

	count, err := h.storage.CommentLike().GetLikesCount(c, comment_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, logger.Error(err))
		fmt.Println("error comment like count:", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *Handler) DeleteCommentLike(c *gin.Context) {
	var comment_id models.DeleteCommentLike

	err := c.ShouldBindJSON(&comment_id)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println(comment_id)
	resp, err := h.storage.CommentLike().DeleteLike(c, &models.DeleteCommentLike{CommentId: comment_id.CommentId})
	if err != nil {
		h.log.Error("error deleting like:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "res": resp})
}
