package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		eventService: service.NewEventService(),
	}
}

// GetAllEvents godoc
// @Summary Get all events
// @Description Get all events for a user, optionally filtered by date range
// @Tags event
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} model.Event
// @Router /api/v1/events [get]
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	events, err := h.eventService.GetAll(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEventByID godoc
// @Summary Get an event by ID
// @Description Get a single event by ID for a user
// @Tags event
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} model.Event
// @Router /api/v1/events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.eventService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event for a user
// @Tags event
// @Accept json
// @Produce json
// @Param event body model.Event true "Event data"
// @Success 201 {object} model.Event
// @Router /api/v1/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.UserID = userID

	if err := h.eventService.Create(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an existing event for a user
// @Tags event
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body model.Event true "Event data"
// @Success 200 {object} model.Event
// @Router /api/v1/events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证ID
	existingEvent, err := h.eventService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// 更新字段
	existingEvent.Title = event.Title
	existingEvent.Date = event.Date
	existingEvent.StartTime = event.StartTime
	existingEvent.Type = event.Type
	existingEvent.Description = event.Description
	existingEvent.Remind = event.Remind

	if err := h.eventService.Update(existingEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, existingEvent)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an event for a user
// @Tags event
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 204 {object} nil
// @Router /api/v1/events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	// TODO: 从认证中获取用户ID
	userID := uint(1) // 临时硬编码

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := h.eventService.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.Status(http.StatusNoContent)
}
