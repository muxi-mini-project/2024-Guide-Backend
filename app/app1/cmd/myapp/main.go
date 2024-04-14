package main

import (
	"fmt"
	"log"

	"app/internal/app/handler"
	"app/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	SmtpServer   string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	DBUsername   string
	DBPassword   string
	DBHost       string
	DBPort       string
	DBName       string
}

func initConfig() *Config {
	var config Config
	viper.SetConfigName("config") // 配置文件名称 (不带后缀)
	viper.SetConfigType("json")   // 明确配置文件的类型
	viper.AddConfigPath(".")      // 例如在当前目录中查找配置
	if err := viper.ReadInConfig(); err != nil { // 查找并读取配置文件
		log.Fatalf("Fatal error config file: %v \n", err)
	}
	if err := viper.Unmarshal(&config); err != nil { // 将配置信息绑定到结构体上
		log.Fatalf("Unable to decode into struct, %v \n", err)
	}
	return &config
}

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupRouter(authService *services.AuthService, taskService *services.TaskService) *gin.Engine {
	r := gin.Default()

	// 应用CORS中间件
	r.Use(CORSMiddleware())

	// 用户认证和管理路由
	r.POST("/users/:userID/upgrade", handlers.UpgradeUserHandler(authService))
	r.GET("/dailyTasks", handlers.GetDailyTaskHandler(authService))
	r.GET("/users/:userID/dailyTasks", handlers.GetDailyTaskByUserIDHandler(authService))
	r.POST("/users/:userID/convertPoints", handlers.ConvertPointsToExperienceHandler(authService))
	r.POST("/forgot_password", handlers.ForgotPasswordHandler)
	r.POST("/reset_password", handlers.ResetPasswordHandler)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/verify_code", handlers.VerifyCodeHandler)

	// 任务相关路由
	r.POST("/task", handlers.CreateTaskHandler(taskService))
	r.GET("/task/:id", handlers.GetTaskHandler(taskService))
	r.GET("/daily_task/:userID", handlers.GetRandomDailyTaskHandler)
	r.POST("/mark_completed/:taskID", handlers.MarkTaskCompletedHandler)
	r.GET("/random_adventure", handlers.GetRandomAdventureTaskHandler)
	r.POST("/create_combination", handlers.CreateCombinationTaskHandler)
	r.DELETE("/delete_completed_tasks", handlers.DeleteCompletedTasksHandler)
	r.POST("/complete_team_task/:id", handlers.CompleteTeamTaskHandler)

	// 团队相关路由
	r.POST("/create_team", handlers.CreateTeamHandler)
	r.POST("/join_team", handlers.JoinTeamHandler)

	return r
}

func main() {
	config := initConfig()

	// 数据库连接字符串配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	authService := services.NewAuthService(db)
	taskService := services.NewTaskService(db)

	r := setupRouter(authService, taskService)
	r.Run(":8080") // 启动HTTP服务器
}
