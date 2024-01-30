package team

import (
	"app/database"
	"gorm.io/gorm"
)

type TeamService struct {
	Team *database.Team
}

func (T *TeamService) CreateTeam(db *gorm.DB) error {
	result := db.Create(T.Team)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (T *TeamService) JoinTeam(db *gorm.DB, username, invCode string) error {
	var member database.Member
	result := db.Where("inv_code = ?", invCode).First(&member)
	if result.Error != nil {
		return result.Error
	}

	member.Email = username
	member.TeamID = T.Team.ID // 使用 Team 的 ID 字段来关联 Member 和 Team
	result = db.Save(&member)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetTeamMembers(db *gorm.DB, teamID uint) ([]database.Member, error) {
	var members []database.Member
	result := db.Where("team_id = ?", teamID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}
