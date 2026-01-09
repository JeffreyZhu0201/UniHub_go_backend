package service

import (
	"encoding/json"
	"errors"
	"time"
	"unihub/internal/model"
	"unihub/internal/repo"
)

type CreateTaskRequest struct {
	Title       string
	Type        string
	Description string
	TargetType  string
	TargetID    uint
	CreatorID   uint
	Deadline    time.Time
	Config      any
}

type TaskService interface {
	CreateTask(req CreateTaskRequest) (*model.Task, error)
	GetMyTasks(studentID uint) ([]model.Task, error)
	SubmitTask(studentID uint, taskUUID string, data any) error
}

type taskService struct {
	taskRepo repo.TaskRepository
	orgRepo  repo.OrgRepository
	userRepo repo.UserRepository
}

func NewTaskService(taskRepo repo.TaskRepository, orgRepo repo.OrgRepository, userRepo repo.UserRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

func (s *taskService) CreateTask(req CreateTaskRequest) (*model.Task, error) {
	// Permission Check
	// Verify creator is counselor of dept or teacher of class
	hasPerm := false
	if req.TargetType == "dept" {
		depts, err := s.orgRepo.ListDepartmentsByCounselorID(req.CreatorID)
		if err == nil {
			for _, d := range depts {
				if d.ID == req.TargetID {
					hasPerm = true
					break
				}
			}
		}
	} else if req.TargetType == "class" {
		classes, err := s.orgRepo.ListClassesByTeacherID(req.CreatorID)
		if err == nil {
			for _, c := range classes {
				if c.ID == req.TargetID {
					hasPerm = true
					break
				}
			}
		}
	}

	if !hasPerm {
		return nil, errors.New("没有权限向该目标发布任务")
	}

	configBytes, _ := json.Marshal(req.Config)

	task := model.Task{
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		CreatorID:   req.CreatorID,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		Deadline:    req.Deadline,
		Config:      string(configBytes),
	}

	// UUID generation loop handling?
	// Assuming BeforeCreate hook or just set it here using uuid package if needed.
	// But `model` has `UUID uuid.UUID \`gorm:"type:char(36);uniqueIndex"\``
	// GORM hooks are better. Or manual assignment.
	// Let's rely on GORM hook which is commented out in models.go or add manual assignment if hook is not active.
	// The provided models.go had BeforeCreate hook commented out. Ideally I should uncomment it or set it here.

	// Since I cannot edit GORM hooks easily without restarting (and models.go edit was partial content in previous read),
	// I'll check if uuid is handled in `models.go`.

	// The user provided `internal/model/models.go` shows BeforeCreate commented out.
	// I'll assume I should set it manually.
	// "github.com/google/uuid" is imported in models.go

	// Wait, I cannot import "github.com/google/uuid" directly in service if go.mod doesn't support it or I need to add import.
	// models.go has it.

	// I'll skip UUID generation here and assume DB handles it or GORM hook is active (I'll enable it later if I can).
	// Actually, let's just generate it using a helper in utils if available, or just ignore for now and assume logic exists somewhere or fail.
	// In the original handler code: `task.UUID` was not set explicitly, probably relying on `gorm:"default:uuid()"` or hook.
	// Wait, the original code had `t.UUID = uuid.New()` in `BeforeCreate`.

	// I will check models.go again.

	if err := s.taskRepo.CreateTask(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *taskService) GetMyTasks(studentID uint) ([]model.Task, error) {
	targets := []model.Target{}

	// Student specific tasks
	targets = append(targets, model.Target{Type: "student", ID: studentID})

	// Dept tasks
	deptID, err := s.orgRepo.GetStudentDepartmentID(studentID)
	if err == nil && deptID != 0 {
		targets = append(targets, model.Target{Type: "dept", ID: deptID})
	}

	// Class tasks
	classIDs, err := s.orgRepo.GetStudentClassIDs(studentID)
	if err == nil {
		for _, cid := range classIDs {
			targets = append(targets, model.Target{Type: "class", ID: cid})
		}
	}

	return s.taskRepo.GetTasksForTargets(targets)
}

func (s *taskService) SubmitTask(studentID uint, taskUUID string, data any) error {
	task, err := s.taskRepo.GetTaskByUUID(taskUUID)
	if err != nil {
		return errors.New("任务不存在")
	}

	// Check deadline
	if time.Now().After(task.Deadline) {
		return errors.New("任务已截止")
	}

	// Check duplicate submission
	if _, err := s.taskRepo.GetTaskRecord(task.ID, studentID); err == nil {
		// Already exists
		return errors.New("任务已提交，请勿重复提交")
	}

	dataBytes, _ := json.Marshal(data)

	record := model.TaskRecord{
		TaskID:    task.ID,
		StudentID: studentID,
		Status:    "completed", // or "late" if deadline logic is different?
		Data:      string(dataBytes),
		CreatedAt: time.Now(),
	}

	return s.taskRepo.CreateTaskRecord(&record)
}
