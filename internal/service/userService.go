package service

import (
	"errors"
	"unihub/internal/model"
	"unihub/internal/repo"
)

type UserService interface {
	GetProfile(userID uint) (*model.User, error)
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
