package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// create comment to post
func (h *Handler) CreateComment(ctx *gin.Context) {
	var comment models.CreateComment

	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Comment().CreateComment(ctx, &comment)
	if err != nil {
		fmt.Println("error comment create:", err.Error())
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "created", "id": resp})
}

// get one comment from post
func (h *Handler) GetComment(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Comment().GetComment(c, &models.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Println("error comment get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// list of post comments
func (h *Handler) GetPostComments(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		h.log.Error("error get page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error get limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}

	search := c.Query("post_id")

	resp, err := h.storage.Comment().GetPostComments(c, &models.GetAllPostComments{
		Page:   &page,
		Limit:  &limit,
		PostId: &search,
	})
	if err != nil {
		h.log.Error("error:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// update post comment
func (h *Handler) UpdateComment(c *gin.Context) {
	var comment models.UpdateComment
	err := c.ShouldBind(&comment)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.Comment().UpdateComment(c, &comment)
	if err != nil {
		h.log.Error("Failed to update comment", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp})
}

// deletes comment
func (h *Handler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Comment().DeleteComment(c, &models.DeleteComment{Id: id})
	if err != nil {
		h.log.Error("error deleting comment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp})
}
