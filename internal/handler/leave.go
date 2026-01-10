package handler

import (
	"net/http"
	"time"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	leaveService service.LeaveService
	dingService  service.DingService
}

func NewLeaveHandler(s service.LeaveService, d service.DingService) *LeaveHandler {
	return &LeaveHandler{leaveService: s}
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

	serviceReq := service.ApplyLeaveRequest{
		StudentID: userID,
		Type:      req.Type,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Reason:    req.Reason,
	}

	leave, err := h.leaveService.Apply(serviceReq)
	if err != nil {
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

	serviceReq := service.AuditLeaveRequest{
		AuditorID: auditorID,
		RoleID:    roleID,
		LeaveID:   req.LeaveID,
		Status:    req.Status,
	}

	if err := h.leaveService.Audit(serviceReq, h.dingService); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "无权限审批" || err.Error() == "无权审批该学生请假" {
			status = http.StatusForbidden
		} else if err.Error() == "请假记录不存在" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "审批成功"})
}

// ListPendingLeaves 辅导员查看待审批请假
func (h *LeaveHandler) ListPendingLeaves(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	leaves, err := h.leaveService.ListPendingLeaves(userID, roleID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaves)
}

// MyLeaves 学生查看自己的请假
func (h *LeaveHandler) MyLeaves(c *gin.Context) {
	userID := c.GetUint("userID")

	leaves, err := h.leaveService.MyLeaves(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaves)
}

func (h *LeaveHandler) LeaveData(context *gin.Context) {
	userId := context.GetUint("userID")

	data, err := h.leaveService.LeaveData(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "查询请假数据失败"})
		return
	}

	context.JSON(http.StatusOK, data)

}
