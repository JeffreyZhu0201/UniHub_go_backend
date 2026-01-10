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

	if _, err := d.Service.CreateDing(req, userID, roleID); err != nil {
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

func (d *DingHandler) ListMyCreatedDingsRecords(context *gin.Context) {
	userID := context.GetUint("userID")
	dingId := context.Param("dingId")

	// get from repo
	studentRecordByDing, err := d.Service.ListMyCreatedDingsRecords(userID, dingId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "获取打卡记录失败"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"records": studentRecordByDing})
}

// GetDingStats 获取打卡统计 (理论总数, 已打卡, 未打卡)
func (d *DingHandler) GetDingStats(c *gin.Context) {
	userID := c.GetUint("userID")
	// 辅导员或教师都可以看，基于 userID 过滤
	stats, err := d.Service.GetDingStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (d *DingHandler) ExportMyCreatedDingRecords(context *gin.Context) {
	// 导出某一次打卡记录
	//userId := context.GetUint("userID")
	dingId := context.Param("dingId")
	filePath, err := d.Service.ExportMyCreatedDingRecords(dingId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "导出打卡记录失败"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "导出所选打卡记录成功", "fileRelativePath": filePath})
}

func (d *DingHandler) Ding(c *gin.Context) { // 打卡
	dingId := c.Param("dingId")
	userID := c.GetUint("userID")

	dingDetail, err := d.Service.Ding(dingId, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取打卡任务详情失败"})
		return
	}
	c.JSON(http.StatusOK, dingDetail)
}
