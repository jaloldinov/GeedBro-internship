package handler

import (
	"auth/models"
	"auth/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateFile(c *gin.Context) {
	var file models.CreateFile
	err := c.ShouldBind(&file)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, er := h.storage.File().CreateFile(c, &file)
	if er != nil {
		fmt.Println("error Post Create:", er)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created", "resp": resp})
}

func (h *Handler) CreateFiles(c *gin.Context) {
	var files models.CreateFiles
	err := c.ShouldBind(&files)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, er := h.storage.File().CreateFiles(c, &files)

	if er != nil {
		fmt.Println("error Post Create:", er)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created", "resp": resp})
}
