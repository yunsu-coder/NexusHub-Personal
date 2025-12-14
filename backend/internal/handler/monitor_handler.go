package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct{}

func NewMonitorHandler() *MonitorHandler {
	return &MonitorHandler{}
}

type CheckRequest struct {
	URL string `json:"url" binding:"required"`
}

type CheckResponse struct {
	Status  string `json:"status"`  // "online" | "offline"
	Latency int64  `json:"latency"` // ms
}

func (h *MonitorHandler) CheckHealth(c *gin.Context) {
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(req.URL)
	latency := time.Since(start).Milliseconds()

	status := "offline"
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
		status = "online"
		resp.Body.Close()
	}

	c.JSON(http.StatusOK, CheckResponse{
		Status:  status,
		Latency: latency,
	})
}
