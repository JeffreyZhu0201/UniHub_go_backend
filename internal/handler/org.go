package handler

import (
	"net/http"
	"unihub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"unihub/internal/model"
	"unihub/internal/service"
)

type OrgHandler struct {
	DB *gorm.DB
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

	// 权限检查：确保用户角色拥有创建部门权限
	// 这里默认假设 "dept:create" 权限已分配给辅导员角色
	isCounselor, _ := service.RequirePermission(c, h.DB, roleID, "dept:create")
	if !isCounselor {
		// 备用逻辑：如果是辅导员角色也可以
		// ...
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviteCode := utils.GenerateInviteCode()
	// 实际生产环境应检查邀请码唯一性
	for err := h.DB.Where("invite_code = ?", inviteCode).First(&model.Department{}).Error; err == nil; {
		inviteCode = utils.GenerateInviteCode()
	}

	dept := model.Department{
		UUID:        uuid.New(),
		Name:        req.Name,
		InviteCode:  inviteCode,
		CounselorID: userID,
	}

	if err := h.DB.Create(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "部门创建成功", "id": dept.UUID, "invite_code": dept.InviteCode})
}

// CreateClass 教师创建班级
func (h *OrgHandler) CreateClass(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查User是否是教师
	var user model.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	if err := h.DB.Where("id = ? AND role_id = (SELECT id FROM roles WHERE `key` = ?)", userID, "teacher").First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "仅教师可创建班级"})
		return
	}

	// 生成唯一邀请码
	inviteCode := utils.GenerateInviteCode()
	for err := h.DB.Where("invite_code = ?", inviteCode).First(&model.Class{}).Error; err == nil; {
		inviteCode = utils.GenerateInviteCode()
	}

	class := model.Class{
		UUID:       uuid.New(),
		Name:       req.Name,
		InviteCode: inviteCode,
		TeacherID:  userID,
	}

	if err := h.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "班级创建成功", "id": class.UUID, "invite_code": class.InviteCode})
}

// StudentJoinDepartment 学生通过邀请码加入部门
func (h *OrgHandler) StudentJoinDepartment(c *gin.Context) {
	userID := c.GetUint("userID")
	var req JoinOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dept model.Department
	if err := h.DB.Where("invite_code = ?", req.InviteCode).First(&dept).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀请码无效"})
		return
	}

	// 检查是否已加入该部门
	var count int64
	h.DB.Model(&model.StudentDepartment{}).Where("student_id = ? AND department_id = ?", userID, dept.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已加入该部门"})
		return
	}

	// 限制：一个学生只能加入一个部门
	var existingCount int64
	h.DB.Model(&model.StudentDepartment{}).Where("student_id = ?", userID).Count(&existingCount)
	if existingCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生只能加入一个部门"})
		return
	}

	link := model.StudentDepartment{
		StudentID:    userID,
		DepartmentID: dept.ID,
	}
	if err := h.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入部门成功"})
}

// StudentJoinClass 学生通过邀请码加入班级
func (h *OrgHandler) StudentJoinClass(c *gin.Context) {
	userID := c.GetUint("userID")
	var req JoinOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var class model.Class
	if err := h.DB.Where("invite_code = ?", req.InviteCode).First(&class).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀请码无效"})
		return
	}

	// 检查是否已加入该班级
	var count int64
	h.DB.Model(&model.StudentClass{}).Where("student_id = ? AND class_id = ?", userID, class.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已加入该班级"})
		return
	}

	link := model.StudentClass{
		StudentID: userID,
		ClassID:   class.ID,
	}
	if err := h.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入班级成功"})
}

// ListMyDepartments 辅导员查看其创建的部门
func (h *OrgHandler) ListMyDepartments(c *gin.Context) {
	userID := c.GetUint("userID")
	var depts []model.Department
	if err := h.DB.Where("counselor_id = ?", userID).Find(&depts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}

// ListMyClasses 教师查看其创建的班级
func (h *OrgHandler) ListMyClasses(c *gin.Context) {
	userID := c.GetUint("userID")
	var classes []model.Class
	if err := h.DB.Where("teacher_id = ?", userID).Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}
