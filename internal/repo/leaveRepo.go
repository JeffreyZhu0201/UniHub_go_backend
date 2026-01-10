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
	ListLeavesWithStudentsByStudentsAndStatus(studentIds []uint, status string) ([]interface{}, interface{})
	ListApprovedLeavesWithStudentsByStudents(students []model.User) (interface{}, interface{})
	ListLeavesWithStudentsByStudentsByDingStatusBeforeEnd(students []model.User) (interface{}, interface{})
	ListLeavesWithStudentsByStudentsAfterEnd(students []model.User) (interface{}, interface{})
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

func (r *leaveRepository) ListLeavesWithStudentsByStudentsAndStatus(studentsIds []uint, status string) ([]interface{}, interface{}) {
	var results []map[string]interface{}
	// 使用连接操作，需要获取学生信息
	if err := r.db.Model(&model.LeaveRequest{}).
		Select("leave_requests.*, users.id as student_id, users.nickname as student_name").
		Joins("join users on leave_requests.student_id = users.id").
		Where("leave_requests.student_id IN ? AND leave_requests.status = ?", studentsIds, status).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	var leaves []interface{}
	for _, result := range results {
		leaves = append(leaves, result)
	}
	return leaves, nil
}

func (r *leaveRepository) ListApprovedLeavesWithStudentsByStudents(students []model.User) (interface{}, interface{}) {
	var results []map[string]interface{}
	var studentIDs []uint
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID)
	}
	// 使用连接操作，需要获取学生信息
	if err := r.db.Model(&model.LeaveRequest{}).
		Select("leave_requests.*, users.id as student_id, users.nickname as student_name").
		Joins("join users on leave_requests.student_id = users.id").
		Where("leave_requests.student_id IN ? AND leave_requests.status = ?", studentIDs, "approved").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	print(studentIDs)
	var leaves []interface{}
	for _, result := range results {
		leaves = append(leaves, result)
	}
	return leaves, nil
}

func (r *leaveRepository) ListLeavesWithStudentsByStudentsByDingStatusBeforeEnd(students []model.User) (interface{}, interface{}) {
	var results []map[string]interface{}
	var studentIDs []uint
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID)
	}
	// 使用连接操作，需要获取学生信息,需要通过leave_requests表的end_time字段和当前时间比较,和连接ding_student表
	if err := r.db.Model(&model.LeaveRequest{}).
		Select("leave_requests.*, users.id as student_id, users.nickname as student_name,ding_students.ding_time").
		Joins("join users on leave_requests.student_id = users.id").
		Joins("join ding_students on leave_requests.ding_id = ding_students.ding_id").
		Where("leave_requests.student_id IN ? AND leave_requests.status = ? AND ding_students.status = ? AND leave_requests.end_time > ding_students.ding_time", studentIDs, "approved", "complete").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	var leaves []interface{}
	for _, result := range results {
		leaves = append(leaves, result)
	}
	//if len(leaves) > 0 {
	//	log.Printf(results[0]["student_name"].(string))
	//}
	return leaves, nil
}

func (r *leaveRepository) ListLeavesWithStudentsByStudentsAfterEnd(students []model.User) (interface{}, interface{}) {
	var results []map[string]interface{}
	var studentIDs []uint
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID)
	}
	// 使用连接操作，需要获取学生信息,需要连接ding_student表
	if err := r.db.Model(&model.LeaveRequest{}).
		Select("leave_requests.*, users.id as student_id, users.nickname as student_name,ding_students.ding_time").
		Joins("join users on leave_requests.student_id = users.id").
		Joins("join ding_students on leave_requests.ding_id = ding_students.ding_id").
		Where("leave_requests.student_id IN ? AND leave_requests.status = ? AND leave_requests.end_time < NOW() AND ding_students.status = ?", studentIDs, "approved", "pending").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	var leaves []interface{}
	for _, result := range results {
		leaves = append(leaves, result)
	}
	return leaves, nil
}
