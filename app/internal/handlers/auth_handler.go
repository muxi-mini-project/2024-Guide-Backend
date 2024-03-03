package handlers

import (
	"app/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// 实现用户注册请求处理逻辑
	var registerData struct {
		Email           string `json:"email" binding:"required,email"`
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,eqfield=Password"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Register(registerData.Email, registerData.Username, registerData.Password, registerData.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	// 实现用户登录请求处理逻辑
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	// 实现忘记密码请求处理逻辑
	var forgotPasswordData struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&forgotPasswordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ForgotPassword(forgotPasswordData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	// 实现重置密码请求处理逻辑
	var resetPasswordData struct {
		Email           string `json:"email" binding:"required,email"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,eqfield=NewPassword"`
		Code            string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&resetPasswordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ResetPassword(resetPasswordData.Email, resetPasswordData.NewPassword, resetPasswordData.ConfirmPassword, resetPasswordData.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *AuthHandler) SetDailyTask(c *gin.Context) {
	// 从请求中获取表单数据
	var setDailyTaskRequest struct {
		UserID    uint   `json:"user_id" binding:"required"`
		DailyTask string `json:"daily_task" binding:"required"`
	}
	if err := c.ShouldBindJSON(&setDailyTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用用户服务设置每日打卡任务
	if err := h.authService.SetDailyTask(setDailyTaskRequest.UserID, setDailyTaskRequest.DailyTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置每日打卡任务成功
	c.JSON(http.StatusOK, gin.H{"message": "daily task set successfully"})
}

func (h *AuthHandler) GetDailyTask(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 调用用户服务获取每日打卡任务
	dailyTask, err := h.authService.GetDailyTaskByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回每日打卡任务
	c.JSON(http.StatusOK, gin.H{"daily_task": dailyTask})
}

func (h *AuthHandler) ConvertPointsToExperience(c *gin.Context) {
	// 从请求中获取表单数据
	var convertPointsRequest struct {
		UserID uint `json:"user_id" binding:"required"`
		Points int  `json:"points" binding:"required"`
	}
	if err := c.ShouldBindJSON(&convertPointsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用用户服务将任务积分转换为经验值
	if err := h.authService.ConvertPointsToExperience(convertPointsRequest.UserID, convertPointsRequest.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换成功
	c.JSON(http.StatusOK, gin.H{"message": "points converted to experience successfully"})
}

func (h *AuthHandler) UpgradeUser(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 调用用户服务升级用户
	if err := h.authService.UpgradeUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 升级成功
	c.JSON(http.StatusOK, gin.H{"message": "user upgraded successfully"})
}
