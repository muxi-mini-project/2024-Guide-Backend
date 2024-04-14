package main

import (
	"app/internal/handlers"
	"app/internal/services"
	"app/web/routes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

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

	// 加载配置文件
	/*config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config file: %s", err)
	}*/

	// 连接数据库
	dsn := fmt.Sprintf("root:kevin123456@tcp(localhost:3306)/theleader?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)

	// 创建服务对象
	authService, err := services.NewAuthService(db, logger)
	if err != nil {
		log.Fatalf("failed to create auth service: %s", err)
	}
	taskService := services.NewTaskService(db)
	teamService := services.NewTeamService(db)

	// 创建处理程序并传递服务对象
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	teamHandler := handlers.NewTeamHandler(teamService)

	// 注册路由并将处理程序传递给路由
	setupRoutes(r, authHandler, taskHandler, teamHandler)

	// 启动服务
	r.Run(":8080")
}

func setupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, taskHandler *handlers.TaskHandler, teamHandler *handlers.TeamHandler) {
	routes.SetupRoutes(r, authHandler, taskHandler, teamHandler)
}
