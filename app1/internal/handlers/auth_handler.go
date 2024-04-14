package handlers

import (
	"app/internal/services"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler with the provided AuthService.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, password, and confirmation password
// @Accept json
// @Produce json
// @Param request body struct {Email string `json:"email" binding:"required,email"` Username string `json:"username" binding:"required"` Password string `json:"password" binding:"required,min=6"` ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,eqfield=Password"`} true "User registration data"
// @Success 200 {object} gin.H{"message": "User registered successfully"}
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// 实现用户注册请求处理逻辑
	var registerData struct {
		Email           string `json:"email" binding:"required,email"`
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,eqfield=Password"`
		Code            string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Register(registerData.Email, registerData.Username, registerData.Password, registerData.ConfirmPassword, registerData.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}


// Login godoc
// @Summary Log in as a user
// @Description Log in with email and password
// @Accept json
// @Produce json
// @Param request body struct {Email string `json:"email" binding:"required,email"` Password string `json:"password" binding:"required"`} true "User login data"
// @Success 200 {object} gin.H{"token": "jwt_token"}
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 401 {object} gin.H{"error": "Unauthorized"}
// @Router /login [post]
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

// ForgotPassword godoc
// @Summary Request to reset password
// @Description Request to reset password by providing email
// @Accept json
// @Produce json
// @Param request body struct {Email string `json:"email" binding:"required,email"`} true "User email for password reset"
// @Success 200 {object} gin.H{"message": "Verification code sent successfully"}
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /forgot-password [post]
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

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset user password by providing email, new password, confirmation password, and verification code
// @Accept json
// @Produce json
// @Param request body struct {Email string `json:"email" binding:"required,email"` NewPassword string `json:"newPassword" binding:"required,min=6"` ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,eqfield=NewPassword"` Code string `json:"code" binding:"required,len=6"`} true "User data for password reset"
// @Success 200 {object} gin.H{"message": "Password reset successfully"}
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /reset-password [post]
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

// SetDailyTask godoc
// @Summary Set daily task for user
// @Description Set daily task for a user by providing user ID and daily task
// @Accept json
// @Produce json
// @Param request body struct {UserID uint `json:"user_id" binding:"required"` DailyTask string `json:"daily_task" binding:"required"`} true "User ID and daily task"
// @Success 200 {object
// @Success 200 {object} gin.H{"message": "daily task set successfully"}
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /set-daily-task [post]
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

// GetDailyTask godoc
// @Summary Get daily task for user
// @Description Get daily task for a user by providing user ID
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} gin.H{"daily_task": "task_content"}
// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /get-daily-task/{id} [get]
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

// ConvertPointsToExperience godoc
// @Summary Convert user points to experience
// @Description Convert user points to experience by providing user ID and points to convert
// @Accept json
// @Produce json
// @Param request body struct {UserID uint `json:"user_id" binding:"required"` Points int `json:"points" binding:"required"`} true "User ID and points to convert"
// @Success 200 {object} gin.H{"message": "points converted to experience successfully"}
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /convert-points [post]
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

// UpgradeUser godoc
// @Summary Upgrade user level
// @Description Upgrade user level by providing user ID
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} gin.H{"message": "user upgraded successfully"}
// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /upgrade-user/{id} [put]
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

// GetUserExperienceAndLevel godoc
// @Summary Get user experience and level
// @Description Get user experience and level by providing user ID
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} gin.H{"user": "user_data"}
// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /user/{id} [get]
func (h *AuthHandler) GetUserExperienceAndLevel(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 调用用户服务获取用户经验和等级信息
	user, err := h.authService.GetUserExperienceAndLevel(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回用户经验和等级信息
	c.JSON(http.StatusOK, gin.H{"user": user})
}

const (
    verificationCodeLength = 6
)


// 生成随机验证码
func generateVerificationCode() string {
    rand.Seed(time.Now().UnixNano())
    code := ""
    for i := 0; i < verificationCodeLength; i++ {
        code += string(rand.Intn(10) + 48) // 48 是 ASCII 码中数字 0 的值
    }
    return code
}

func (h *AuthHandler) SendVerificationEmailHandler(c *gin.Context) {
    email := c.PostForm("email") // 从请求中获取电子邮件地址
    code := generateVerificationCode() // 生成验证码，您需要实现这个函数

    err := h.authService.SendVerificationEmail(email, code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}