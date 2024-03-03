package models

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	UserID       uint             // 用户ID，用于关联用户
	TeamID       uint             // 团队ID，用于关联团队
	Title        string           // 任务标题
	Description  string           // 任务描述
	Points       int              // 任务积分
	Completed    bool             // 是否已完成
	TaskType     string           // 任务类型，可以是 "personal" 或 "team"
	Contributors map[uint]float64 `json:"contributors"`
}

type SubTask struct {
	Title       string // 子任务标题
	Description string // 子任务描述
	Points      int    // 子任务积分
	Completed   bool   // 是否已完成
	TaskID      uint   // 任务ID，用于关联父任务
}
