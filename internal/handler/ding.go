package handler

import (
	"log"
	"net/http"
	"time"
	"unihub/internal/model"
	"unihub/internal/service"
	"unihub/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DingHandler struct {
	DB *gorm.DB
}

type CreateDingRequest struct {
	Title     string    `json:"title" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Radius    uint      `json:"radius" binding:"required"`
	StudentId uint      `json:"student_id"`
	DeptId    uint      `json:"dept_id"`
	ClassId   uint      `json:"class_id"`
}

func (d *DingHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	// 权限检查：确保用户角色拥有创建打卡任务权限
	// 这里默认假设 "ding:create" 权限已分配给教师和辅导员角色
	if havePermission, _ := service.RequirePermission(c, d.DB, roleID, "ding:create"); !havePermission {
		c.JSON(403, gin.H{"error": "无权限创建打卡任务"})
		return
	}

	var req CreateDingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var studentIDs []uint
	// 如果是向部门发布
	if req.DeptId != 0 {
		if err := d.DB.Model(&model.StudentDepartment{}).Where("department_id = ?", req.DeptId).Pluck("student_id", &studentIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发生错误"})
			return
		}
	} else if req.ClassId != 0 {
		if err := d.DB.Model(&model.StudentClass{}).Where("class_id = ?", req.ClassId).Pluck("student_id", &studentIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发生错误"})
			return
		}
	} else if req.StudentId != 0 {
		if err := d.DB.Model(&model.User{}).Where("id = ?", req.StudentId).Pluck("id", &studentIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发生错误"})
			return
		}
	}
	if len(studentIDs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未找到目标学生"})
		return
	}
	ding := model.Ding{
		LauncherID: userID,
		Title:      req.Title,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Radius:     float64(req.Radius),
		UserID:     req.StudentId,
		DeptID:     req.DeptId,
		ClassID:    req.ClassId,
	}

	// Save Ding
	if err := d.DB.Create(&ding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, studentID := range studentIDs {
		dingStudent := model.DingStudent{
			DingID:    ding.ID,
			StudentID: studentID,
			Status:    "pending",
			DingTime:  time.Now(),
		}
		if err := d.DB.Create(&dingStudent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		notif := model.Notification{
			Title:      "新的打卡任务：" + req.Title,
			Content:    "请在规定时间内完成打卡任务。",
			SenderID:   userID,
			TargetType: "student",
			TargetID:   studentID,
		}
		if _, err := utils.PushNotification(notif, d.DB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		log.Printf("已向学生 %d 发送打卡任务通知", studentID)
	}
	c.JSON(http.StatusOK, gin.H{"message": "打卡任务创建成功"})
}
