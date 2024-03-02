package main

import (
	"app/internal/handlers"
	"app/internal/services"
	"app/web/routes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
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
	config, err := LoadConfig("config.json")
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %s", err))
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %s", err))
	}

	// 创建服务对象
	authService, err := services.NewAuthService(db)
	if err != nil {
		panic(fmt.Errorf("failed to create auth service: %s", err))
	}
	taskService := services.NewTaskService(db)
	teamService := services.NewTeamService(db)

	// 创建处理程序并传递服务对象
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	teamHandler := handlers.NewTeamHandler(teamService)

	// 注册路由并将处理程序传递给路由
	routes.SetupRoutes(r, authHandler, taskHandler, teamHandler)

	// 启动服务
	r.Run(":8080")
}
