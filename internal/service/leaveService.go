package service

import (
	"errors"
	"time"
	"unihub/internal/DTO"
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
	Audit(req AuditLeaveRequest, d DingService) error
	ListPendingLeaves(counselorID, roleID uint) ([]model.LeaveRequest, error)
	MyLeaves(studentID uint) ([]model.LeaveRequest, error)
	LeaveData(userId uint) (interface{}, interface{})
	LeaveBackInfo(userId uint) (interface{}, interface{})
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

func (s *leaveService) Audit(req AuditLeaveRequest, dscv DingService) error {
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

	if req.Status == "approved" {
		dingEntity := DTO.CreateDingRequest{
			StudentId:  leave.StudentID,
			Title:      "返校签到",
			LauncherId: req.AuditorID,
			StartTime:  leave.EndTime.Add(-1 * time.Hour),
			EndTime:    leave.EndTime,
			Latitude:   200,
			Longitude:  200,
			Radius:     50,
		}
		dingId, err := dscv.CreateDing(dingEntity, req.AuditorID, req.RoleID)
		if err != nil {
			return err
		}
		leave.DingId = dingId
		if err := s.leaveRepo.UpdateLeaveRequest(leave); err != nil {
			return err
		}
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

func (s *leaveService) LeaveData(userId uint) (interface{}, interface{}) {
	//approvedCount, _ := s.leaveRepo.CountLeavesByStatus("approved")
	// get all my students
	students, _ := s.orgRepo.ListStudentsByCounselorID(userId)
	// select id
	var studentIDs []uint
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID)
	}
	// find leaves by student IDs and status
	result := make(map[string][]interface{})
	for _, status := range []string{"approved", "pending", "rejected"} {
		leavesDetail, _ := s.leaveRepo.ListLeavesWithStudentsByStudentsAndStatus(studentIDs, status)
		result[status] = leavesDetail
	}
	return result, nil
}

func (s *leaveService) LeaveBackInfo(userId uint) (interface{}, interface{}) {
	// get all my students
	students, _ := s.orgRepo.ListStudentsByCounselorID(userId)
	// select id
	var studentIDs []uint
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID)
	}
	// find leaves by student IDs and status
	result := make(map[string]interface{})
	allApproved, _ := s.leaveRepo.ListApprovedLeavesWithStudentsByStudents(students)
	returned, _ := s.leaveRepo.ListLeavesWithStudentsByStudentsByDingStatusBeforeEnd(students)
	lateReturned, _ := s.leaveRepo.ListLeavesWithStudentsByStudentsAfterEnd(students)
	result["approved"] = allApproved
	result["returned"] = returned
	result["late_returned"] = lateReturned
	return result, nil
}
