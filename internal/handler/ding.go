package handler

import (
	"net/http"
	"unihub/internal/DTO"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DingHandler struct {
	DB *gorm.DB
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

	var req DTO.CreateDingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.LauncherId = userID // 设置发布者为当前用户

	if err := service.CreateDing(req, d.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建打卡任务失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建打卡任务成功"})
}
