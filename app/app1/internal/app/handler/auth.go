package handlers

import (
	services "app/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 升级用户的处理器
func UpgradeUserHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		err = authService.UpgradeUser(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "升级用户失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "用户成功升级"})
	}
}

// 获取所有每日任务的处理器
func GetDailyTaskHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := authService.GetDailyTask()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取每日任务失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"dailyTasks": tasks})
	}
}

// 获取指定用户的每日任务的处理器
func GetDailyTaskByUserIDHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		tasks, err := authService.GetDailyTaskByUserID(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户每日任务失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"userDailyTasks": tasks})
	}
}

// 将积分转换为用户经验的处理器
func ConvertPointsToExperienceHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		var req struct {
			Points int `json:"points"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
			return
		}

		err = authService.ConvertPointsToExperience(uint(userID), req.Points)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "积分转换经验失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "积分成功转换为经验"})
	}
}
