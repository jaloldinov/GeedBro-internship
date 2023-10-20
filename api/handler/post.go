package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePost godoc
// @Router       /post [POST]
// @Summary      CREATES POST
// @Description  creates a new post based on the given postname amd password
// @Tags         POST
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreatePost  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreatePost(c *gin.Context) {
	var post models.CreatePost
	err := c.ShouldBindJSON(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Post().CreatePost(c.Request.Context(), &post)
	if err != nil {
		fmt.Println("error Post Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created", "id": resp})
}

// GetPost godoc
// @Router       /post/{id} [GET]
// @Summary      GET POST BY ID
// @Description  gets the post by ID
// @Tags         POST
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID" format(uuid)
// @Success      200  {object}  models.Post
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetPost(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Post().GetPost(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Println("error Post Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPosts godoc
// @Router       /post [GET]
// @Summary      GET  ALL POSTS
// @Description  gets all post based on limit, page and search by postname
// @Tags         POST
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllPost
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
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

	resp, err := h.storage.Post().GetAllActivePost(c.Request.Context(), &models.GetAllPostRequest{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.log.Error("error:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdatePost godoc
// @Router       /post/{id} [PUT]
// @Summary      UPDATES POST BY ID
// @Description  UPDATES POST BASED ON GIVEN DATA AND ID
// @Tags         POST
// @Accept       json
// @Produce      json
// @Param        data  body      models.UpdatePost  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
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

// DeletePost godoc
// @Router       /post [DELETE]
// @Summary      DELETE POST BY ID
// @Description  DELETES POST BASED ON ID
// @Tags         POST
// @Accept       json
// @Produce      json
// @Param        data  body      models.DeletePost  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeletePost(c *gin.Context) {
	var post models.DeletePost
	err := c.ShouldBind(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.Post().DeletePost(c.Request.Context(), &models.DeletePost{Id: post.Id, UserId: post.UserId})
	if err != nil {
		h.log.Error("error deleting post:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "deleted post id": resp})
}

// ListPosts godoc
// @Router       /deleted-posts [GET]
// @Summary      GETS ALL DELETED POSTS
// @Description  gets all post based on limit, page and search by postname
// @Tags         POST
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllPost
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
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

	resp, err := h.storage.Post().GetAllDeletedPost(c.Request.Context(), &models.GetAllPostRequest{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.log.Error("error:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)
}
