package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unihub/internal/config"
	"unihub/internal/handler"
	"unihub/internal/repo"
	"unihub/internal/service"
	"unihub/pkg/middleware"
)

// Register registers all routes.
func Register(r *gin.Engine, cfg *config.Config, db *gorm.DB) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 初始化 Repositories
	userRepo := repo.NewUserRepository(db)
	orgRepo := repo.NewOrgRepository(db)
	notifRepo := repo.NewNotificationRepository(db)
	leaveRepo := repo.NewLeaveRepository(db)
	taskRepo := repo.NewTaskRepository(db)
	openRepo := repo.NewOpenRepository(db)
	dingRepo := repo.NewDingRepository(db)

	// 初始化 Services
	authSvc := service.NewAuthService(userRepo, orgRepo, cfg)
	orgSvc := service.NewOrgService(orgRepo, userRepo)
	userSvc := service.NewUserService(userRepo, orgRepo)
	notifSvc := service.NewNotificationService(notifRepo, orgRepo, userRepo, db)
	leaveSvc := service.NewLeaveService(leaveRepo, orgRepo, userRepo)
	taskSvc := service.NewTaskService(taskRepo, orgRepo, userRepo)
	openSvc := service.NewOpenService(openRepo)
	dingSvc := service.NewDingService(dingRepo, orgRepo, userRepo, db)

	// 初始化 Handlers
	authH := handler.NewAuthHandler(authSvc)
	orgH := handler.NewOrgHandler(orgSvc)
	userH := handler.NewUserHandler(userSvc)
	notifH := handler.NewNotificationHandler(notifSvc)
	leaveH := handler.NewLeaveHandler(leaveSvc)
	taskH := handler.NewTaskHandler(taskSvc)
	openH := handler.NewOpenHandler(openSvc)
	dingH := handler.NewDingHandler(dingSvc, userRepo)

	api := r.Group("/api/v1")
	{
		// 认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}

		// 受保护路由 (Internal Users)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/user/profile", userH.GetProfile)

			// 组织管理 (Org Management)
			// 辅导员相关 (Counselor)
			protected.POST("/departments", orgH.CreateDepartment)      // 创建部门
			protected.GET("/departments/mine", orgH.ListMyDepartments) // 我的部门
			protected.GET("/leaves/pending", leaveH.ListPendingLeaves) // 待审批请假
			protected.POST("/leaves/:uuid/audit", leaveH.Audit)        // 审批请假

			// 教师相关 (Teacher)
			protected.POST("/classes", orgH.CreateClass)       // 创建班级
			protected.GET("/classes/mine", orgH.ListMyClasses) // 我的班级

			// 通用发布 (Counselor & Teacher)
			protected.POST("/notifications", notifH.Create) // 发布通知
			//protected.POST("/tasks", taskH.CreateTask)      // 发布任务

			// 学生相关 (Student)
			protected.POST("/departments/join", orgH.StudentJoinDepartment) // 加入部门
			protected.POST("/classes/join", orgH.StudentJoinClass)          // 加入班级
			protected.POST("/leaves", leaveH.Apply)                         // 申请请假
			protected.GET("/leaves/mine", leaveH.MyLeaves)                  // 我的请假
			protected.GET("/notifications/mine", notifH.GetMyNotifications) // 我的通知
			protected.GET("/tasks/mine", taskH.GetMyTasks)                  // 我的任务
			protected.POST("/tasks/:uuid/submit", taskH.SubmitTask)         // 提交任务

			// 列表查看 (List View)
			protected.GET("/students", userH.ListStudents)

			// 打卡任务 (Ding Tasks)
			protected.POST("/createdings", dingH.Create)
			protected.GET("/mydings", dingH.ListMyDings)
			protected.GET("/mycreateddings", dingH.ListMyCreatedDings)
		}

		// 开放平台注册 (Open Platform Registration)
		openReg := api.Group("/open")
		{
			openReg.POST("/register", openH.RegisterDeveloper) // 注册开发者
			openReg.POST("/apps", openH.CreateApp)             // 创建应用
		}

		// 开放平台 API (Open Platform API)
		// 外部应用调用，使用 AppID/Secret 鉴权
		openapi := api.Group("/start/v1")
		openapi.Use(middleware.OpenAPIMiddleware(db))
		{
			// 系统向开发者开放的部分接口
			openapi.GET("/user/:id/public_profile", func(c *gin.Context) {
				// 获取公开信息示例
				c.JSON(200, gin.H{"message": "Public Profile Data", "user_id": c.Param("id")})
			})
			// 可以暴露更多只读接口...
		}
	}
}
