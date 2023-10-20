package handler

import (
	"auth/models"
	"auth/pkg/helper"
	"auth/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Router       /user [POST]
// @Summary      CREATES USER
// @Description  creates a new user based on the given username amd password
// @Tags         USER
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateUser  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.CreateUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	hashedPass, err := helper.GeneratePasswordHash(user.Password)
	if err != nil {
		h.log.Error("error while generating hash password:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	user.Password = string(hashedPass)

	resp, err := h.storage.User().CreateUser(c.Request.Context(), &user)
	if err != nil {
		fmt.Println("error User Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "username is already used, enter another one")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created", "id": resp})
}

// GetUser godoc
// @Router       /user/{id} [GET]
// @Summary      GET USER BY ID
// @Description  gets the user by ID
// @Tags         USER
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID" format(uuid)
// @Success      200  {object}  models.User
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.User().GetUser(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Println("error User Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListUsers godoc
// @Router       /user [GET]
// @Summary      GET  ALL USERS
// @Description  gets all user based on limit, page and search by username
// @Tags         USER
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllUser
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllUser(c *gin.Context) {
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

	resp, err := h.storage.User().GetAllActiveUser(c.Request.Context(), &models.GetAllUserRequest{
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

// UpdateUser godoc
// @Router       /user/{id} [PUT]
// @Summary      UPDATES USER BY ID
// @Description  UPDATES USER BASED ON GIVEN DATA AND ID
// @Tags         USER
// @Accept       json
// @Produce      json
// @Param        data  body      models.UpdateUser  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.UpdateUser
	err := c.ShouldBind(&user)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.User().UpdateUser(c.Request.Context(), &user)
	if err != nil {
		h.log.Error("error User Update:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "updated user id": resp})
}

// DeleteUser godoc
// @Router       /user/{id} [DELETE]
// @Summary      DELETE USER BY ID
// @Description  DELETES USER BASED ON ID
// @Tags         USER
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of user" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.User().DeleteUser(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting user:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "deleted user id": resp})
}

// ListUsers godoc
// @Router       /deleted-users [GET]
// @Summary      GETS ALL DELETED USERS
// @Description  gets all user based on limit, page and search by username
// @Tags         USER
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllUser
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllDeletedUser(c *gin.Context) {
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

	resp, err := h.storage.User().GetAllDeletedUser(c.Request.Context(), &models.GetAllUserRequest{
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
