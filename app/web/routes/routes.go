package routes

import (
	"app/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	authHandler := handlers.NewAuthHandler()

	// 用户认证相关路由
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
	}

	// 任务相关路由
	taskRoutes := router.Group("/task")
	{
		taskRoutes.POST("/create", taskHandler.CreateTask)
		taskRoutes.DELETE("/delete/:id", taskHandler.DeleteTask)
		taskRoutes.PUT("/update/:id", taskHandler.UpdateTask)
		// 其他任务相关路由...
	}

	// 团队相关路由
	teamRoutes := router.Group("/team")
	{
		teamRoutes.POST("/create", teamHandler.CreateTeam)
		teamRoutes.POST("/join", teamHandler.JoinTeam)
		teamRoutes.GET("/members/:id", teamHandler.GetTeamMembers)
		// 其他团队相关路由...
	}

	// 其他路由...
}
