package handler

import (
	"nexushub-personal/internal/constants"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CRUDService 定义CRUD服务接口
type CRUDService[T any] interface {
	GetAll(userID uint) ([]T, error)
	GetByID(id, userID uint) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id, userID uint) error
}

// Entity 定义实体接口 - 需要有ID和UserID字段
type Entity interface {
	SetID(id uint)
	SetUserID(userID uint)
}

// BaseHandler 通用CRUD Handler基类
type BaseHandler[T Entity] struct {
	service     CRUDService[T]
	resourceName string // 资源名称，用于错误消息
}

// NewBaseHandler 创建基础Handler
func NewBaseHandler[T Entity](service CRUDService[T], resourceName string) *BaseHandler[T] {
	return &BaseHandler[T]{
		service:      service,
		resourceName: resourceName,
	}
}

// GetAll 获取所有资源
func (h *BaseHandler[T]) GetAll(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	items, err := h.service.GetAll(userID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.SuccessResponse(c, items)
}

// GetByID 根据ID获取资源
func (h *BaseHandler[T]) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestError(c, constants.ErrInvalidID)
		return
	}

	item, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundError(c, h.resourceName+" not found")
		} else {
			utils.InternalServerError(c, err.Error())
		}
		return
	}
	utils.SuccessResponse(c, item)
}

// Create 创建资源
func (h *BaseHandler[T]) Create(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	entity.SetUserID(userID)
	if err := h.service.Create(&entity); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.CreatedResponse(c, entity)
}

// Update 更新资源
func (h *BaseHandler[T]) Update(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestError(c, constants.ErrInvalidID)
		return
	}

	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	entity.SetID(uint(id))
	entity.SetUserID(userID)
	if err := h.service.Update(&entity); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.SuccessResponse(c, entity)
}

// Delete 删除资源
func (h *BaseHandler[T]) Delete(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestError(c, constants.ErrInvalidID)
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundError(c, h.resourceName+" not found")
		} else {
			utils.InternalServerError(c, err.Error())
		}
		return
	}
	utils.MessageResponse(c, h.resourceName+" deleted successfully")
}
