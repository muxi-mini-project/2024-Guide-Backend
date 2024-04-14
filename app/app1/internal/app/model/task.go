package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID           uint             `json:"id"`
	UserID       uint             `json:"user_id"`      // 用户ID，用于关联用户
	TeamID       uint             `json:"team_id"`      // 团队ID，用于关联团队
	Title        string           `json:"title"`        // 任务标题
	Description  string           `json:"description"`  // 任务描述
	Points       int              `json:"points"`       // 任务积分
	Completed    bool             `json:"completed"`    // 是否已完成
	TaskType     string           `json:"task_type"`    // 任务类型，可以是 "personal" 或 "team"
	Contributors map[uint]float64 `json:"contributors"` // 参与者，key为用户ID，value为贡献度
}

type SubTask struct {
	gorm.Model
	Title       string `json:"title"`       // 子任务标题
	Description string `json:"description"` // 子任务描述
	Points      int    `json:"points"`      // 子任务积分
	Completed   bool   `json:"completed"`   // 是否已完成
	TaskID      uint   `json:"task_id"`     // 任务ID，用于关联父任务
}

type CombinationTask struct {
	UserID      uint      `json:"userId"` // 假设使用uint类型的ID
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SubTasks    []SubTask `json:"subTasks"`
}
