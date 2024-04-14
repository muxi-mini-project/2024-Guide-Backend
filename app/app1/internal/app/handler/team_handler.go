package handlers

import (
	models "app/internal/app/model"
	services "app/internal/app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"net/http"
)

type TeamService struct {
	db *gorm.DB
	DB *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{DB: db}
}

var db *gorm.DB

// CreateTeamHandler 创建团队处理函数
func CreateTeamHandler(c *gin.Context) {
	var team models.Team
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 创建 TeamService 的实例
	teamService := services.NewTeamService(db) // 这里的 db 是你的数据库连接，你需要替换成实际的数据库连接

	// 调用 CreateTeam 方法
	if err := teamService.CreateTeam(team.Name, team.Invitation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// JoinTeamHandler 加入团队处理函数
func JoinTeamHandler(c *gin.Context) {
	var member models.TeamMember
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 类型断言，确保 invitation 是 string 类型
	invitation, ok := member.Invitation.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation must be a string"})
		return
	}

	// 创建 TeamService 的实例
	teamService := services.NewTeamService(db)

	// 使用断言后的 invitation 字符串加入团队
	if err := teamService.JoinTeam(invitation, member.UserID, member.Position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Joined team successfully"})
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

func (s *TeamService) GetTeamMembers(teamID uint) ([]models.TeamMember, error) {
	// 查询团队成员
	var members []models.TeamMember
	if err := s.db.Where("team_id = ?", teamID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}
