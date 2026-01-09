package service

import (
	"errors"
	"log"
	"unihub/internal/model"
	"unihub/internal/repo"
	"unihub/internal/utils"
)

type OrgService interface {
	CreateDepartment(creatorID, roleID uint, name string) (*model.Department, error)
	CreateClass(creatorID, roleID uint, name string) (*model.Class, error)
	StudentJoinDepartment(studentID, roleID uint, inviteCode string) error
	StudentJoinClass(studentID, roleID uint, inviteCode string) error
	ListMyDepartments(counselorID, roleID uint) ([]model.Department, error)
	ListMyClasses(teacherID uint) ([]model.Class, error)
	ListMyClassStudent(userId uint, classId string) (interface{}, interface{})
	ListMyDepartmentStudent(userId uint, deptId string) (interface{}, interface{})
}

type orgService struct {
	orgRepo  repo.OrgRepository
	userRepo repo.UserRepository
}

func NewOrgService(orgRepo repo.OrgRepository, userRepo repo.UserRepository) OrgService {
	return &orgService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

func (s *orgService) CreateDepartment(creatorID, roleID uint, name string) (*model.Department, error) {
	if allowed, _ := s.userRepo.CheckPermission(roleID, "dept:create"); !allowed {
		return nil, errors.New("无权限创建部门")
	}

	inviteCode := utils.GenerateInviteCode()
	for {
		exists, _ := s.orgRepo.CheckDepartmentExists(inviteCode)
		if !exists {
			break
		}
		inviteCode = utils.GenerateInviteCode()
	}

	dept := model.Department{
		Name:        name,
		InviteCode:  inviteCode,
		CounselorID: creatorID,
	}

	if err := s.orgRepo.CreateDepartment(&dept); err != nil {
		return nil, err
	}
	return &dept, nil
}

func (s *orgService) CreateClass(creatorID, roleID uint, name string) (*model.Class, error) {
	if allowed, _ := s.userRepo.CheckPermission(roleID, "class:create"); !allowed {
		return nil, errors.New("无权限创建班级")
	}

	inviteCode := utils.GenerateInviteCode()
	for {
		exists, _ := s.orgRepo.CheckClassExists(inviteCode)
		if !exists {
			break
		}
		inviteCode = utils.GenerateInviteCode()
	}

	class := model.Class{
		Name:       name,
		InviteCode: inviteCode,
		TeacherID:  creatorID,
	}

	if err := s.orgRepo.CreateClass(&class); err != nil {
		return nil, err
	}
	return &class, nil
}

func (s *orgService) StudentJoinDepartment(studentID, roleID uint, inviteCode string) error {
	if allowed, _ := s.userRepo.CheckPermission(roleID, "dept:join"); !allowed {
		return errors.New("无权限加入部门")
	}

	dept, err := s.orgRepo.GetDepartmentByInviteCode(inviteCode)
	if err != nil {
		return errors.New("邀请码无效")
	}

	// Check if already joined THIS department
	count, _ := s.orgRepo.GetStudentDepartmentCount(studentID, dept.ID)
	if count > 0 {
		return errors.New("已加入该部门")
	}

	// Check if joined ANY department
	anyCount, _ := s.orgRepo.GetStudentAnyDepartmentCount(studentID)
	if anyCount > 0 {
		return errors.New("学生只能加入一个部门")
	}

	link := model.StudentDepartment{
		StudentID:    studentID,
		DepartmentID: dept.ID,
	}

	if err := s.orgRepo.AddStudentToDepartment(&link); err != nil {
		return err
	}

	log.Printf("Student %d joined Department %d", studentID, dept.ID)

	// Update user's department_id field for quick access
	user, err := s.userRepo.GetUserByID(studentID)
	if err != nil {
		return errors.New("更新学生部门信息失败")
	}
	user.DepartmentID = dept.ID
	if err := s.userRepo.UpdateUser(user); err != nil {
		return errors.New("更新学生部门信息失败")
	}

	return nil
}

func (s *orgService) StudentJoinClass(studentID, roleID uint, inviteCode string) error {
	if allowed, _ := s.userRepo.CheckPermission(roleID, "class:join"); !allowed {
		return errors.New("无权限加入班级")
	}

	class, err := s.orgRepo.GetClassByInviteCode(inviteCode)
	if err != nil {
		return errors.New("邀请码无效")
	}

	count, _ := s.orgRepo.GetStudentClassCount(studentID, class.ID)
	if count > 0 {
		return errors.New("已加入该班级")
	}

	link := model.StudentClass{
		StudentID: studentID,
		ClassID:   class.ID,
	}

	if err := s.orgRepo.AddStudentToClass(&link); err != nil {
		return err
	}

	return nil
}

func (s *orgService) ListMyDepartments(counselorID, roleID uint) ([]model.Department, error) {
	if allowed, _ := s.userRepo.CheckPermission(roleID, "dept:list"); !allowed {
		return nil, errors.New("无权限查看部门")
	}
	return s.orgRepo.ListDepartmentsByCounselorID(counselorID)
}

func (s *orgService) ListMyDepartmentStudent(userId uint, deptId string) (interface{}, interface{}) {
	// TODO : check Permission

	deptDetails, err := s.orgRepo.GetDepartmentDetailsByID(deptId) // interface{}
	if err != nil {
		return nil, err
	}
	students, err := s.orgRepo.ListStudentsByDepartmentID(deptId) // []
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	result["deptDetails"] = deptDetails
	result["students"] = students
	return result, nil
}

func (s *orgService) ListMyClasses(teacherID uint) ([]model.Class, error) {
	// Assuming teachers have permission to list their classes by default or checked earlier.
	// Adding simple permission check if needed, but original code didn't check permission explicitly for class listing
	// (Check handler/org.go: ListMyClasses doesn't call RequirePermission, unlike ListMyDepartments)
	return s.orgRepo.ListClassesByTeacherID(teacherID)
}

func (s *orgService) ListMyClassStudent(userId uint, classId string) (interface{}, interface{}) {
	// TODO : check Permission

	classDetails, err := s.orgRepo.GetClassDetailsByID(classId) // interface{}
	if err != nil {
		return nil, err
	}
	students, err := s.orgRepo.ListStudentsByClassID(classId) // []
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	result["classDetails"] = classDetails
	result["students"] = students
	return result, nil
}
