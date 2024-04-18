package handlers

import (
	models "app/internal/app/model"
	services "app/internal/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 发送验证码
	verificationCode, err := services.SendVerificationCode(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	// 设置验证码到上下文中
	c.Set("verificationCode", verificationCode)

	c.Status(http.StatusOK)
}
