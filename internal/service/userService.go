package service

import (
	"errors"
	"unihub/internal/model"
	"unihub/internal/repo"
)

type UserService interface {
	// Register(req RegisterRequest) (*model.User, error)
	// Login(req LoginRequest) (string, *model.User, error)
	GetProfile(userID uint) (*model.User, error)
	GetUserOrgInfo(userID uint) (map[string]interface{}, error) // 新增接口方法
	ListStudents(userID, roleID uint) ([]model.User, error)
}

type userService struct {
	userRepo repo.UserRepository
	orgRepo  repo.OrgRepository
}

func NewUserService(userRepo repo.UserRepository, orgRepo repo.OrgRepository) UserService {
	return &userService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
	}
}

func (s *userService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetUserByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	user.Password = ""
	return user, nil
}

// 新增实现：获取用户组织信息
func (s *userService) GetUserOrgInfo(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	var classes []model.Class
	// 只有学生角色才去查询加入的班级
	if user.RoleID != 0 { // 假设 RoleID 存在，简单判断，严谨可判断 Role Key
		// 这里简单处理，查询该用户作为学生加入的班级
		classes, _ = s.orgRepo.ListClassesByStudentID(userID)
	}

	var dept *model.Department
	if user.DepartmentID != 0 {
		dept, _ = s.orgRepo.GetDepartmentByID(user.DepartmentID)
	}

	// 组装返回数据
	result := map[string]interface{}{
		"department": dept,
		"classes":    classes,
	}
	return result, nil

}

func (s *userService) ListStudents(userID, roleID uint) ([]model.User, error) {
	// We need the role key. We can get it from user with role.
	user, err := s.userRepo.GetUserByIDWithRole(userID)
	if err != nil {
		return nil, err
	}

	roleKey := user.Role.Key
	var students []model.User

	if roleKey == "counselor" {
		depts, err := s.orgRepo.ListDepartmentsByCounselorID(userID)
		if err != nil {
			return nil, err
		}
		var deptIDs []uint
		for _, d := range depts {
			deptIDs = append(deptIDs, d.ID)
		}
		if len(deptIDs) > 0 {
			students, err = s.userRepo.ListStudentsByDepartmentIDs(deptIDs)
			if err != nil {
				return nil, err
			}
		}
	} else if roleKey == "teacher" {
		classes, err := s.orgRepo.ListClassesByTeacherID(userID) // 找出所有该教师负责的班级
		if err != nil {
			return nil, err
		}
		var classIDs []uint
		for _, c := range classes {
			classIDs = append(classIDs, c.ID)
		}
		if len(classIDs) > 0 {
			students, err = s.userRepo.ListStudentsByClassIDs(classIDs) // 返回这些班级所有学生
			if err != nil {
				return nil, err
			}
		}
	} else if roleKey == "admin" || roleKey == "super_admin" {
		students, err = s.userRepo.ListAllStudents()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("无权查看")
	}

	for i := range students {
		students[i].Password = ""
	}

	return students, nil
}
