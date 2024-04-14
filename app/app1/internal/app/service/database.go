package services

import (
	models "app/internal/app/model"
	"crypto/rand"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
)

var db *sql.DB

var (
	// verificationCodes 用于存储验证码的全局变量
	verificationCodes = make(map[string]string)
	// mu 用于确保对 verificationCodes 的并发安全访问
	mu sync.Mutex
)



func RegisterUser(user models.User) error {
	// 检查用户名或邮箱是否已存在
	var count int
	row := db.QueryRow("SELECT COUNT(id) FROM users WHERE username = ? OR email = ?", user.Username, user.Email)
	row.Scan(&count)
	if count > 0 {
		return errors.New("Username or email already exists")
	}

	// 生成盐和加密密码
	salt := generateRandomString(16)
	hashedPassword := hashPassword(user.Password, salt)

	// 插入用户到数据库
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(user models.User) error {
	// 查找用户
	var storedPassword, salt string
	row := db.QueryRow("SELECT password FROM users WHERE username = ? OR email = ?", user.Username, user.Username)
	err := row.Scan(&storedPassword)
	if err != nil {
		return errors.New("Invalid username or password")
	}

	// 获取盐并验证密码
	salt = strings.Split(storedPassword, "$")[1]
	if storedPassword != hashPassword(user.Password, salt) {
		return errors.New("Invalid username or password")
	}

	return nil
}

func ResetPassword(email string, code string, password string) error {
	// 生成随机密码重置令牌
	resetToken := generateRandomString(32)

	// 将重置令牌存储在数据库中，通常需要设置过期时间等
	_, err := db.Exec("UPDATE users SET reset_token = ? WHERE email = ?", resetToken, email)
	if err != nil {
		return err
	}

	// 发送包含重置令牌的电子邮件到用户邮箱

	return nil
}

func SendVerificationCode(email string) (string, error) {
	// 生成随机的6位数字验证码
	verificationCode := generateRandomNumericCode(6)

	// 将验证码存储在数据库中，通常需要设置过期时间等
	_, err := db.Exec("UPDATE users SET verification_code = ? WHERE email = ?", verificationCode, email)
	if err != nil {
		return "", err
	}

	// 发送包含验证码的电子邮件到用户邮箱或通过短信发送验证码给用户

	return verificationCode, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func generateRandomNumericCode(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func StoreVerificationCode(email, code string) error {
	// 加锁以确保并发安全访问
	mu.Lock()
	defer mu.Unlock()

	// 将验证码存储到缓存中
	verificationCodes[email] = code

	return nil
}

func GetVerificationCode(email string) (string, error) {
	// 加锁以确保并发安全访问
	mu.Lock()
	defer mu.Unlock()

	// 从缓存中获取验证码
	code, ok := verificationCodes[email]
	if !ok {
		return "", errors.New("verification code not found")
	}

	return code, nil
}

func DeleteVerificationCode(email string) error {
	// 加锁以确保并发安全访问
	mu.Lock()
	defer mu.Unlock()

	// 从缓存中删除验证码
	delete(verificationCodes, email)

	return nil
}
