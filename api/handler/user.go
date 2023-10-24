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
