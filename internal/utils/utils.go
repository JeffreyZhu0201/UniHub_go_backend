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

// GenerateInviteCode 生成8位随机大写字母邀请码
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
	// find target users
	var studentIDs []uint
	if notification.TargetType == "dept" {
		// find users in department
		if err := DB.Model(&model.StudentDepartment{}).Where("department_id = ?", notification.TargetID).Pluck("student_id", &studentIDs).Error; err != nil {
			return "查询部门学生失败", err
		}
	} else if notification.TargetType == "class" {
		// find users in class
		if err := DB.Model(&model.StudentClass{}).Where("class_id = ?", notification.TargetID).Pluck("student_id", &studentIDs).Error; err != nil {
			return "查询班级学生失败", err
		}
	} else {
		return "未知目标类型", nil
	}

	if len(studentIDs) == 0 {
		return "未找到目标学生", nil
	}

	// push to each student
	// for _, sid := range studentIDs {
	// 	// TODO: lookup user push token and send
	// }

	return "消息推送成功", nil
}
