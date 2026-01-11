package repo

import (
	"log"
	"strconv"
	"unihub/internal/model"

	"gorm.io/gorm"
)

type DingRepository interface {
	CreateDing(ding *model.Ding) error
	CreateDingStudent(ds *model.DingStudent) error
	GetDingsByStudentIDAndStatus(studentID uint, status string) ([]model.Ding, error)
	GetDingsByLauncherID(launcherID uint) ([]model.Ding, error)
	GetDingRecordsByDingID(dingId string) (interface{}, interface{})
	UpdateDingStudent(dingId string, userId uint) (interface{}, interface{})
	GetDingStats(launcherID uint) (int64, int64, error)
}

type dingRepository struct {
	db *gorm.DB
}

func NewDingRepository(db *gorm.DB) DingRepository {
	return &dingRepository{db: db}
}

func (r *dingRepository) CreateDing(ding *model.Ding) error {
	return r.db.Create(ding).Error
}

func (r *dingRepository) CreateDingStudent(ds *model.DingStudent) error {
	return r.db.Create(ds).Error
}

func (r *dingRepository) GetDingsByStudentIDAndStatus(studentID uint, status string) ([]model.Ding, error) {
	var dingIDs []uint
	if err := r.db.Model(&model.DingStudent{}).Where("student_id = ? AND status = ?", studentID, status).Pluck("ding_id", &dingIDs).Error; err != nil {
		return nil, err
	}

	if len(dingIDs) == 0 {
		return []model.Ding{}, nil
	}

	var dings []model.Ding
	if err := r.db.Where("id IN ?", dingIDs).Find(&dings).Error; err != nil {
		return nil, err
	}
	return dings, nil
}

func (r *dingRepository) GetDingsByLauncherID(launcherID uint) ([]model.Ding, error) {
	var dings []model.Ding
	if err := r.db.Where("launcher_id = ?", launcherID).Order("created_at desc").Find(&dings).Error; err != nil {
		return nil, err
	}
	return dings, nil
}

// GetDingRecordsByDingID modified to include student name and no
func (r *dingRepository) GetDingRecordsByDingID(dingId string) (interface{}, interface{}) {
	var results []map[string]interface{}
	// Join ding_students with users to get student details
	err := r.db.Table("ding_students").
		Select("ding_students.*, users.nickname as student_name, users.student_no").
		Joins("JOIN users ON ding_students.student_id = users.id").
		Where("ding_students.ding_id = ?", dingId).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *dingRepository) UpdateDingStudent(dingId string, userId uint) (interface{}, interface{}) {
	var dingStudent model.DingStudent
	if err := r.db.Where("ding_id = ? AND student_id = ?", dingId, userId).First(&dingStudent).Error; err != nil {
		return nil, err
	}
	dingStudent.Status = "complete"
	if err := r.db.Save(&dingStudent).Error; err != nil {
		return nil, err
	}
	return dingStudent, nil
}

func (r *dingRepository) GetDingStats(launcherID uint) (int64, int64, error) {
	var total int64
	var checked int64

	// 1. 计算理论总人数 (所有非这类任务的学生关联记录)
	// 假设 ding_students 表名为 ding_students, dings 表名为 dings
	// 过滤 Type != 'leave_return'
	if err := r.db.Table("ding_students").
		Joins("JOIN dings ON dings.id = ding_students.ding_id").
		Where("dings.launcher_id = ? AND dings.title != ?", launcherID, "返校签到").
		Count(&total).Error; err != nil {
		return 0, 0, err
	}

	// 2. 计算已打卡人数 (Status = true/1)
	if err := r.db.Table("ding_students").
		Joins("JOIN dings ON dings.id = ding_students.ding_id").
		Where("dings.launcher_id = ? AND dings.title != ? AND ding_students.status = ?", launcherID, "返校签到", "complete").
		Count(&checked).Error; err != nil {
		return 0, 0, err
	}
	log.Printf(strconv.FormatInt(checked, 10), total)
	return total, checked, nil
}
