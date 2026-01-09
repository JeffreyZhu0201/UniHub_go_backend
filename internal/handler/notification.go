package handler

import (
	"net/http"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	Service service.NotificationService
}

func NewNotificationHandler(s service.NotificationService) *NotificationHandler {
	return &NotificationHandler{Service: s}
}

type CreateNotifRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	TargetType string `json:"target_type" binding:"required,oneof=dept class"` // 目标类型：dept, class
	TargetID   uint   `json:"target_id" binding:"required"`
}

// Create 发布通知 (辅导员/教师)
func (h *NotificationHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	var req CreateNotifRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceReq := service.CreateNotifRequest{
		Title:      req.Title,
		Content:    req.Content,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		SenderID:   userID,
	}

	if err := h.Service.Create(serviceReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知发布成功"})
}

// GetMyNotifications 学生查看通知
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	notifs, err := h.Service.GetMyNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifs)
}
