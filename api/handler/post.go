package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var post models.CreatePost
	err := c.ShouldBindJSON(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Post().CreatePost(c, &post)
	if err != nil {
		fmt.Println("error Post Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created", "id": resp})
}

func (h *Handler) GetPost(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Post().GetPost(c, &models.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Println("error Post Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAllMyPost(c *gin.Context) {

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

	search := c.Query("search")

	resp, err := h.storage.Post().GetAllMyActivePost(c, &models.GetAllMyPostRequest{
		Page:   &page,
		Limit:  &limit,
		Search: &search,
	})
	if err != nil {
		h.log.Error("error:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAllPost(c *gin.Context) {
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

	username := c.Query("search")
	resp, err := h.storage.Post().GetAllActivePost(c.Request.Context(), &models.GetAllPostRequest{
		Page:   &page,
		Limit:  &limit,
		Search: &username,
	})
	if err != nil {
		h.log.Error("error:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdatePost(c *gin.Context) {
	var post models.UpdatePost
	err := c.ShouldBind(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.Post().UpdatePost(c.Request.Context(), &post)
	if err != nil {
		h.log.Error("error Post Update:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "updated post id": resp})
}

func (h *Handler) DeletePost(c *gin.Context) {
	var post models.DeletePost
	err := c.ShouldBindJSON(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.Post().DeletePost(c, &models.DeletePost{Id: post.Id})
	if err != nil {
		h.log.Error("error deleting post:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "deleted post id": resp})
}

func (h *Handler) GetAllDeletedPost(c *gin.Context) {
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

	username := c.Query("search")

	resp, err := h.storage.Post().GetAllDeletedPost(c.Request.Context(), &models.GetAllPostRequest{
		Page:   &page,
		Limit:  &limit,
		Search: &username,
	})

	if err != nil {
		h.log.Error("error:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)
}
