package handler

import (
	"net/http"
	"unihub/internal/DTO"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
)

type DingHandler struct {
	Service service.DingService
	Auth    service.AuthService // or UserService for permission check helper if needed
	// Actually we can use service.RequirePermission helper directly if we pass DB, but better to put permission check in Service or keep using helper with DB if helper is standalone.
	// But `service.RequirePermission` uses *gorm.DB. We should probably move RBAC to a service or repository.
	// It is in `service/rbac.go`. Let's assume we can access UserRepo which has CheckPermission.
	UserRepo interface {
		CheckPermission(roleID uint, permCode string) (bool, error)
	}
}

func NewDingHandler(s service.DingService, uRepo interface {
	CheckPermission(roleID uint, permCode string) (bool, error)
}) *DingHandler {
	return &DingHandler{Service: s, UserRepo: uRepo}
}

func (d *DingHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	// 权限检查
	if havePermission, _ := d.UserRepo.CheckPermission(roleID, "ding:create"); !havePermission {
		c.JSON(403, gin.H{"error": "无权限创建打卡任务"})
		return
	}

	var req DTO.CreateDingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.LauncherId = userID // 设置发布者为当前用户

	if err := d.Service.CreateDing(req, userID, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建打卡任务失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建打卡任务成功"})
}

func (d *DingHandler) ListMyDings(context *gin.Context) {
	userID := context.GetUint("userID")

	dingsDetails, err := d.Service.ListAllMyDings(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "获取打卡任务失败"})
		return
	}
	context.JSON(http.StatusOK, dingsDetails)
}

func (d *DingHandler) ListMyCreatedDings(context *gin.Context) {
	userID := context.GetUint("userID")
	//roleID := context.GetUint("roleID")

	// 权限检查 skipped as per original code commented out

	dings, err := d.Service.ListMyCreatedDings(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "获取我创建的打卡任务失败"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"dings": dings})
}
