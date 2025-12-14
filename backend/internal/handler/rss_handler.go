package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
)

type RSSHandler struct{}

func NewRSSHandler() *RSSHandler {
	return &RSSHandler{}
}

type RSSRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *RSSHandler) GetFeed(c *gin.Context) {
	var req RSSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse RSS feed"})
		return
	}

	// 简化返回数据，只取前 5 条
	items := make([]map[string]string, 0)
	for i, item := range feed.Items {
		if i >= 5 {
			break
		}
		items = append(items, map[string]string{
			"title": item.Title,
			"link":  item.Link,
			"date":  item.Published,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"title": feed.Title,
		"items": items,
	})
}
