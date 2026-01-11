package handler

import (
	"net/http"
	"unihub/internal/service"
	"unihub/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// GetProfile 获取个人资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.Service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetOrgInfo 获取用户组织详细信息
func (h *UserHandler) GetOrgInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	data, err := h.Service.GetUserOrgInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取组织信息失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ListStudents 列出自己管理的部门或班级的学生
func (h *UserHandler) ListStudents(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	students, err := h.Service.ListStudents(userID, roleID)
	if err != nil {
		if err.Error() == "无权查看" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, students)
}

func (h *UserHandler) ExportListOfObjectsUpload(context *gin.Context) {

	// bind query parameters
	var queryParams struct {
		data []interface{}
	}
	if context.ShouldBindJSON(&queryParams) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "参数绑定失败"})
	}
	filePath, err := utils.ExportToExcel(queryParams.data, "导出")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"filePath": filePath})
}
