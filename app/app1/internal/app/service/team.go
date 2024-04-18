package services

import (
	models "app/internal/app/model"
	"gorm.io/gorm"
)

type TeamService struct {
	db *gorm.DB
}

// NewTeamService 创建一个新的团队服务实例
func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) CreateTeam(name, invitation string) error {
	team := &models.Team{
		Name:       name,
		Invitation: invitation,
		// 可根据需要设置其他属性
	}
	if err := s.db.Create(team).Error; err != nil {
		return err
	}
	return nil
}

func (s *TeamService) JoinTeam(invitation string, userID uint, position string) error {
	// 查找对应的团队
	var team models.Team
	if err := s.db.Where("invitation = ?", invitation).First(&team).Error; err != nil {
		return err
	}

	// 创建新的团队成员并保存到数据库
	member := models.TeamMember{
		TeamID:   team.ID,
		UserID:   userID,
		Position: position,
	}
	if err := s.db.Create(&member).Error; err != nil {
		return err
	}

	return nil
}

func (s *TeamService) GetTeamMembers(teamID uint) ([]models.TeamMember, error) {
	// 查询团队成员
	var members []models.TeamMember
	if err := s.db.Where("team_id = ?", teamID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (s *TeamService) GetUserTeams(userID uint) ([]models.Team, error) {
	// 查询用户所在的所有团队
	var teams []models.Team
	if err := s.db.Table("teams").Select("teams.*").
		Joins("INNER JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userID).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *TeamService) AddComment(memberID, commentedUserID uint, comment string) error {
	// 查询被评论的团队成员
	var commentedUser models.TeamMember
	if err := s.db.First(&commentedUser, commentedUserID).Error; err != nil {
		return err
	}

	// 添加评论到被评论的团队成员的评论列表中
	commentedUser.Comments[memberID] = comment
	if err := s.db.Save(&commentedUser).Error; err != nil {
		return err
	}

	return nil
}

func (s *TeamService) GetComments(memberID uint) (map[uint]string, error) {
	// 查询团队成员
	var member models.TeamMember
	if err := s.db.First(&member, memberID).Error; err != nil {
		return nil, err
	}

	// 返回团队成员的评论列表
	return member.Comments, nil
}
