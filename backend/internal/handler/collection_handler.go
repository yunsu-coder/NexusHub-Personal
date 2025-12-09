package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"
)

type CollectionHandler struct {
	collectionService *service.CollectionService
}

func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{
		collectionService: service.NewCollectionService(),
	}
}

// GetAllCollections godoc
// @Summary Get all collections
// @Description Get all collections for a user, optionally filtered by type
// @Tags collection
// @Accept json
// @Produce json
// @Param type query string false "Collection type filter"
// @Success 200 {array} model.Collection
// @Router /api/v1/collections [get]
func (h *CollectionHandler) GetAllCollections(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码
	collectionType := c.Query("type")

	collections, err := h.collectionService.GetAll(userID, collectionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}

	c.JSON(http.StatusOK, collections)
}

// GetCollectionByID godoc
// @Summary Get a collection by ID
// @Description Get a single collection by ID for a user
// @Tags collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} model.Collection
// @Router /api/v1/collections/{id} [get]
func (h *CollectionHandler) GetCollectionByID(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	collection, err := h.collectionService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// CreateCollection godoc
// @Summary Create a new collection
// @Description Create a new collection for a user
// @Tags collection
// @Accept json
// @Produce json
// @Param collection body model.Collection true "Collection data"
// @Success 201 {object} model.Collection
// @Router /api/v1/collections [post]
func (h *CollectionHandler) CreateCollection(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	var collection model.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection.UserID = userID

	if err := h.collectionService.Create(&collection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}

	c.JSON(http.StatusCreated, collection)
}

// UpdateCollection godoc
// @Summary Update a collection
// @Description Update an existing collection for a user
// @Tags collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param collection body model.Collection true "Collection data"
// @Success 200 {object} model.Collection
// @Router /api/v1/collections/{id} [put]
func (h *CollectionHandler) UpdateCollection(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var collection model.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证ID
	existingCollection, err := h.collectionService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// 更新字段
	existingCollection.Title = collection.Title
	existingCollection.URL = collection.URL
	existingCollection.Type = collection.Type
	existingCollection.Thumbnail = collection.Thumbnail
	existingCollection.Description = collection.Description
	existingCollection.Tags = collection.Tags

	if err := h.collectionService.Update(existingCollection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection"})
		return
	}

	c.JSON(http.StatusOK, existingCollection)
}

// DeleteCollection godoc
// @Summary Delete a collection
// @Description Delete a collection for a user
// @Tags collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} nil
// @Router /api/v1/collections/{id} [delete]
func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	if err := h.collectionService.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}

	c.Status(http.StatusNoContent)
}
