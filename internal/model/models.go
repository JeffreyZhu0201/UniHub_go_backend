package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents system roles (super admin, school admin, college admin, counselor, homeroom teacher, student).
type Role struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:50;uniqueIndex"`
	Key       string `gorm:"size:50;uniqueIndex"`
	DataScope string `gorm:"size:20"` // all, custom, dept, dept_and_sub, self
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Permission represents an action string, e.g., student:list
type Permission struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"size:100;uniqueIndex"`
	Name      string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// RolePermission join table
// many-to-many
// gorm will auto create table role_permissions

// OrgUnit represents hierarchical units (school/college/grade/class).
type OrgUnit struct {
	ID        uint      `gorm:"primaryKey"`
	UUID      uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Name      string    `gorm:"size:100;not null"`
	Type      string    `gorm:"size:30;not null"` // school, college, grade, class
	ParentID  *uint     `gorm:"index"`
	Path      string    `gorm:"size:255;index"` // materialized path, e.g., /1/5/9
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User represents login identity.
type User struct {
	ID        uint    `gorm:"primaryKey"`
	Nickname  string  `gorm:"size:100;not null"`
	Email     string  `gorm:"size:120;uniqueIndex;not null"`
	Password  string  `gorm:"size:255;not null"`
	RoleID    uint    `gorm:"index"`
	OrgUnitID *uint   `gorm:"index"`
	StaffNo   *string `gorm:"size:50"`  // for admins/teachers/counselors
	StudentNo *string `gorm:"size:50"`  // for students
	PushToken string  `gorm:"size:255"` // for push notifications
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Role Role    `gorm:"foreignKey:RoleID"`
	Org  OrgUnit `gorm:"foreignKey:OrgUnitID"`
}

// RolePermission defines many-to-many link.
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

// Department represents counselor-managed department with invite code.
type Department struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	InviteCode  string `gorm:"size:12;uniqueIndex;not null"`
	CounselorID uint   `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Class represents teacher-managed class with invite code.
type Class struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:100;not null"`
	InviteCode string `gorm:"size:12;uniqueIndex;not null"`
	TeacherID  uint   `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// StudentDepartment maps student to a single department.
type StudentDepartment struct {
	ID           uint `gorm:"primaryKey"`
	StudentID    uint `gorm:"index;not null"`
	DepartmentID uint `gorm:"index;not null"`
	CreatedAt    time.Time
}

// StudentClass maps student to multiple classes.
type StudentClass struct {
	ID        uint `gorm:"primaryKey"`
	StudentID uint `gorm:"index;not null"`
	ClassID   uint `gorm:"index;not null"`
	CreatedAt time.Time
}

// Developer represents open platform developer.
type Developer struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:120;uniqueIndex;not null"`
	Secret    string `gorm:"size:64;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// App represents registered application with rate limit.
type App struct {
	ID          uint   `gorm:"primaryKey"`
	DeveloperID uint   `gorm:"index;not null"`
	Name        string `gorm:"size:100;not null"`
	AppID       string `gorm:"size:32;uniqueIndex;not null"`
	AppSecret   string `gorm:"size:64;not null"`
	RateLimit   int    `gorm:"not null"` // requests per minute
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Notification 通知
type Notification struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"size:100;not null"`  // 标题
	Content    string `gorm:"type:text;not null"` // 内容
	SenderID   uint   `gorm:"index"`              // 发送者ID
	TargetType string `gorm:"size:20;not null"`   // dept 或 class
	TargetID   uint   `gorm:"index"`              // 目标部门或班级ID
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// LeaveRequest 请假申请
type LeaveRequest struct {
	ID        uint      `gorm:"primaryKey"`
	StudentID uint      `gorm:"index;not null"`
	Type      string    `gorm:"size:50"` // 病假、事假等
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Reason    string    `gorm:"size:255"`
	Status    string    `gorm:"size:20;default:'pending'"` // pending, approved, rejected, active, completed, overdue
	AuditorID *uint     `gorm:"index"`                     // 审批人(辅导员)
	AuditTime *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Task 任务 (签到/查寝)
type Task struct {
	ID          uint      `gorm:"primaryKey"`
	UUID        uuid.UUID `gorm:"type:char(36);uniqueIndex"` // Added UUID
	Title       string    `gorm:"size:100;not null"`
	Type        string    `gorm:"size:50;not null"` // sign_in, dorm_check, leave_check
	Description string    `gorm:"size:255"`
	CreatorID   uint      `gorm:"index"`   // 发布者
	TargetType  string    `gorm:"size:20"` // dept, class
	TargetID    uint      `gorm:"index"`
	Deadline    time.Time `gorm:"not null"`
	Config      string    `gorm:"type:json"` // 任务配置：如签到经纬度、距离限制等
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.UUID = uuid.New()
	return
}

// TaskRecord 任务记录 (学生提交)
type TaskRecord struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"index;not null"`
	StudentID uint      `gorm:"index;not null"`
	Status    string    `gorm:"size:20"`   // completed, late
	Data      string    `gorm:"type:json"` // 提交的数据：location, photo_url
	CreatedAt time.Time // 提交时间
}

// 打卡任务实体
type Ding struct {
	ID         uint   `gorm:"primaryKey"`
	LauncherID uint   `gorm:"index;not null"`
	Title      string `gorm:"size:100;not null"`
	StartTime  time.Time
	EndTime    time.Time
	// 经纬度
	Latitude  float64 `gorm:"not null"`
	Longitude float64 `gorm:"not null"`
	// 允许的最大距离，单位米
	Radius    float64 `gorm:"not null"`
	UserID    uint    `gorm:"index;not null"`
	DeptID    uint    `gorm:"index;not null"`
	ClassID   uint    `gorm:"index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DingStudent struct {
	ID            uint `gorm:"primaryKey"`
	DingID        uint `gorm:"index;not null"`
	StudentID     uint `gorm:"index;not null"`
	DingTime      time.Time
	DingLatitude  float64 `gorm:"not null"`
	DingLongitude float64 `gorm:"not null"`
	Status        string  `gorm:"size:20"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// AutoMigrate migrates all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Role{}, &Permission{}, &OrgUnit{}, &User{}, &RolePermission{},
		&Department{}, &Class{}, &StudentDepartment{}, &StudentClass{},
		&Developer{}, &App{},
		&Notification{}, &LeaveRequest{}, &Task{}, &TaskRecord{},
		&Ding{}, &DingStudent{},
	)
}
