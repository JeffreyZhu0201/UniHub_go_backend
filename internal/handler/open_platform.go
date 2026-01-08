package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"unihub/internal/model"
)

type OpenHandler struct {
	DB *gorm.DB
}

type RegisterDevRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CreateAppRequest struct {
	Name string `json:"name" binding:"required"`
}

// RegisterDeveloper 注册开发者
func (h *OpenHandler) RegisterDeveloper(c *gin.Context) {
	var req RegisterDevRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成 Secret
	bytes := make([]byte, 32)
	rand.Read(bytes)
	secret := hex.EncodeToString(bytes)

	dev := model.Developer{
		UUID:   uuid.New(),
		Name:   req.Name,
		Email:  req.Email,
		Secret: secret,
	}

	if err := h.DB.Create(&dev).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮箱已被注册或系统错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "注册成功",
		"dev_id":     dev.ID,
		"dev_secret": dev.Secret, // 仅显示一次
	})
}

// CreateApp 创建应用
func (h *OpenHandler) CreateApp(c *gin.Context) {
	// 简单验证：通过 Header "X-Dev-Secret" 验证开发者身份 (实际应更复杂)
	secret := c.GetHeader("X-Dev-Secret")
	var dev model.Developer
	if err := h.DB.Where("secret = ?", secret).First(&dev).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的开发者密钥"})
		return
	}

	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成 AppID, AppSecret
	appIDData := make([]byte, 8)
	rand.Read(appIDData)
	appID := hex.EncodeToString(appIDData)

	appSecretData := make([]byte, 16)
	rand.Read(appSecretData)
	appSecret := hex.EncodeToString(appSecretData)

	app := model.App{
		UUID:        uuid.New(),
		DeveloperID: dev.ID,
		Name:        req.Name,
		AppID:       appID,
		AppSecret:   appSecret,
		RateLimit:   60, // 默认 60 req/min
	}

	if err := h.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "应用创建成功",
		"app_id":     app.AppID,
		"app_secret": app.AppSecret,
	})
}
