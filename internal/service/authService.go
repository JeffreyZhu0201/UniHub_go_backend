package service

import (
	"errors"
	"unihub/internal/config"
	"unihub/internal/model"
	"unihub/internal/repo"
	"unihub/internal/utils"
	"unihub/pkg/jwtutil"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Nickname   string  `json:"nickname"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	RoleKey    string  `json:"role_key"`
	StaffNo    *string `json:"staff_no"`
	StudentNo  *string `json:"student_no"`
	InviteCode string  `json:"invite_code"`
	PushToken  string  `json:"push_token"`
}

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	PushToken string `json:"push_token"`
}

type AuthService interface {
	Register(req RegisterRequest) (string, error)
	Login(req LoginRequest) (string, error)
}

type authService struct {
	userRepo repo.UserRepository
	orgRepo  repo.OrgRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repo.UserRepository, orgRepo repo.OrgRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(req RegisterRequest) (string, error) {
	var department *model.Department

	// Check invite code
	if req.InviteCode != "" {
		dept, err := s.orgRepo.GetDepartmentByInviteCode(req.InviteCode)
		if err != nil {
			return "", errors.New("无效的邀请码")
		}
		department = dept
	}

	// Validate Email
	if len(req.Email) < 8 && !utils.EndsWith(req.Email, ".com") {
		return "", errors.New("邮箱长度必须至少为8位")
	}
	if !utils.EndsWith(req.Email, ".com") {
		return "", errors.New("邮箱必须以 .com 结尾")
	}

	// Validate Role
	role, err := s.userRepo.GetRoleByKey(req.RoleKey)
	if err != nil {
		return "", errors.New("无效的角色")
	}

	// Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	deptID := uint(0)
	if department != nil {
		deptID = department.ID
	}

	user := model.User{
		Nickname:     req.Nickname,
		Email:        req.Email,
		DepartmentID: deptID, // This might be redundant if we use StudentDepartment table, but keeping as per original code
		Password:     string(hashed),
		RoleID:       role.ID,
		StaffNo:      req.StaffNo,
		StudentNo:    req.StudentNo,
		PushToken:    req.PushToken,
	}

	if err := s.userRepo.CreateUser(&user); err != nil {
		return "", err
	}

	if department != nil {
		studentDept := model.StudentDepartment{
			StudentID:    user.ID,
			DepartmentID: department.ID,
		}
		if err := s.orgRepo.AddStudentToDepartment(&studentDept); err != nil {
			return "", err
		}
	}

	// Generate Token
	token, err := jwtutil.Generate(s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours, user.ID, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Login(req LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("邮箱或密码错误")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return "", errors.New("邮箱或密码错误")
	}

	// Update Push Token
	if req.PushToken != "" && req.PushToken != user.PushToken {
		user.PushToken = req.PushToken
		if err := s.userRepo.UpdateUser(user); err != nil {
			// Log error but proceed?
		}
	}

	token, err := jwtutil.Generate(s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours, user.ID, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}
