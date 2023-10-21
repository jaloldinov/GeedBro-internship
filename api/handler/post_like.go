package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateLike godoc
// @Security ApiKeyAuth
// @Router       /like [POST]
// @Summary      CREATES LIKE
// @Description  addes like to existing like
// @Tags         LIKE
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateLike  true  "like data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateLike(c *gin.Context) {
	var like models.CreateLike
	err := c.ShouldBindJSON(&like)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	err = h.storage.Like().AddLike(c.Request.Context(), &like)
	if err != nil {
		fmt.Println("error Like Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "like added"})
}

// GetLike godoc
// @Security ApiKeyAuth
// @Router       /like-count/{post_id} [GET]
// @Summary      GET LIKE BY ID
// @Description  gets likes count based on post id
// @Tags         LIKE
// @Accept       json
// @Produce      json
// @Param        post_id    path     string  true  "post_id" format(uuid)
// @Success      200  {object}  int
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetLike(c *gin.Context) {
	postId := c.Param("post_id")
	fmt.Println(postId)

	resp, err := h.storage.Like().GetLikesCount(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, logger.Error(err))
		fmt.Println("error Like Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteLike godoc
// @Security ApiKeyAuth
// @Router       /like [DELETE]
// @Summary      DELETE LIKE
// @Description  DELETES LIKE BASED ON user_id and post_id
// @Tags         LIKE
// @Accept       json
// @Produce      json
// @Param        data  body      models.DeleteLike  true  "like data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteLike(c *gin.Context) {
	var like models.DeleteLike
	err := c.ShouldBind(&like)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.Like().DeleteLike(c.Request.Context(), &models.DeleteLike{UserId: like.UserId, PostId: like.PostId})
	if err != nil {
		h.log.Error("error deleting like:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "res": resp})
}
