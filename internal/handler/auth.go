package handler

import (
	"net/http"
	"unihub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"unihub/internal/config"
	"unihub/internal/model"
	"unihub/pkg/jwtutil"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username  string  `json:"username" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	RoleKey   string  `json:"role_key" binding:"required"`
	OrgID     *uint   `json:"org_id"`
	StaffNo   *string `json:"staff_no"`   // 教职工号 (辅导员/教师)
	StudentNo *string `json:"student_no"` // 学号 (学生)
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	// 如果Username长度小于8且不是以.com结尾，返回错误
	if len(req.Username) < 8 && !utils.EndsWith(req.Username, ".com") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱长度必须至少为8位"})
		return
	}
	if !utils.EndsWith(req.Username, ".com") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱必须以 .com 结尾"})
		return
	}

	var role model.Role
	if err := h.DB.Where("key = ?", req.RoleKey).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色"})
		return
	}

	// 密码加密
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		UUID:      uuid.New(),
		Username:  req.Username,
		Password:  string(hashed),
		RoleID:    role.ID, // 通过roleKey获取RoleID
		OrgUnitID: req.OrgID,
		StaffNo:   req.StaffNo,
		StudentNo: req.StudentNo,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 自动登录返回Token (恢复为正常的小时过期)
	// 之前的问题可能是配置未正确加载导致ExpirationHours为0
	token, err := jwtutil.Generate(h.Cfg.JWT.Secret, h.Cfg.JWT.ExpirationHours, user.ID, user.RoleID, user.OrgUnitID)
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
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	// 恢复为正常的小时过期
	token, err := jwtutil.Generate(h.Cfg.JWT.Secret, h.Cfg.JWT.ExpirationHours, user.ID, user.RoleID, user.OrgUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
