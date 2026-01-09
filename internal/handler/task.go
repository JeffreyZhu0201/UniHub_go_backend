package handler

import (
	"net/http"
	"time"
	"unihub/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Service service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
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

	serviceReq := service.CreateTaskRequest{
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		CreatorID:   userID,
		Deadline:    req.Deadline,
		Config:      req.Config,
	}

	task, err := h.Service.CreateTask(serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务发布成功", "id": task.ID, "uuid": task.UUID})
}

// GetMyTasks 学生查看任务列表
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	userID := c.GetUint("userID")

	tasks, err := h.Service.GetMyTasks(userID)
	if err != nil {
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

	if err := h.Service.SubmitTask(userID, taskUUID, req.Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提交成功"})
}
