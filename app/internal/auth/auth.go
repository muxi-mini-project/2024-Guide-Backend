package auth

import (
	"app/common/usercommon"
	"app/database"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const verificationCodeLength = 6 // 设置验证码长度

func RegisterWithEmailVerification(email, password, confirmPassword string) error {
	// 检查密码是否匹配
	if password != confirmPassword {
		return fmt.Errorf("两次输入的密码不匹配")
	}

	// 检查密码复杂性（这里可以根据需要添加更多规则）
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少为8个字符")
	}

	// 生成哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("哈希密码生成失败: %v", err)
	}

	// 生成并发送验证码邮件
	verificationCode := generateRandomCode()
	err = usercommon.GenerateAndSendVerificationCode(email, verificationCode)
	if err != nil {
		return fmt.Errorf("验证码邮件发送失败: %v", err)
	}

	// 注册用户
	err = CreateUser(email, string(hashedPassword), verificationCode)
	if err != nil {
		return fmt.Errorf("用户注册失败: %v", err)
	}

	return nil
}

func CreateUser(email, password, verificationCode string) error {
	// 初始化一个新用户
	newUser := database.User{
		Email:            email,
		Password:         password,
		VerificationCode: verificationCode,
	}

	// 在数据库中创建用户
	result := database.Db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserByEmail 根据 email 获取用户
func GetUserByEmail(email string) (database.User, bool) {
	var user database.User
	result := database.Db.Where("email = ?", email).First(&user)

	// 如果找到用户，返回用户信息和 true；如果找不到用户，返回空用户信息和 false
	return user, result.RowsAffected > 0
}

// LoginUser 省略部分代码
func LoginUser(email, password string) error {
	// 省略部分代码

	// 获取用户信息
	user, userExists := GetUserByEmail(email)
	if !userExists {
		return fmt.Errorf("用户不存在")
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("密码不正确")
	}

	return nil
}

func ForgotPassword(email string) error {
	// 生成并发送验证码邮件
	verificationCode := generateRandomCode()
	err := usercommon.GenerateAndSendVerificationCode(email, verificationCode)
	if err != nil {
		return fmt.Errorf("验证码邮件发送失败: %v", err)
	}

	// 存储重置密码验证码
	usercommon.ResetCodeRepo.Codes[email] = verificationCode

	return nil
}

func ResetPassword(email, newPassword, confirmPassword, userInputCode string) error {
	// 验证用户输入的验证码是否正确
	if userInputCode != fmt.Sprint(usercommon.ResetCodeRepo.Codes[email]) {
		return fmt.Errorf("验证码不正确")
	}

	// 检查密码是否匹配
	if newPassword != confirmPassword {
		return fmt.Errorf("两次输入的密码不匹配")
	}

	// 检查密码复杂性（这里可以根据需要添加更多规则）
	if len(newPassword) < 8 {
		return fmt.Errorf("密码长度至少为8个字符")
	}

	// 生成哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("哈希密码生成失败: %v", err)
	}

	// 更新用户密码
	err = usercommon.UpdateUserPassword(email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("密码更新失败: %v", err)
	}

	return nil
}

func generateRandomCode() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, verificationCodeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
