package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unihub/internal/model"
)

type UserHandler struct {
	DB *gorm.DB
}

// GetProfile 获取个人资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	var user model.User
	// 预加载角色信息
	if err := h.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 隐藏密码
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// ListStudents 列出自己管理的部门或班级的学生
// 逻辑：
// 1. 如果是辅导员，查找自己创建的所有部门，再查找这些部门下的学生。
// 2. 如果是教师，查找自己创建的所有班级，再查找这些班级下的学生。
func (h *UserHandler) ListStudents(c *gin.Context) {
	userID := c.GetUint("userID")
	roleID := c.GetUint("roleID")

	var role model.Role
	if err := h.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "角色错误"})
		return
	}

	var students []model.User

	if role.Key == "counselor" {
		// 查找该辅导员管理的部门
		var deptIDs []uint
		h.DB.Model(&model.Department{}).Where("counselor_id = ?", userID).Pluck("id", &deptIDs)

		if len(deptIDs) == 0 {
			c.JSON(http.StatusOK, []model.User{})
			return
		}

		// 查找部门下的学生
		h.DB.Table("users").
			Joins("JOIN student_departments sd ON sd.student_id = users.id").
			Where("sd.department_id IN ?", deptIDs).
			Find(&students)

	} else if role.Key == "teacher" {
		// 查找该教师管理的班级
		var classIDs []uint
		h.DB.Model(&model.Class{}).Where("teacher_id = ?", userID).Pluck("id", &classIDs)

		if len(classIDs) == 0 {
			c.JSON(http.StatusOK, []model.User{})
			return
		}

		// 查找班级下的学生
		h.DB.Table("users").
			Joins("JOIN student_classes sc ON sc.student_id = users.id").
			Where("sc.class_id IN ?", classIDs).
			Find(&students)

	} else if role.Key == "admin" || role.Key == "super_admin" {
		// 管理员可以看到所有学生
		h.DB.Joins("JOIN roles r ON r.id = users.role_id").
			Where("r.key = ?", "student").
			Find(&students)
	} else {
		// 学生或其他角色无权查看
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看"})
		return
	}

	// 隐藏密码
	for i := range students {
		students[i].Password = ""
	}

	c.JSON(http.StatusOK, students)
}
