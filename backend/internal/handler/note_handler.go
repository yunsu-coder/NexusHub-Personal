package handler

import (
	"net/http"
	"strconv"

	"nexushub-personal/internal/common"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler() *NoteHandler {
	return &NoteHandler{
		service: service.NewNoteService(),
	}
}

func (h *NoteHandler) GetAll(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	notes, err := h.service.GetAll(userID)
	if err != nil {
		common.InternalServerError(c, err.Error())
		return
	}
	common.Success(c, notes)
}

func (h *NoteHandler) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid ID")
		return
	}

	note, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		common.NotFound(c, "Note not found")
		return
	}
	common.Success(c, note)
}

func (h *NoteHandler) Create(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	note.UserID = userID
	if err := h.service.Create(&note); err != nil {
		common.InternalServerError(c, err.Error())
		return
	}
	common.Created(c, note)
}

func (h *NoteHandler) Update(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid ID")
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.ID = uint(id)
	note.UserID = userID
	if err := h.service.Update(&note); err != nil {
		common.InternalServerError(c, err.Error())
		return
	}
	common.Success(c, note)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid ID")
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		common.InternalServerError(c, err.Error())
		return
	}
	common.SuccessWithMessage(c, "Note deleted successfully", nil)
}
