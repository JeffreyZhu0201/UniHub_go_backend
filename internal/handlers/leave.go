package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unihub-backend/internal/models"
)

type LeaveHandler struct {
	db *gorm.DB
}

func NewLeaveHandler(db *gorm.DB) *LeaveHandler {
	return &LeaveHandler{db: db}
}

func (h *LeaveHandler) GetLeaves(c *gin.Context) {
	var leaves []models.Leave

	result := h.db.Preload("Student").
		Order("created_at desc").
		Find(&leaves)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": leaves})
}

func (h *LeaveHandler) AuditLeave(c *gin.Context) {
	type Request struct {
		LeaveID string `json:"leave_id"`
		Status  int    `json:"status"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var leave models.Leave
	if err := h.db.First(&leave, "id = ?", req.LeaveID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		return
	}

	leave.Status = req.Status
	h.db.Save(&leave)
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}