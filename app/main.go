package main

import (
	"app/internal/handlers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

type Config struct {
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
func main() {
	r := gin.Default()

	// 从配置文件加载数据库连接信息
	config, err := LoadConfig("config.json")
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %s", err))
	}

	// 连接数据库
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %s", err))
	}

	// 注册路由
	authHandler := handlers.NewAuthHandler()
	taskHandler := handlers.NewTaskHandler()
	teamHandler := handlers.NewTeamHandler()

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.POST("/set-daily-task", authHandler.SetDailyTask)
		authRoutes.GET("/get-daily-task/:id", authHandler.GetDailyTask)
		authRoutes.POST("/convert-points-to-experience", authHandler.ConvertPointsToExperience)
		authRoutes.PUT("/upgrade-user/:id", authHandler.UpgradeUser)
	}

	taskRoutes := r.Group("/tasks")
	{
		taskRoutes.POST("/create-task", taskHandler.CreateTask)
		taskRoutes.DELETE("/delete-task/:id", taskHandler.DeleteTask)
		taskRoutes.PUT("/update-task/:id", taskHandler.UpdateTask)
		taskRoutes.GET("/personal-tasks/:id", taskHandler.GetPersonalTasks)
		taskRoutes.GET("/team-tasks/:id", taskHandler.GetTeamTasks)
		taskRoutes.GET("/random-adventure-task", taskHandler.GetRandomAdventureTask)
		taskRoutes.POST("/create-combination-task", taskHandler.CreateCombinationTask)
		taskRoutes.PUT("/mark-task-as-completed/:id", taskHandler.MarkTaskAsCompleted)
		taskRoutes.DELETE("/delete-completed-tasks", taskHandler.DeleteCompletedTasks)
		taskRoutes.PUT("/complete-team-task/:id", taskHandler.CompleteTeamTask)
	}

	teamRoutes := r.Group("/teams")
	{
		teamRoutes.POST("/create-team", teamHandler.CreateTeam)
		teamRoutes.POST("/join-team", teamHandler.JoinTeam)
		teamRoutes.GET("/team-members/:id", teamHandler.GetTeamMembers)
		teamRoutes.GET("/user-teams/:id", teamHandler.GetUserTeams)
		teamRoutes.POST("/add-comment", teamHandler.AddComment)
		teamRoutes.GET("/get-comments/:id", teamHandler.GetComments)
	}

	// 启动服务
	r.Run(":8080")
}