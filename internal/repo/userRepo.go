package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetRoleByKey(key string) (*model.Role, error)
	UpdateUser(user *model.User) error
	CheckPermission(roleID uint, permCode string) (bool, error)
	GetUserByIDWithRole(id uint) (*model.User, error)
	ListStudentsByDepartmentIDs(deptIDs []uint) ([]model.User, error)
	ListStudentsByClassIDs(classIDs []uint) ([]model.User, error)
	ListAllStudents() ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetRoleByKey(key string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("`key` = ?", key).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) CheckPermission(roleID uint, permCode string) (bool, error) {
	var count int64
	err := r.db.Table("permissions").Select("permissions.id").
		Joins("join role_permissions rp on rp.permission_id = permissions.id").
		Where("rp.role_id = ? AND permissions.code = ?", roleID, permCode).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) GetUserByIDWithRole(id uint) (*model.User, error) {
	// 通过UserId查询用户信息，并预加载角色信息
	var user model.User
	if err := r.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ListStudentsByDepartmentIDs(deptIDs []uint) ([]model.User, error) {
	var students []model.User
	err := r.db.Table("users").
		Joins("JOIN student_departments sd ON sd.student_id = users.id").
		Where("sd.department_id IN ?", deptIDs).
		Find(&students).Error
	return students, err
}

func (r *userRepository) ListStudentsByClassIDs(classIDs []uint) ([]model.User, error) {
	var students []model.User
	err := r.db.Table("users").
		Joins("JOIN student_classes sc ON sc.student_id = users.id").
		Where("sc.class_id IN ?", classIDs).
		Find(&students).Error
	return students, err
}

func (r *userRepository) ListAllStudents() ([]model.User, error) {
	var students []model.User
	err := r.db.Joins("JOIN roles r ON r.id = users.role_id").
		Where("r.key = ?", "student").
		Find(&students).Error
	return students, err
}
