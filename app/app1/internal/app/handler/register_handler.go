package handlers

import (
	models "app/internal/app/model"
	services "app/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"net/http"
)

type Config struct {
	SmtpServer   string `json:"smtp_server"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp_username"`
	SmtpPassword string `json:"smtp_password"`
}

// RegisterHandler 处理用户注册请求
func RegisterHandler(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 发送验证码到用户填写的邮箱
	verificationCode, err := services.SendVerificationCode(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	// 存储验证码到缓存中
	err = services.StoreVerificationCode(user.Email, verificationCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

// VerifyCodeHandler 处理验证码验证请求
func VerifyCodeHandler(c *gin.Context) {
	// 从请求中获取用户填写的验证码
	verificationCode := c.PostForm("verification_code")
	userEmail := c.PostForm("email")

	// 从缓存中获取之前发送的验证码
	storedCode, err := services.GetVerificationCode(userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve verification code"})
		return
	}

	// 比较用户填写的验证码和之前发送的验证码
	if verificationCode != storedCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// 验证码正确，注册用户
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := services.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func SendVerificationEmail(email, code string, config *Config) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.SmtpUsername)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/plain", "Your verification code is: "+code)

	d := gomail.NewDialer(config.SmtpServer, config.SmtpPort, config.SmtpUsername, config.SmtpPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
