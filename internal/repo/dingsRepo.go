package repo

import (
	"errors"
	"unihub/internal/model"

	"gorm.io/gorm"
)

// 学生打卡任务,返回该学生未完成打卡任务
func GetMyDingsByStatus(id uint, db *gorm.DB, status string) ([]model.Ding, error) {
	var dings []uint
	if err := db.Model(&model.DingStudent{}).Where("student_id = ? AND status = ?", id, status).Pluck("ding_id", &dings).Error; err != nil {
		return nil, errors.New("发生错误")
	}
	var dingDetails []model.Ding
	if err := db.Where("id IN ?", dings).Find(&dingDetails).Error; err != nil {
		return nil, errors.New("发生错误")
	}
	return dingDetails, nil
}

//
//// 教师/辅导员查看自己创建的打卡任务
//func ListMyCreatedDings(id uint, db *gorm.DB) (interface{}, error) {
//	var dings []model.Ding
//	if err := db.Where("launcher_id = ?", id).Find(&dings).Error; err != nil {
//		return nil, errors.New("发生错误")
//	}
//	return dings, nil
//}
