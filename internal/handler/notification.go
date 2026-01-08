package handler

import (
	"net/http"
	"unihub/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unihub/internal/model"
)

type NotificationHandler struct {
	DB *gorm.DB
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

	// 权限检查：需确认用户是否有权向该部门/班级发通知
	// 简化逻辑：查询是否是该部门/班级的创建者
	hasPerm := false
	if req.TargetType == "dept" {
		var count int64
		h.DB.Model(&model.Department{}).Where("id = ? AND counselor_id = ?", req.TargetID, userID).Count(&count)
		hasPerm = count > 0
	} else {
		var count int64
		h.DB.Model(&model.Class{}).Where("id = ? AND teacher_id = ?", req.TargetID, userID).Count(&count)
		hasPerm = count > 0
	}

	if !hasPerm {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限向该目标发送通知"})
		return
	}

	notif := model.Notification{
		Title:      req.Title,
		Content:    req.Content,
		SenderID:   userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	if err := h.DB.Create(&notif).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := utils.PushNotification(notif, h.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知发布成功"})
}

// GetMyNotifications 学生查看通知
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	// 找到学生所属的部门ID
	var deptID uint
	h.DB.Model(&model.StudentDepartment{}).Where("student_id = ?", userID).Pluck("department_id", &deptID)

	// 找到学生所属的班级IDs
	var classIDs []uint
	h.DB.Model(&model.StudentClass{}).Where("student_id = ?", userID).Pluck("class_id", &classIDs)

	// 查询目标为这些部门或班级的通知
	var notifs []model.Notification
	query := h.DB.Where("(target_type = ? AND target_id = ?)", "dept", deptID)
	if len(classIDs) > 0 {
		query = query.Or("(target_type = ? AND target_id IN ?)", "class", classIDs)
	}

	if err := query.Order("created_at desc").Find(&notifs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifs)
}
