package handler

import (
	"net/http"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"

	"github.com/gin-gonic/gin"
)

type ThemeHandler struct {
	service *service.ThemeService
}

func NewThemeHandler() *ThemeHandler {
	return &ThemeHandler{
		service: service.NewThemeService(),
	}
}

func (h *ThemeHandler) Get(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	theme, err := h.service.Get(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	c.JSON(http.StatusOK, theme)
}

func (h *ThemeHandler) Update(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var theme model.Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	theme.UserID = userID
	if err := h.service.Update(&theme); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, theme)
}
