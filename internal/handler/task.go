package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unihub/internal/model"
)

type TaskHandler struct {
	DB *gorm.DB
}

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Type        string    `json:"type" binding:"required,oneof=sign_in dorm_check"`
	Description string    `json:"description"`
	TargetType  string    `json:"target_type" binding:"required,oneof=dept class"`
	TargetID    uint      `json:"target_id" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	Config      any       `json:"config"` // JSON Object
}

type SubmitTaskRequest struct {
	Data any `json:"data" binding:"required"`
}

// CreateTask 发布任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 权限检查 (略，同 Notification)

	configBytes, _ := json.Marshal(req.Config)

	task := model.Task{
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		CreatorID:   userID,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		Deadline:    req.Deadline,
		Config:      string(configBytes),
	}

	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务发布成功", "id": task.ID, "uuid": task.UUID})
}

// GetMyTasks 学生查看任务列表
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	userID := c.GetUint("userID")

	// 获取所属部门ID和班级ID
	var deptID uint
	h.DB.Model(&model.StudentDepartment{}).Where("student_id = ?", userID).Pluck("department_id", &deptID)
	var classIDs []uint
	h.DB.Model(&model.StudentClass{}).Where("student_id = ?", userID).Pluck("class_id", &classIDs)

	var tasks []model.Task
	query := h.DB.Where("(target_type = ? AND target_id = ?)", "dept", deptID)
	if len(classIDs) > 0 {
		query = query.Or("(target_type = ? AND target_id IN ?)", "class", classIDs)
	}

	// 也查询 target_type = student AND target_id = userID (例如销假签到)
	query = query.Or("(target_type = ? AND target_id = ?)", "student", userID)

	if err := query.Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// SubmitTask 学生提交任务
func (h *TaskHandler) SubmitTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskUUID := c.Param("uuid")

	var req SubmitTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task model.Task
	if err := h.DB.Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if time.Now().After(task.Deadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务已截止"})
		return
	}

	// 检查是否已提交
	var count int64
	h.DB.Model(&model.TaskRecord{}).Where("task_id = ? AND student_id = ?", task.ID, userID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请勿重复提交"})
		return
	}

	dataBytes, _ := json.Marshal(req.Data)

	record := model.TaskRecord{
		TaskID:    task.ID,
		StudentID: userID,
		Status:    "completed",
		Data:      string(dataBytes),
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提交成功"})
}
