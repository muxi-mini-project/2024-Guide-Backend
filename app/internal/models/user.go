package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email              string `gorm:"uniqueIndex"`
	Username           string
	Password           string
	VerificationCode   string // 添加验证码字段
	DailyTask          string // 每日打卡任务
	DailyTaskID        uint
	Level              int
	Experience         int // 用户经验值
	SelfImprovementExp int // 自我提升经验
	WorkExp            int // 工作事务经验
	HabitExp           int // 习惯养成经验
	TodoExp            int // 待办杂事经验
}
