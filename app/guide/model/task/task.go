package task

import (
	"app/database"
	"gorm.io/gorm"
	"time"
)

type CompletedTask struct {
	gorm.Model
	Task        database.Task `json:"task" gorm:"foreignKey:TaskID"`
	CompletedBy uint          `json:"completedBy"`
	TaskID      uint          `json:"taskId"`
}

func InsertTask(db *gorm.DB, email, taskname, description string, status int, thetype string, deadline time.Time) error {
	newTask := database.Task{
		Email:       email,
		TaskName:    taskname,
		Description: description,
		Status:      status,
		Type:        thetype,
		Deadline:    deadline,
	}

	result := db.Create(&newTask)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DisplayTasks(db *gorm.DB) []database.Task {
	var tasks []database.Task
	db.Find(&tasks)
	return tasks
}

func CompleteTask(db *gorm.DB, taskID uint) error {
	var task database.Task
	result := db.First(&task, taskID)
	if result.Error != nil {
		return result.Error
	}

	// 在这里可以添加一些逻辑，例如任务完成时的处理

	// 将任务移到已完成任务表
	completedTask := CompletedTask{
		Task:        task,
		CompletedBy: 1, // 你可能需要根据实际情况设置 CompletedBy
	}
	db.Create(&completedTask)

	// 删除原任务
	db.Delete(&task)

	return nil
}

func CreateCombineTask(db *gorm.DB, userID uint, name string) error {
	task := database.CombineTask{
		UserID:   userID,
		Name:     name,
		Complete: false,
	}

	result := db.Create(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetCombineTasks(db *gorm.DB, userID uint) ([]database.CombineTask, error) {
	var tasks []database.CombineTask
	result := db.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
