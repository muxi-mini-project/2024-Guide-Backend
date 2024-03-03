package routes

import (
	"app/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authService *handlers.AuthHandler, taskService *handlers.TaskHandler, teamService *handlers.TeamHandler) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authService.Register)
		authRoutes.POST("/login", authService.Login)
		authRoutes.POST("/forgot-password", authService.ForgotPassword)
		authRoutes.POST("/reset-password", authService.ResetPassword)
	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.POST("/create-task", taskService.CreateTask)
		taskRoutes.DELETE("/delete-task/:id", taskService.DeleteTask)
		taskRoutes.PUT("/update-task/:id", taskService.UpdateTask)
		taskRoutes.PUT("/mark-task-as-completed/:id", taskService.MarkTaskAsCompleted)
	}

	teamRoutes := router.Group("/teams")
	{
		teamRoutes.POST("/create-team", teamService.CreateTeam)
		teamRoutes.POST("/join-team", teamService.JoinTeam)
		teamRoutes.GET("/team-members/:id", teamService.GetTeamMembers)
		teamRoutes.POST("/add-comment", teamService.AddComment)
	}
}
