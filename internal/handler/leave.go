package handler

import (
	"net/http"
	"time"
	"unihub/internal/DTO"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unihub/internal/model"
)

type LeaveHandler struct {
	DB *gorm.DB
}

type ApplyLeaveRequest struct {
	Type      string    `json:"type" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
}

type AuditLeaveRequest struct {
	Status  string `json:"status" binding:"required,oneof=approved rejected"`
	LeaveID uint   `json:"leave_id" binding:"required"`
}

// Apply 申请请假
func (h *LeaveHandler) Apply(c *gin.Context) {
	userID := c.GetUint("userID")
	var req ApplyLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave := model.LeaveRequest{
		StudentID: userID,
		Type:      req.Type,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Reason:    req.Reason,
		Status:    "pending",
	}

	if err := h.DB.Create(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "申请成功", "id": leave.ID})
}

// Audit 审批请假 (辅导员)
func (h *LeaveHandler) Audit(c *gin.Context) {
	auditorID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	var req AuditLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if havePermission, _ := service.RequirePermission(c, h.DB, roleID, "leave:approve"); !havePermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限审批"})
	}

	var leave model.LeaveRequest
	if err := h.DB.Where("id = ?", req.LeaveID).First(&leave).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "请假记录不存在"})
		return
	}

	// 权限检查：审批人是否是该学生的辅导员
	// 找到学生所在部门
	var dept model.StudentDepartment
	if err := h.DB.Where("student_id = ?", leave.StudentID).First(&dept).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生未加入部门"})
		return
	}

	// 检查辅导员是否管理该部门
	var count int64
	h.DB.Model(&model.Department{}).Where("id = ? AND counselor_id = ?", dept.DepartmentID, auditorID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权审批该学生请假"})
		return
	}

	now := time.Now()
	leave.Status = req.Status
	leave.AuditorID = &auditorID
	leave.AuditTime = &now

	// 插入
	if err := h.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 附加功能：审批通过后自动生成签到任务 (销假签到)
	if req.Status == "approved" {
		ding := DTO.CreateDingRequest{
			StudentId:  leave.StudentID,
			LauncherId: leave.StudentID,
			Title:      "已完成请假，请在规定时间内销假签到",
			StartTime:  leave.EndTime,
			EndTime:    leave.EndTime.Add(2 * time.Hour), // 销假签到时间窗口为请假结束后2小时内
			Latitude:   200.0,
			Longitude:  200.0,
			Radius:     200.0,
		}

		if err := service.CreateDing(ding, h.DB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "审批成功，但创建销假签到任务失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "审批完成"})
}

// MyLeaves 学生查看请假记录
func (h *LeaveHandler) MyLeaves(c *gin.Context) {
	userID := c.GetUint("userID")
	var leaves []model.LeaveRequest
	if err := h.DB.Where("student_id = ?", userID).Order("created_at desc").Find(&leaves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaves)
}

// ListPendingLeaves 辅导员查看待审批
func (h *LeaveHandler) ListPendingLeaves(c *gin.Context) {
	counselorID := c.GetUint("userID")
	roleID := c.GetUint("roleID")
	// 权限检查
	if havePermission, _ := service.RequirePermission(c, h.DB, roleID, "leave:approve"); !havePermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限审批"})
	}

	// 找到管理的部门
	var deptIDs []uint

	// 检查我的部门
	h.DB.Model(&model.Department{}).Where("counselor_id = ?", counselorID).Pluck("id", &deptIDs)

	if len(deptIDs) == 0 {
		c.JSON(http.StatusOK, []model.LeaveRequest{})
		return
	}

	// 找到这些部门的学生
	var studentIDs []uint
	h.DB.Model(&model.StudentDepartment{}).Where("department_id IN ?", deptIDs).Pluck("student_id", &studentIDs)

	if len(studentIDs) == 0 {
		c.JSON(http.StatusOK, []model.LeaveRequest{})
		return
	}

	// 找出我的部门未审核的学生请假记录
	var leaves []model.LeaveRequest
	h.DB.Where("student_id IN ? AND status = ?", studentIDs, "pending").Find(&leaves)

	c.JSON(http.StatusOK, leaves)
}
