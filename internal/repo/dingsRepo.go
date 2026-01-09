package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type DingRepository interface {
	CreateDing(ding *model.Ding) error
	CreateDingStudent(ds *model.DingStudent) error
	GetDingsByStudentIDAndStatus(studentID uint, status string) ([]model.Ding, error)
	GetDingsByLauncherID(launcherID uint) ([]model.Ding, error)
	GetDingRecordsByDingID(dingId string) (interface{}, interface{})
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
	if err := r.db.Where("launcher_id = ?", launcherID).Find(&dings).Error; err != nil {
		return nil, err
	}
	return dings, nil
}

func (r *dingRepository) GetDingRecordsByDingID(dingId string) (interface{}, interface{}) {
	var dingRecords []model.DingStudent
	if err := r.db.Where("ding_id = ?", dingId).Find(&dingRecords).Error; err != nil {
		return nil, err
	}
	return dingRecords, nil
}
