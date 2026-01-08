package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"unihub/internal/model"
)

// RequirePermission checks if role has permission code.
func RequirePermission(ctx context.Context, db *gorm.DB, roleID uint, perm string) (bool, error) {
	var count int64
	err := db.Table("permissions").Select("permissions.id").
		Joins("join role_permissions rp on rp.permission_id = permissions.id").
		Where("rp.role_id = ? AND permissions.code = ?", roleID, perm).
		Count(&count).Error
	return count > 0, err
}

// DataScopeFilter returns a list of org IDs accessible given role data scope and org path.
func DataScopeFilter(db *gorm.DB, role model.Role, orgID *uint) ([]uint, error) {
	// Simple implementation: support all, dept, dept_and_sub, self
	switch role.DataScope {
	case "all":
		return nil, nil
	case "dept":
		if orgID == nil {
			return []uint{}, nil
		}
		return []uint{*orgID}, nil
	case "dept_and_sub":
		if orgID == nil {
			return []uint{}, nil
		}
		var ids []uint
		// match path prefix
		if err := db.Model(&model.OrgUnit{}).
			Where("path LIKE ?", "%/"+toStr(*orgID)+"/%").
			Or("id = ?", *orgID).
			Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		return ids, nil
	case "self":
		return []uint{}, nil
	default:
		return []uint{}, nil
	}
}

func toStr(id uint) string {
	return fmt.Sprintf("%d", id)
}
