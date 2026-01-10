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

type DingService interface {
	CreateDing(req DTO.CreateDingRequest, launcherID uint, roleID uint) (uint, error)
	ListAllMyDings(studentID uint) (map[string][]model.Ding, error)
	ListMyCreatedDings(launcherID uint) ([]model.Ding, error)
	ListMyCreatedDingsRecords(userId uint, dingID string) (interface{}, interface{})
	ExportMyCreatedDingRecords(dingId string) (interface{}, interface{})
	Ding(userId string, dingId uint) (interface{}, interface{})
	GetDingStats(launcherID uint) (map[string]int64, error)
}

type dingService struct {
	dingRepo repo.DingRepository
	orgRepo  repo.OrgRepository
	userRepo repo.UserRepository
	db       *gorm.DB // Kept for transaction or utils.PushNotification if refactoring notification is not done yet
}

func NewDingService(dingRepo repo.DingRepository, orgRepo repo.OrgRepository, userRepo repo.UserRepository, db *gorm.DB) DingService {
	return &dingService{
		dingRepo: dingRepo,
		orgRepo:  orgRepo,
		userRepo: userRepo,
		db:       db,
	}
}

func (s *dingService) CreateDing(req DTO.CreateDingRequest, launcherID uint, _ uint) (uint, error) {
	var studentIDs []uint

	var err error
	if req.DeptId != 0 {
		studentIDs, err = s.orgRepo.GetStudentIDsByDepartmentID(req.DeptId)
	} else if req.ClassId != 0 {
		studentIDs, err = s.orgRepo.GetStudentIDsByClassID(req.ClassId)
	} else if req.StudentId != 0 {
		// Verify student exists
		user, uErr := s.userRepo.GetUserByID(req.StudentId)
		if uErr == nil && user != nil {
			studentIDs = []uint{user.ID}
		} else {
			err = uErr
		}
	}

	if err != nil || len(studentIDs) == 0 {
		return 0, errors.New("目标学生不存在或发生错误")
	}

	ding := model.Ding{
		LauncherID: launcherID, // From arg
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

	if err := s.dingRepo.CreateDing(&ding); err != nil {
		return 0, err
	}

	for _, studentID := range studentIDs {
		dingStudent := model.DingStudent{
			DingID:    ding.ID,
			StudentID: studentID,
			Status:    "pending",
			DingTime:  time.Now(),
		}
		if err := s.dingRepo.CreateDingStudent(&dingStudent); err != nil {
			return 0, err
		}

		notif := model.Notification{
			Title:      "新的打卡任务：" + req.Title,
			Content:    "请在规定时间内完成打卡任务。",
			SenderID:   launcherID,
			TargetType: "student",
			TargetID:   studentID,
		}
		// TODO: Refactor Notification logic into NotificationService
		if _, err := utils.PushNotification(notif, s.db); err != nil {
			// logging error but not failing the whole process?
			log.Printf("Failed to push notification to student %d: %v", studentID, err)
		} else {
			log.Printf("已向学生 %d 发送打卡任务通知", studentID)
		}
	}
	return ding.ID, nil
}

func (s *dingService) ListAllMyDings(studentID uint) (map[string][]model.Ding, error) {
	result := make(map[string][]model.Ding)

	pendingDings, err := s.dingRepo.GetDingsByStudentIDAndStatus(studentID, "pending")
	if err == nil {
		result["pending"] = pendingDings
	}
	completeDings, err := s.dingRepo.GetDingsByStudentIDAndStatus(studentID, "complete")
	if err == nil {
		result["complete"] = completeDings
	}
	return result, nil
}

func (s *dingService) ListMyCreatedDings(launcherID uint) ([]model.Ding, error) {
	return s.dingRepo.GetDingsByLauncherID(launcherID)
}

func (s *dingService) ListMyCreatedDingsRecords(userId uint, dingID string) (interface{}, interface{}) {
	// 查询该ding所有学生的状态
	dingRecords, err := s.dingRepo.GetDingRecordsByDingID(dingID) // []dingStudent
	if err != nil {
		return nil, err
	}
	return dingRecords, nil
}

func (s *dingService) ExportMyCreatedDingRecords(dingId string) (interface{}, interface{}) {
	// filePath,err
	dingsRecords, err := s.dingRepo.GetDingRecordsByDingID(dingId) // []dingStudent
	if err != nil {
		return nil, err
	}
	exportedFilePath, err := utils.ExportToExcel(dingsRecords, "ding_records_"+dingId+".xlsx")
	if err != nil {
		return nil, err
	}
	return exportedFilePath, nil
}

func (s *dingService) Ding(userId string, dingId uint) (interface{}, interface{}) {
	// 打卡逻辑
	dingStudent, err := s.dingRepo.UpdateDingStudent(dingId, userId)
	if err != nil {
		return nil, err
	}
	return dingStudent, nil
}

func (s *dingService) GetDingStats(launcherID uint) (map[string]int64, error) {
	total, checked, err := s.dingRepo.GetDingStats(launcherID)
	if err != nil {
		return nil, err
	}
	return map[string]int64{
		"total_count":   total,
		"checked_count": checked,
		"missed_count":  total - checked,
	}, nil
}
