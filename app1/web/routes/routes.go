package routes

import (
	"app/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authService *handlers.AuthHandler, taskService *handlers.TaskHandler, teamService *handlers.TeamHandler) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authService.Register)
		authRoutes.POST("/send-verification-email", authService.SendVerificationEmailHandler)
		authRoutes.POST("/login", authService.Login)
		authRoutes.POST("/forgot-password", authService.ForgotPassword)
		authRoutes.POST("/reset-password", authService.ResetPassword)
		authRoutes.POST("/set-daily-task", authService.SetDailyTask)
		authRoutes.GET("/daily-task", authService.GetDailyTask) // Updated route
		authRoutes.POST("/convert-points-to-experience", authService.ConvertPointsToExperience)
		authRoutes.PUT("/upgrade-user/:userID", authService.UpgradeUser)                            // Updated route
		authRoutes.GET("/user-experience-and-level/:userID", authService.GetUserExperienceAndLevel) // Updated route
	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.POST("/create-task", taskService.CreateTask)
		taskRoutes.DELETE("/delete-task/:taskID", taskService.DeleteTask)                           // Updated route
		taskRoutes.PUT("/update-task/:taskID", taskService.UpdateTask)                              // Updated route
		taskRoutes.PUT("/mark-task-as-completed/:taskID", taskService.MarkTaskAsCompleted)          // Updated route
		taskRoutes.GET("/completion-percentage/:userID", taskService.CalculateCompletionPercentage) // Updated route
	}

	teamRoutes := router.Group("/teams")
	{
		teamRoutes.POST("/create-team", teamService.CreateTeam)
		teamRoutes.POST("/join-team", teamService.JoinTeam)
		teamRoutes.GET("/team-members/:teamID", teamService.GetTeamMembers) // Updated route
		teamRoutes.GET("/user-teams/:userID", teamService.GetUserTeams)     // Updated route
		teamRoutes.POST("/add-comment", teamService.AddComment)
		teamRoutes.GET("/comments/:memberID", teamService.GetComments) // Updated route
	}
}
