package services

import (
	"app/internal/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"net/smtp"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

const (
	verificationCodeLength = 6
)

type AuthService struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// 生成随机验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < verificationCodeLength; i++ {
		code += string(rand.Intn(10) + 48) // 48 是 ASCII 码中数字 0 的值
	}
	return code
}

func (s *AuthService) Register(email, username, password, confirmPassword string) error {
	// 检查密码和确认密码是否匹配
	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	// 检查邮箱是否已经被注册
	existingUser := &models.User{}
	if err := s.db.Where("email = ?", email).First(existingUser).Error; err == nil {
		return errors.New("email already registered")
	}

	// 创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := &models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}
	if err := s.db.Create(newUser).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	// 根据邮箱查询用户
	user := &models.User{}
	if err := s.db.Where("email = ?", email).First(user).Error; err != nil {
		return "", errors.New("user not found")
	}

	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// 生成JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	// 实现忘记密码逻辑

	// 检查邮箱是否存在
	user := &models.User{}
	if err := s.db.Where("email = ?", email).First(user).Error; err != nil {
		return errors.New("user not found")
	}

	// 生成随机验证码
	verificationCode := generateVerificationCode()

	// 存储验证码到数据库
	user.VerificationCode = verificationCode
	if err := s.db.Save(user).Error; err != nil {
		return err
	}

	// 发送包含验证码的邮件给用户
	err := sendVerificationEmail(email, verificationCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(email, newPassword, confirmPassword, code string) error {
	// 实现重置密码逻辑

	// 验证验证码是否正确
	user := &models.User{}
	if err := s.db.Where("email = ?", email).First(user).Error; err != nil {
		return errors.New("user not found")
	}

	// 获取存储的验证码，这里需要根据实际情况从数据库或其他存储介质中获取验证码
	storedCode := user.VerificationCode

	// 检查用户输入的验证码是否正确
	if code != storedCode {
		return errors.New("invalid verification code")
	}

	// 更新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.db.Model(user).Update("password", string(hashedPassword)).Error; err != nil {
		return err
	}

	return nil
}

func generateToken(userID uint) (string, error) {
	// 定义 JWT Token 的过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建一个新的 JWT Token
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "your_app_name", // 发行者名称
		Subject:   string(userID),  // 将用户ID作为主题
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名 JWT Token
	secretKey := []byte("your_secret_key") // 用于签名的密钥，务必保密
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func sendVerificationEmail(email, code string) error {
	// 从环境变量中读取 SMTP 服务器地址和端口
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	// 从环境变量中读取发件人邮箱和密码
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	// 设置收件人邮箱
	to := []string{email}

	// 构建电子邮件内容
	subject := "Verification Code"
	body := "Your verification code is: " + code

	// 组装 MIME 消息
	message := []byte("Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// 连接到 SMTP 服务器
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), auth, senderEmail, to, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SetDailyTask(userID uint, taskIDStr string) error {
	// 将 taskIDStr 转换为 uint 类型
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}

	// 查询用户
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 查询任务
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 更新用户的每日任务ID
	user.DailyTaskID = uint(taskID)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetDailyTask() ([]models.Task, error) {
	// 在此处编写逻辑以获取每日任务
	// 例如，从数据库中查询所有每日任务并返回

	var dailyTasks []models.Task
	if err := s.db.Where("is_daily = ?", true).Find(&dailyTasks).Error; err != nil {
		return nil, err
	}

	return dailyTasks, nil
}

func (s *AuthService) GetDailyTaskByUserID(userID uint) ([]models.Task, error) {
	// 在数据库中查询特定用户的每日任务并返回
	var dailyTasks []models.Task
	if err := s.db.Where("user_id = ? AND is_daily = ?", userID, true).Find(&dailyTasks).Error; err != nil {
		return nil, err
	}

	return dailyTasks, nil
}

func (s *AuthService) ConvertPointsToExperience(userID uint, points int) error {
	// 根据任务积分计算经验值并更新用户经验值
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}
	user.Experience += points
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *AuthService) UpgradeUser(userID uint) error {
	// 获取用户的经验值和等级
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 计算阈值
	threshold := int(user.Level)*4*2 - 3*2

	// 根据经验值判断是否满足升级条件
	if user.SelfImprovementExp >= threshold &&
		user.WorkExp >= threshold &&
		user.HabitExp >= threshold &&
		user.TodoExp >= threshold {
		// 如果满足升级条件，则将用户等级加一，并扣除相应的经验值
		user.Level++
		user.SelfImprovementExp -= threshold
		user.WorkExp -= threshold
		user.HabitExp -= threshold
		user.TodoExp -= threshold
		// 保存用户信息
		if err := s.db.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
