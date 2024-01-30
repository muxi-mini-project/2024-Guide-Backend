package daily

import (
	"app/database"
	"gorm.io/gorm"
)

type DailyTask struct {
	gorm.Model
	UserID   uint   `json:"userId"`
	TaskName string `json:"taskName"`
}

func SetDailyTask(db *gorm.DB, userID uint, taskName string) error {
	task := DailyTask{
		UserID:   userID,
		TaskName: taskName,
	}

	result := db.Create(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllUsers(db *gorm.DB) ([]database.User, error) {
	var users []database.User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetDailyTask(db *gorm.DB, userID uint) (*DailyTask, error) {
	var task DailyTask
	result := db.Where("user_id = ?", userID).Find(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}
