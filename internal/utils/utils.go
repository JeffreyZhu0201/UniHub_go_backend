package utils

import (
	"crypto/rand"
	"math/big"
	"unihub/internal/model"

	"gorm.io/gorm"
)

func EndsWith(username string, s string) bool {
	if len(username) < len(s) {
		return false
	}
	return username[len(username)-len(s):] == s
}

// generateInviteCode 生成8位随机大写字母邀请码
func GenerateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 8)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func PushNotification(notification model.Notification, DB *gorm.DB) (string, error) {
	//notif := model.Notification{
	//	Title:      req.Title,
	//	Content:    req.Content,
	//	SenderID:   userID,
	//	TargetType: req.TargetType,
	//	TargetID:   req.TargetID,
	//}

	// find target users
	// send notification to target users
	var studentIDs []uint
	if notification.TargetType == "dept" {
		// find users in department
		if err := DB.Where("department_id = ?", notification.TargetID).Pluck("student_id", &studentIDs).Error; err != nil {
			return "未找到目标学生", err
		}
	}
	if notification.TargetType == "class" {
		// find users in class
		if err := DB.Where("class_id = ?", notification.TargetID).Pluck("student_id", &studentIDs).Error; err != nil {
			return "未找到目标学生", err
		}
	}
	// push to each student
	//for _, sid := range studentIDs
	return "消息推送成功", nil
}
