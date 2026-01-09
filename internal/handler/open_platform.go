package handler

import (
	"net/http"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
)

type OpenHandler struct {
	Service service.OpenService
}

func NewOpenHandler(s service.OpenService) *OpenHandler {
	return &OpenHandler{Service: s}
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

	dev, err := h.Service.RegisterDeveloper(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	// 简单验证：通过 Header "X-Dev-Secret" 验证开发者身份
	secret := c.GetHeader("X-Dev-Secret")
	if secret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要开发者密钥"})
		return
	}

	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.Service.CreateApp(secret, req.Name)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "无效的开发者密钥" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "应用创建成功",
		"app_id":     app.AppID,
		"app_secret": app.AppSecret,
	})
}
