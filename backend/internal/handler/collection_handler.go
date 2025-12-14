package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/charset"
)

type ParseURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type Resource struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// ParseURLHandler godoc
// @Summary Parse URL metadata
// @Description Fetch and parse a URL to get title and description
// @Tags collection
// @Accept json
// @Produce json
// @Param url body ParseURLRequest true "URL to parse"
// @Success 200 {object} object{title=string,description=string}
// @Router /api/v1/collections/parse [post]
func (h *CollectionHandler) ParseURLHandler(c *gin.Context) {
	var req ParseURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	// 添加 User-Agent，避免被反爬
	r, _ := http.NewRequest("GET", req.URL, nil)
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch URL"})
		return
	}
	defer resp.Body.Close()

	// 自动检测并转换编码
	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		// 如果转换失败，降级使用原始 body
		utf8Reader = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML"})
		return
	}

	// Parse base URL for resolving relative paths
	baseURL, err := url.Parse(req.URL)
	if err != nil {
		// If base URL parsing fails, we just won't be able to resolve relative paths correctly,
		// but we can still proceed with best effort.
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	desc, _ := doc.Find("meta[name='description']").Attr("content")
	if desc == "" {
		desc, _ = doc.Find("meta[property='og:description']").Attr("content")
	}

	// Helper to resolve URL
	resolveURL := func(src string) string {
		if baseURL == nil {
			return src
		}
		u, err := url.Parse(src)
		if err != nil {
			return src
		}
		return baseURL.ResolveReference(u).String()
	}

	// Try to find an image (og:image or first img)
	image, _ := doc.Find("meta[property='og:image']").Attr("content")
	if image == "" {
		image, _ = doc.Find("img").First().Attr("src")
	}
	if image != "" {
		image = resolveURL(image)
	}

	// Find all potential resources
	var resources []gin.H

	// Images
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			resources = append(resources, gin.H{"type": "image", "url": resolveURL(src)})
		}
	})

	// Videos
	doc.Find("video").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			resources = append(resources, gin.H{"type": "video", "url": resolveURL(src)})
		}
		s.Find("source").Each(func(j int, source *goquery.Selection) {
			src, exists := source.Attr("src")
			if exists && src != "" {
				resources = append(resources, gin.H{"type": "video", "url": resolveURL(src)})
			}
		})
	})

	// Audios
	doc.Find("audio").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			resources = append(resources, gin.H{"type": "audio", "url": resolveURL(src)})
		}
		s.Find("source").Each(func(j int, source *goquery.Selection) {
			src, exists := source.Attr("src")
			if exists && src != "" {
				resources = append(resources, gin.H{"type": "audio", "url": resolveURL(src)})
			}
		})
	})

	c.JSON(http.StatusOK, gin.H{
		"title":       title,
		"description": desc,
		"image":       image,
		"resources":   resources,
	})
}

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
	userID := middleware.GetCurrentUserID(c)
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
	userID := middleware.GetCurrentUserID(c)

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
	userID := middleware.GetCurrentUserID(c)

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
	userID := middleware.GetCurrentUserID(c)

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
	userID := middleware.GetCurrentUserID(c)

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
