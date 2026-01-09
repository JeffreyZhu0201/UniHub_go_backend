package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type OrgRepository interface {
	GetDepartmentByInviteCode(inviteCode string) (*model.Department, error)
	GetClassByInviteCode(inviteCode string) (*model.Class, error)
	CreateDepartment(dept *model.Department) error
	CreateClass(class *model.Class) error
	CheckDepartmentExists(inviteCode string) (bool, error)
	CheckClassExists(inviteCode string) (bool, error)
	AddStudentToDepartment(link *model.StudentDepartment) error
	AddStudentToClass(link *model.StudentClass) error
	GetStudentDepartmentCount(studentID, deptID uint) (int64, error)
	GetStudentClassCount(studentID, classID uint) (int64, error)
	GetStudentAnyDepartmentCount(studentID uint) (int64, error)
	ListDepartmentsByCounselorID(counselorID uint) ([]model.Department, error)
	ListClassesByTeacherID(teacherID uint) ([]model.Class, error)
	GetStudentIDsByDepartmentID(deptID uint) ([]uint, error)
	GetStudentIDsByClassID(classID uint) ([]uint, error)
	GetStudentDepartmentID(studentID uint) (uint, error)
	GetStudentClassIDs(studentID uint) ([]uint, error)
}

type orgRepository struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

func (r *orgRepository) GetDepartmentByInviteCode(inviteCode string) (*model.Department, error) {
	var dept model.Department
	if err := r.db.Where("invite_code = ?", inviteCode).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *orgRepository) GetClassByInviteCode(inviteCode string) (*model.Class, error) {
	var class model.Class
	if err := r.db.Where("invite_code = ?", inviteCode).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *orgRepository) CreateDepartment(dept *model.Department) error {
	return r.db.Create(dept).Error
}

func (r *orgRepository) CreateClass(class *model.Class) error {
	return r.db.Create(class).Error
}

func (r *orgRepository) CheckDepartmentExists(inviteCode string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Department{}).Where("invite_code = ?", inviteCode).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *orgRepository) CheckClassExists(inviteCode string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Class{}).Where("invite_code = ?", inviteCode).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *orgRepository) AddStudentToDepartment(link *model.StudentDepartment) error {
	return r.db.Create(link).Error
}

func (r *orgRepository) AddStudentToClass(link *model.StudentClass) error {
	return r.db.Create(link).Error
}

func (r *orgRepository) GetStudentDepartmentCount(studentID, deptID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.StudentDepartment{}).Where("student_id = ? AND department_id = ?", studentID, deptID).Count(&count).Error
	return count, err
}

func (r *orgRepository) GetStudentClassCount(studentID, classID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.StudentClass{}).Where("student_id = ? AND class_id = ?", studentID, classID).Count(&count).Error
	return count, err
}

func (r *orgRepository) GetStudentAnyDepartmentCount(studentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.StudentDepartment{}).Where("student_id = ?", studentID).Count(&count).Error
	return count, err
}

func (r *orgRepository) ListDepartmentsByCounselorID(counselorID uint) ([]model.Department, error) {
	var depts []model.Department
	err := r.db.Where("counselor_id = ?", counselorID).Find(&depts).Error
	return depts, err
}

func (r *orgRepository) ListClassesByTeacherID(teacherID uint) ([]model.Class, error) {
	var classes []model.Class
	err := r.db.Where("teacher_id = ?", teacherID).Find(&classes).Error
	return classes, err
}

func (r *orgRepository) GetStudentIDsByDepartmentID(deptID uint) ([]uint, error) {
	var studentIDs []uint
	err := r.db.Model(&model.StudentDepartment{}).Where("department_id = ?", deptID).Pluck("student_id", &studentIDs).Error
	return studentIDs, err
}

func (r *orgRepository) GetStudentIDsByClassID(classID uint) ([]uint, error) {
	var studentIDs []uint
	err := r.db.Model(&model.StudentClass{}).Where("class_id = ?", classID).Pluck("student_id", &studentIDs).Error
	return studentIDs, err
}

func (r *orgRepository) GetStudentDepartmentID(studentID uint) (uint, error) {
	var deptID uint
	err := r.db.Model(&model.StudentDepartment{}).Where("student_id = ?", studentID).Pluck("department_id", &deptID).Error
	return deptID, err
}

func (r *orgRepository) GetStudentClassIDs(studentID uint) ([]uint, error) {
	var classIDs []uint
	err := r.db.Model(&model.StudentClass{}).Where("student_id = ?", studentID).Pluck("class_id", &classIDs).Error
	return classIDs, err
}
