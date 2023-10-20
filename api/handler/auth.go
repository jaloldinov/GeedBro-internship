package handler

import (
	"auth/config"
	"auth/models"
	"auth/pkg/helper"
	"auth/pkg/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignUp godoc
// @Router       /auth/sign-up [POST]
// @Summary      CREATES USER
// @Description  CREATES USER BASED ON GIVEN DATA
// @Tags         AUTH
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateUser  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) SignUp(c *gin.Context) {
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

	if err != nil {
		id, err := h.storage.User().CreateUser(c.Request.Context(), &user)
		if err != nil {
			fmt.Println("error User Create:", err.Error())
			c.JSON(http.StatusInternalServerError, "login or phone number is already used, enter another one")
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "created", "id": id})
		return
	}
	c.JSON(http.StatusBadGateway, gin.H{"error": "username is used, enter another one"})
}

// loginUser godoc
// @Router       /auth/login [POST]
// @Summary      auth
// @Description  login
// @Tags         AUTH
// @Accept       json
// @Produce      json
// @Param        user    body   models.LoginRequest  true  "data of user"
// @Success      200  {object}  models.LoginRespond
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fields in body"})
		return
	}

	resp, err := h.storage.User().GetByUsername(context.Background(), &models.LoginRequest{
		Username: req.Username,
	})
	if err != nil {
		h.log.Error("error get by username:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not found username"})
		return
	}

	// Compare hashed password with plain text password
	err = helper.ComparePasswords([]byte(resp.Password), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "login or password didn't match"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password comparison failed"})
		}
		return
	}

	m := make(map[string]interface{})
	m["username"] = resp.Username

	token, _ := helper.GenerateJWT(m, config.TokenExpireTime, config.JWTSecretKey)
	c.JSON(http.StatusOK, models.LoginRespond{Token: token})
}
