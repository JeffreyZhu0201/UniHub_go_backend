package handler

import (
	"net/http"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	Service service.OrgService
}

func NewOrgHandler(s service.OrgService) *OrgHandler {
	return &OrgHandler{Service: s}
}

type CreateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

type JoinOrgRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

// CreateDepartment 辅导员创建部门
func (h *OrgHandler) CreateDepartment(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dept, err := h.Service.CreateDepartment(userID, roleID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "部门创建成功", "id": dept.ID, "invite_code": dept.InviteCode})
}

// CreateClass 教师\导员创建班级
func (h *OrgHandler) CreateClass(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class, err := h.Service.CreateClass(userID, roleID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "班级创建成功", "id": class.ID, "invite_code": class.InviteCode})
}

// StudentJoinDepartment 学生通过邀请码加入部门
func (h *OrgHandler) StudentJoinDepartment(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")
	var req JoinOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.StudentJoinDepartment(userID, roleID, req.InviteCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入部门成功"})
}

// StudentJoinClass 学生通过邀请码加入班级
func (h *OrgHandler) StudentJoinClass(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	var req JoinOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.StudentJoinClass(userID, roleID, req.InviteCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入班级成功"})
}

// ListMyDepartments 辅导员查看其创建的部门
func (h *OrgHandler) ListMyDepartments(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	depts, err := h.Service.ListMyDepartments(userID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}

// ListMyClasses 教师查看其创建的班级
func (h *OrgHandler) ListMyClasses(c *gin.Context) {
	userID := c.GetUint("userID")

	classes, err := h.Service.ListMyClasses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}
