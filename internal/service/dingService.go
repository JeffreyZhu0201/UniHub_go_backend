package service

import (
	"errors"
	"log"
	"time"
	"unihub/internal/DTO"
	"unihub/internal/model"
	"unihub/internal/repo"
	"unihub/internal/utils"

	"gorm.io/gorm"
)

func CreateDing(req DTO.CreateDingRequest, DB *gorm.DB) error {
	var studentIDs []uint
	// 如果是向部门发布
	if req.DeptId != 0 {
		if err := DB.Model(&model.StudentDepartment{}).Where("department_id = ?", req.DeptId).Pluck("student_id", &studentIDs).Error; err != nil {
			return errors.New("发生错误")
		}
	} else if req.ClassId != 0 {
		if err := DB.Model(&model.StudentClass{}).Where("class_id = ?", req.ClassId).Pluck("student_id", &studentIDs).Error; err != nil {
			return errors.New("发生错误")
		}
	} else if req.StudentId != 0 {
		if err := DB.Model(&model.User{}).Where("id = ?", req.StudentId).Pluck("id", &studentIDs).Error; err != nil {
			return errors.New("发生错误")
		}
	}
	if len(studentIDs) == 0 {
		return errors.New("发生错误")
	}
	ding := model.Ding{
		LauncherID: req.LauncherId,
		Title:      req.Title,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Radius:     float64(req.Radius),
		UserID:     req.StudentId,
		DeptID:     req.DeptId,
		ClassID:    req.ClassId,
	}

	// Save Ding
	if err := DB.Create(&ding).Error; err != nil {
		return errors.New("发生错误")
	}

	for _, studentID := range studentIDs {
		dingStudent := model.DingStudent{
			DingID:    ding.ID,
			StudentID: studentID,
			Status:    "pending",
			DingTime:  time.Now(),
		}
		if err := DB.Create(&dingStudent).Error; err != nil {
			return errors.New("发生错误")
		}

		notif := model.Notification{
			Title:      "新的打卡任务：" + req.Title,
			Content:    "请在规定时间内完成打卡任务。",
			SenderID:   req.LauncherId,
			TargetType: "student",
			TargetID:   studentID,
		}
		if _, err := utils.PushNotification(notif, DB); err != nil {
			return errors.New("发生错误")
		}
		log.Printf("已向学生 %d 发送打卡任务通知", studentID)
	}
	return nil
}

func ListAllMyDings(id uint, db *gorm.DB) (interface{}, error) {
	// pending + completed
	var result = make(map[string][]model.Ding)

	if pendingDings, err := repo.GetMyDingsByStatus(id, db, "pending"); err == nil {
		result["pending"] = pendingDings
	}
	if completeDings, err := repo.GetMyDingsByStatus(id, db, "complete"); err == nil {
		result["complete"] = completeDings
	}
	return result, nil
}

func ListMyCreatedDings(id uint, db *gorm.DB) (interface{}, interface{}) {
	var dings []model.Ding
	if err := db.Where("launcher_id = ?", id).Find(&dings).Error; err != nil {
		return nil, errors.New("发生错误")
	}
	return dings, nil
}
