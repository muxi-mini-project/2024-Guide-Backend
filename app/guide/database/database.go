package database

import (
	"app/common/usercommon"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
	"time"
)

var Db *gorm.DB
var DbMutex sync.Mutex

type Task struct {
	gorm.Model
	Email       string    `json:"email"`
	TaskName    string    `json:"taskName"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	Type        string    `json:"type"`
	Deadline    time.Time `json:"deadline"`
	Score       int       `json:"score"`
}

type TeamTask struct {
	gorm.Model
	TeamID   uint      `json:"teamID"`
	TaskName string    `json:"taskName"`
	Status   int       `json:"status"`
	Type     string    `json:"type"`
	Deadline time.Time `json:"deadline"`
	Score    int       `json:"score"`
}

type Team struct {
	gorm.Model
	Creator string     `json:"creator"`
	Name    string     `json:"teamName"`
	TeamID  uint       `gorm:"not null"` // 所属团队的 ID，不能为空
	Members []Member   `json:"members,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tasks   []TeamTask `json:"tasks,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Member 是团队成员模型
type Member struct {
	gorm.Model
	Email   string `gorm:"unique;not null"` // 成员邮箱，唯一且不能为空
	TeamID  uint   `gorm:"not null"`        // 所属团队的 ID，不能为空
	InvCode string `gorm:"unique;not null"`
	// 可以根据需要添加其他成员属性
}

type CombineTask struct {
	gorm.Model
	UserID   uint   `json:"userId"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

type User struct {
	gorm.Model
	Email            string `gorm:"unique;not null"`
	Password         string `gorm:"not null"`
	VerificationCode string `gorm:"not niull"`
}

type CompletedTask struct {
	gorm.Model
	TaskName    string    `json:"taskName"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	Type        string    `json:"type"`
	Deadline    time.Time `json:"deadline"`
	UserID      uint      `json:"userID"`
}

func SetTaskScore(db *gorm.DB, taskID uint, score int) error {
	result := db.Model(&Task{}).Where("id = ?", taskID).Update("score", score)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetTaskScore(db *gorm.DB, taskID uint) (int, error) {
	var score int
	result := db.Model(&Task{}).Select("score").Where("id = ?", taskID).Scan(&score)
	if result.Error != nil {
		return 0, result.Error
	}

	return score, nil
}

func InitDB(config *usercommon.Config, codeRepo *usercommon.VerificationCodeRepository, resetCodeRepo struct{ Codes map[string]string }) (*gorm.DB, error) {
	// 连接到数据库，这里使用 SQLite 作为示例
	var err error
	Db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 执行自动迁移
	if err := AutoMigrate(Db); err != nil {
		return nil, err
	}

	return Db, nil
}

func GetDB() *gorm.DB {
	return Db
}

func AutoMigrate(models ...interface{}) error {
	return Db.AutoMigrate(models...)
}
