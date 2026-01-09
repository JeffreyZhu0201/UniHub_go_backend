package service

import (
	"errors"
	"time"
	"unihub/internal/model"
	"unihub/internal/repo"
)

type ApplyLeaveRequest struct {
	StudentID uint
	Type      string
	StartTime time.Time
	EndTime   time.Time
	Reason    string
}

type AuditLeaveRequest struct {
	AuditorID uint
	RoleID    uint
	LeaveID   uint
	Status    string
}

type LeaveService interface {
	Apply(req ApplyLeaveRequest) (*model.LeaveRequest, error)
	Audit(req AuditLeaveRequest) error
	ListPendingLeaves(counselorID, roleID uint) ([]model.LeaveRequest, error)
	MyLeaves(studentID uint) ([]model.LeaveRequest, error)
}

type leaveService struct {
	leaveRepo repo.LeaveRepository
	orgRepo   repo.OrgRepository
	userRepo  repo.UserRepository
}

func NewLeaveService(leaveRepo repo.LeaveRepository, orgRepo repo.OrgRepository, userRepo repo.UserRepository) LeaveService {
	return &leaveService{
		leaveRepo: leaveRepo,
		orgRepo:   orgRepo,
		userRepo:  userRepo,
	}
}

func (s *leaveService) Apply(req ApplyLeaveRequest) (*model.LeaveRequest, error) {
	leave := model.LeaveRequest{
		StudentID: req.StudentID,
		Type:      req.Type,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Reason:    req.Reason,
		Status:    "pending",
	}

	if err := s.leaveRepo.CreateLeaveRequest(&leave); err != nil {
		return nil, err
	}
	return &leave, nil
}

func (s *leaveService) Audit(req AuditLeaveRequest) error {
	if allowed, _ := s.userRepo.CheckPermission(req.RoleID, "leave:approve"); !allowed {
		return errors.New("无权限审批")
	}

	leave, err := s.leaveRepo.GetLeaveRequestByID(req.LeaveID)
	if err != nil {
		return errors.New("请假记录不存在")
	}

	// Validation: Is Auditor the student's counselor?
	studentDeptID, err := s.orgRepo.GetStudentDepartmentID(leave.StudentID)
	if err != nil || studentDeptID == 0 {
		return errors.New("学生未加入部门")
	}

	depts, err := s.orgRepo.ListDepartmentsByCounselorID(req.AuditorID)
	// Check if studentDeptID is in depts
	isManaged := false
	if err == nil {
		for _, d := range depts {
			if d.ID == studentDeptID {
				isManaged = true
				break
			}
		}
	}

	if !isManaged {
		return errors.New("无权审批该学生请假")
	}

	now := time.Now()
	leave.Status = req.Status
	leave.AuditorID = &req.AuditorID
	leave.AuditTime = &now

	if err := s.leaveRepo.UpdateLeaveRequest(leave); err != nil {
		return err
	}

	// Create "Sign In Task" if approved?
	// Requirement: "请假成功后自动生成一个签到任务，规定时间段在请假截止时间前。"
	if req.Status == "approved" {
		// Implementation of creating auto sign-in task
		// This requires TaskService or direct repo access.
		// Since services shouldn't cyclically depend, we might use TaskRepo here or Event Bus.
		// Or simpler: just do it here.
		// "规定时间段在请假截止时间前" -> Deadline = leave.EndTime ?
		// "Sign In Task" usually means strict location check.
		// Or just a task for them to click "I'm back".
		// TODO: Implement automatic task creation
	}

	return nil
}

func (s *leaveService) ListPendingLeaves(counselorID, roleID uint) ([]model.LeaveRequest, error) {
	if allowed, _ := s.userRepo.CheckPermission(roleID, "leave:approve"); !allowed {
		return nil, errors.New("无权限查看待审批请假")
	}

	// Find all students in counselor's departments
	depts, err := s.orgRepo.ListDepartmentsByCounselorID(counselorID)
	if err != nil {
		return nil, err
	}

	var allStudentIDs []uint
	for _, d := range depts {
		ids, err := s.orgRepo.GetStudentIDsByDepartmentID(d.ID)
		if err == nil {
			allStudentIDs = append(allStudentIDs, ids...)
		}
	}

	if len(allStudentIDs) == 0 {
		return []model.LeaveRequest{}, nil
	}

	return s.leaveRepo.ListPendingLeavesByStudentIDs(allStudentIDs)
}

func (s *leaveService) MyLeaves(studentID uint) ([]model.LeaveRequest, error) {
	return s.leaveRepo.ListLeavesByStudentID(studentID)
}
