package handlers

import (
	services "app/internal/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForgotPasswordHandler(c *gin.Context) {
	email := c.PostForm("email")

	// 发送验证码邮件给用户
	verificationCode, err := services.SendVerificationCode(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	// 返回成功响应，提示用户检查邮箱获取验证码
	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent to your email. Please check your inbox."})

	// 保存验证码到上下文中，以便后续验证密码重置请求
	c.Set("verificationCode", verificationCode)
}

func ResetPasswordHandler(c *gin.Context) {
	email := c.PostForm("email")
	verificationCode := c.PostForm("verification_code")
	newPassword := c.PostForm("new_password")

	// 从上下文中获取之前发送的验证码
	savedVerificationCode, exists := c.Get("verificationCode")
	if !exists || savedVerificationCode.(string) != verificationCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// 重置密码
	err := services.ResetPassword(email, verificationCode, newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
