package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(notif *model.Notification) error
	GetNotifications(targetType string, targetID uint) ([]model.Notification, error)
	GetNotificationsForTargets(targets []model.Target) ([]model.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(notif *model.Notification) error {
	return r.db.Create(notif).Error
}

func (r *notificationRepository) GetNotifications(targetType string, targetID uint) ([]model.Notification, error) {
	var notifs []model.Notification
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).Order("created_at desc").Find(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) GetNotificationsForTargets(targets []model.Target) ([]model.Notification, error) {
	if len(targets) == 0 {
		return []model.Notification{}, nil
	}

	query := r.db.Model(&model.Notification{})
	for i, t := range targets {
		if i == 0 {
			query = query.Where("target_type = ? AND target_id = ?", t.Type, t.ID)
		} else {
			query = query.Or("target_type = ? AND target_id = ?", t.Type, t.ID)
		}
	}

	var notifs []model.Notification
	err := query.Order("created_at desc").Find(&notifs).Error
	return notifs, err
}
