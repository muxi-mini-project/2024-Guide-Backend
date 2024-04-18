package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name       string // 团队名称
	Invitation string // 邀请码
	Members    []User // 团队成员，可以根据具体需求修改
}

type UserTeam struct {
	UserID uint `gorm:"primaryKey"`
	TeamID uint `gorm:"primaryKey"`
}

type TeamMember struct {
	gorm.Model
	UserID       uint   // 用户ID，用于关联用户
	TeamID       uint   // 团队ID，用于关联团队
	Position     string // 职位信息
	Contributors map[uint]float64
	Comments     map[uint]string
	Invitation   interface{}
}
