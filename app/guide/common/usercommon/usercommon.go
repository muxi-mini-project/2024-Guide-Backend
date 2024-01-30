package usercommon

import (
	"app/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"sync"
	"time"
)

var Db *sql.DB
var DbMutex sync.Mutex
var CodeRepo *VerificationCodeRepository

type Database struct {
	Users   map[string]*Users
	Figures map[string]string
}

type Users struct {
	Email    string
	Username string
	Password string
}

// 登录请求结构体
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 重置密码请求体结构体
type ResetPasswordRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"verification_code"`
}

// 注册请求体结构体
type RegistrationRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	VerificationCode string `json:"verification_code"`
}

type SMTPConfig struct {
	Server   string `json:"smtpServer"`
	Port     int    `json:"smtpPort"`
	Username string `json:"smtpUsername"`
	Password string `json:"smtpPassword"`
}

// 注册人物信息请求体结构体
type FigureRegistrationRequest struct {
	Username string `json:"username"`
}

type ResetCodeRepository interface {
	CreateResetCode(code *ResetCode) error
	GetResetCodeByEmail(email string) (*ResetCode, error)
	DeleteResetCodeByEmail(email string) error
}

// ResetCode 包含了重置码的信息
type ResetCode struct {
	Email string
	Code  string
}

type EmailConfig struct {
	SMTPServer   string `json:"smtpServer"`
	SMTPPort     int    `json:"smtpPort"`
	SMTPUsername string `json:"smtpUsername"`
	SMTPPassword string `json:"smtpPassword"`
}

type Config struct {
	Email EmailConfig `json:"email"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type VerificationCodeRepository struct {
	Codes map[string]string
}

// SendVerificationCode 用于发送验证码邮件
func SendVerificationCode(email, verificationCode, configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %v", err)
	}

	smtpServer := config.Email.SMTPServer
	smtpPort := config.Email.SMTPPort
	smtpUsername := config.Email.SMTPUsername
	smtpPassword := config.Email.SMTPPassword

	subject := "Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", verificationCode)
	message := "Subject: " + subject + "\r\n\r\n" + body

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	return nil
}

// GenerateAndSendVerificationCode 生成并发送验证码邮件
func GenerateAndSendVerificationCode(email, configPath string) error {
	verificationCode := GenerateVerificationCode()
	return SendVerificationCode(email, verificationCode, configPath)
}

// GenerateVerificationCode 生成默认长度的验证码
func GenerateVerificationCode() string {
	const defaultVerificationCodeLength = 6
	return generateRandomCode(defaultVerificationCodeLength)
}

func UpdateUserPassword(email, newPassword string) error {
	var user database.User
	result := database.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return fmt.Errorf("找不到用户：%v", result.Error)
	}

	// 更新密码
	user.Password = newPassword
	result = database.Db.Save(&user)
	if result.Error != nil {
		return fmt.Errorf("更新密码失败：%v", result.Error)
	}

	return nil
}

// generateRandomCode 生成指定长度的随机验证码
func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func NewResetCodeRepo() *VerificationCodeRepository {
	return &VerificationCodeRepository{
		Codes: make(map[string]string),
	}
}
