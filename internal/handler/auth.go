package handler

import (
	"log"
	"net/http"
	"unihub/internal/repo"
	"unihub/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"unihub/internal/config"
	"unihub/internal/model"
	"unihub/pkg/jwtutil"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Nickname   string  `json:"nickname" binding:"required"`
	Email      string  `json:"email" binding:"required"`
	Password   string  `json:"password" binding:"required"`
	RoleKey    string  `json:"role_key" binding:"required"`
	StaffNo    *string `json:"staff_no"`    // 教职工号 (辅导员/教师)
	StudentNo  *string `json:"student_no"`  // 学号 (学生)
	InviteCode string  `json:"invite_code"` // 邀请码 (可选)
	PushToken  string  `json:"push_token"`  // 推送令牌
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	PushToken string `json:"push_token"`
}

// AuthHandler 认证处理器
type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var department model.Department
	// 如果邀请码无效，返回错误
	log.Printf(req.InviteCode)
	if req.InviteCode != "" {
		dept, err := repo.GetDepartmentByInviteCode(h.DB, req.InviteCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邀请码"})
			return
		}
		department = dept
	}

	// 如果Username长度小于8且不是以.com结尾，返回错误
	if len(req.Email) < 8 && !utils.EndsWith(req.Email, ".com") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱长度必须至少为8位"})
		return
	}
	if !utils.EndsWith(req.Email, ".com") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱必须以 .com 结尾"})
		return
	}

	var role model.Role
	if err := h.DB.Where("`key` = ?", req.RoleKey).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色"})
		return
	}

	// 密码加密
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Nickname:     req.Nickname,
		Email:        req.Email,
		DepartmentID: department.ID,
		Password:     string(hashed),
		RoleID:       role.ID, // 通过roleKey获取RoleID
		StaffNo:      req.StaffNo,
		StudentNo:    req.StudentNo,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if department.ID != 0 {
		studentDept := model.StudentDepartment{
			StudentID:    user.ID,
			DepartmentID: department.ID,
		}
		if err := h.DB.Create(&studentDept).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	//如果推送令牌存在则 更新推送令牌
	user.PushToken = req.PushToken
	if req.PushToken != "" {
		user.PushToken = req.PushToken
		h.DB.Save(&user)
	}

	// 自动登录返回Token (恢复为正常的小时过期)
	// 之前的问题可能是配置未正确加载导致ExpirationHours为0
	token, err := jwtutil.Generate(h.Cfg.JWT.Secret, h.Cfg.JWT.ExpirationHours, user.ID, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	// 更新推送令牌
	if req.PushToken != "" && req.PushToken != user.PushToken {
		user.PushToken = req.PushToken
		h.DB.Save(&user)
	}

	// 恢复为正常的小时过期
	token, err := jwtutil.Generate(h.Cfg.JWT.Secret, h.Cfg.JWT.ExpirationHours, user.ID, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
