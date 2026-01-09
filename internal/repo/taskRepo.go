package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *model.Task) error
	GetTaskByUUID(uuid string) (*model.Task, error)
	GetTasksForTargets(targets []model.Target) ([]model.Task, error)
	CreateTaskRecord(record *model.TaskRecord) error
	GetTaskRecord(taskID, studentID uint) (*model.TaskRecord, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetTaskByUUID(uuidStr string) (*model.Task, error) {
	var task model.Task
	if err := r.db.Where("uuid = ?", uuidStr).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetTasksForTargets(targets []model.Target) ([]model.Task, error) {
	if len(targets) == 0 {
		return []model.Task{}, nil
	}

	query := r.db.Model(&model.Task{})
	for i, t := range targets {
		if i == 0 {
			query = query.Where("target_type = ? AND target_id = ?", t.Type, t.ID)
		} else {
			query = query.Or("target_type = ? AND target_id = ?", t.Type, t.ID)
		}
	}

	var tasks []model.Task
	err := query.Order("created_at desc").Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) CreateTaskRecord(record *model.TaskRecord) error {
	return r.db.Create(record).Error
}

func (r *taskRepository) GetTaskRecord(taskID, studentID uint) (*model.TaskRecord, error) {
	var record model.TaskRecord
	if err := r.db.Where("task_id = ? AND student_id = ?", taskID, studentID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}
