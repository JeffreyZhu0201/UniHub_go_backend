package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type LeaveRepository interface {
	CreateLeaveRequest(leave *model.LeaveRequest) error
	GetLeaveRequestByID(id uint) (*model.LeaveRequest, error)
	UpdateLeaveRequest(leave *model.LeaveRequest) error
	ListPendingLeavesByStudentIDs(studentIDs []uint) ([]model.LeaveRequest, error) // For counselor to audit
	ListLeavesByStudentID(studentID uint) ([]model.LeaveRequest, error)
	ListLeavesWithStudentsByStudentsAndStatus(students interface{}, status string) ([]interface{}, interface{})
}

type leaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{db: db}
}

func (r *leaveRepository) CreateLeaveRequest(leave *model.LeaveRequest) error {
	return r.db.Create(leave).Error
}

func (r *leaveRepository) GetLeaveRequestByID(id uint) (*model.LeaveRequest, error) {
	var leave model.LeaveRequest
	if err := r.db.First(&leave, id).Error; err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *leaveRepository) UpdateLeaveRequest(leave *model.LeaveRequest) error {
	return r.db.Save(leave).Error
}

func (r *leaveRepository) ListPendingLeavesByStudentIDs(studentIDs []uint) ([]model.LeaveRequest, error) {
	var leaves []model.LeaveRequest
	err := r.db.Where("student_id IN ? AND status = ?", studentIDs, "pending").Find(&leaves).Error
	return leaves, err
}

func (r *leaveRepository) ListLeavesByStudentID(studentID uint) ([]model.LeaveRequest, error) {
	var leaves []model.LeaveRequest
	err := r.db.Where("student_id = ?", studentID).Order("created_at desc").Find(&leaves).Error
	return leaves, err
}

func (r *leaveRepository) ListLeavesWithStudentsByStudentsAndStatus(studentsIds interface{}, status string) ([]interface{}, interface{}) {
	var leaves []interface{}
	// 使用连接操作，需要获取学生信息
	if err := r.db.Model(&model.LeaveRequest{}).
		Select("leave_requests.*, users.id as student_id, users.name as student_name").
		Joins("join users on leave_requests.student_id = users.id").
		Where("leave_requests.student_id IN ? AND leave_requests.status = ?", studentsIds, status).
		Scan(&leaves).Error; err != nil {
		return nil, err
	}
	return leaves, nil
}
