package handlers

import (
	"app/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	// 从请求中获取表单数据
	var createTeamRequest struct {
		Name       string `json:"name" binding:"required"`
		Invitation string `json:"invitation" binding:"required"`
	}
	if err := c.ShouldBindJSON(&createTeamRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用团队服务创建团队
	if err := h.teamService.CreateTeam(createTeamRequest.Name, createTeamRequest.Invitation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建团队成功
	c.JSON(http.StatusOK, gin.H{"message": "team created successfully"})
}

func (h *TeamHandler) JoinTeam(c *gin.Context) {
	// 从请求中获取表单数据
	var joinTeamRequest struct {
		Invitation string `json:"invitation" binding:"required"`
		UserID     uint   `json:"user_id" binding:"required"`
		Position   string `json:"position" binding:"required"`
	}
	if err := c.ShouldBindJSON(&joinTeamRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用团队服务加入团队
	if err := h.teamService.JoinTeam(joinTeamRequest.Invitation, joinTeamRequest.UserID, joinTeamRequest.Position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 加入团队成功
	c.JSON(http.StatusOK, gin.H{"message": "joined team successfully"})
}

func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	// 从URL参数中获取团队ID
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	// 调用团队服务获取团队成员信息
	members, err := h.teamService.GetTeamMembers(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回团队成员信息
	c.JSON(http.StatusOK, members)
}

func (h *TeamHandler) GetUserTeams(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 调用团队服务获取用户所在的所有团队
	teams, err := h.teamService.GetUserTeams(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回用户所在的所有团队
	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) AddComment(c *gin.Context) {
	// 从请求中获取表单数据
	var addCommentRequest struct {
		MemberID        uint   `json:"member_id" binding:"required"`
		CommentedUserID uint   `json:"commented_user_id" binding:"required"`
		Comment         string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&addCommentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用团队成员服务添加评论
	if err := h.teamService.AddComment(addCommentRequest.MemberID, addCommentRequest.CommentedUserID, addCommentRequest.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 添加评论成功
	c.JSON(http.StatusOK, gin.H{"message": "comment added successfully"})
}

func (h *TeamHandler) GetComments(c *gin.Context) {
	// 从URL参数中获取团队成员ID
	memberIDStr := c.Param("id")
	memberID, err := strconv.ParseUint(memberIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	// 调用团队成员服务获取评论列表
	comments, err := h.teamService.GetComments(uint(memberID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回评论列表
	c.JSON(http.StatusOK, comments)
}
