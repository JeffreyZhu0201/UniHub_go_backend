package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

func GetDepartmentByInviteCode(db *gorm.DB, inviteCode string) (model.Department, error) {

	var dept model.Department
	if err := db.Where("invite_code = ?", inviteCode).First(&dept).Error; err != nil {
		return dept, err
	}
	return dept, nil
}
