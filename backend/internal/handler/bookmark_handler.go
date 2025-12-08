package handler

import (
	"net/http"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	service *service.BookmarkService
}

func NewBookmarkHandler() *BookmarkHandler {
	return &BookmarkHandler{
		service: service.NewBookmarkService(),
	}
}

func (h *BookmarkHandler) GetAll(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	bookmarks, err := h.service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmarks)
}

func (h *BookmarkHandler) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	bookmark, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (h *BookmarkHandler) Create(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var bookmark model.Bookmark
	if err := c.ShouldBindJSON(&bookmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookmark.UserID = userID
	if err := h.service.Create(&bookmark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bookmark)
}

func (h *BookmarkHandler) Update(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var bookmark model.Bookmark
	if err := c.ShouldBindJSON(&bookmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookmark.ID = uint(id)
	bookmark.UserID = userID
	if err := h.service.Update(&bookmark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (h *BookmarkHandler) Delete(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
}
