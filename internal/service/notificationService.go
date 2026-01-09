package service

import (
	"errors"
	"unihub/internal/model"
	"unihub/internal/repo"
	"unihub/internal/utils"

	"gorm.io/gorm"
)

type CreateNotifRequest struct {
	Title      string
	Content    string
	TargetType string
	TargetID   uint
	SenderID   uint
}

type NotificationService interface {
	Create(req CreateNotifRequest) error
	GetMyNotifications(studentID uint) ([]model.Notification, error)
}

type notificationService struct {
	notifRepo repo.NotificationRepository
	orgRepo   repo.OrgRepository
	userRepo  repo.UserRepository
	db        *gorm.DB
}

func NewNotificationService(notifRepo repo.NotificationRepository, orgRepo repo.OrgRepository, userRepo repo.UserRepository, db *gorm.DB) NotificationService {
	return &notificationService{
		notifRepo: notifRepo,
		orgRepo:   orgRepo,
		userRepo:  userRepo,
		db:        db,
	}
}

func (s *notificationService) Create(req CreateNotifRequest) error {
	// Permission Check inside Service
	// Verify sender is counselor of dept or teacher of class
	hasPerm := false
	if req.TargetType == "dept" {
		depts, err := s.orgRepo.ListDepartmentsByCounselorID(req.SenderID)
		if err == nil {
			for _, d := range depts {
				if d.ID == req.TargetID {
					hasPerm = true
					break
				}
			}
		}
	} else if req.TargetType == "class" {
		classes, err := s.orgRepo.ListClassesByTeacherID(req.SenderID)
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
		return errors.New("没有权限向该目标发送通知")
	}

	notif := model.Notification{
		Title:      req.Title,
		Content:    req.Content,
		SenderID:   req.SenderID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	if err := s.notifRepo.CreateNotification(&notif); err != nil {
		return err
	}

	if _, err := utils.PushNotification(notif, s.db); err != nil {
		return err
	}

	return nil
}

func (s *notificationService) GetMyNotifications(studentID uint) ([]model.Notification, error) {
	targets := []model.Target{}

	// Get Student's Department
	deptID, err := s.orgRepo.GetStudentDepartmentID(studentID)
	if err == nil && deptID != 0 {
		targets = append(targets, model.Target{Type: "dept", ID: deptID})
	}

	// Get Student's Classes
	classIDs, err := s.orgRepo.GetStudentClassIDs(studentID)
	if err == nil {
		for _, cid := range classIDs {
			targets = append(targets, model.Target{Type: "class", ID: cid})
		}
	}

	return s.notifRepo.GetNotificationsForTargets(targets)
}
