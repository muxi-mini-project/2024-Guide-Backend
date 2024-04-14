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

// @Summary 创建团队
// @Description 创建一个新的团队
// @Tags Teams
// @Accept json
// @Produce json
//
//	@Param team body {object} struct {
//	  Name string `json:"name" binding:"required"` // 团队名称
//	  Invitation string `json:"invitation" binding:"required"` // 邀请码
//	}
//
// @Success 200 {object} gin.H{"message": "团队创建成功"} "创建团队成功"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /teams [post]
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

// @Summary 加入团队
// @Description 根据邀请码加入团队
// @Tags Teams
// @Accept json
// @Produce json
//
//	@Param team body {object} struct {
//	  Invitation string `json:"invitation" binding:"required"` // 邀请码
//	  UserID uint `json:"user_id" binding:"required"` // 用户ID
//	  Position string `json:"position" binding:"required"` // 职位
//	}
//
// @Success 200 {object} gin.H{"message": "成功加入团队"} "加入团队成功"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /teams/join [post]
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

// @Summary 获取团队成员列表
// @Description 根据团队ID获取团队成员列表
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path uint true "团队ID"
// @Success 200 {object} interface{} "团队成员信息"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /teams/members/{id} [get]
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

// @Summary 获取用户所在团队列表
// @Description 根据用户ID获取用户所在的所有团队
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} interface{} "用户所在团队列表"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /teams/user/{id} [get]
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

// @Summary 添加评论
// @Description 给团队成员添加评论
// @Tags Teams
// @Accept json
// @Produce json
//
//	@Param comment body {object} struct {
//	  MemberID uint `json:"member_id" binding:"required"` // 团队成员ID
//	  CommentedUserID uint `json:"commented_user_id" binding:"required"` // 被评论的用户ID
//	  Comment string `json:"comment" binding:"required"` // 评论内容
//	}
//
// @Success 200 {object} gin.H{"message": "评论添加成功"} "添加评论成功"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /teams/comments [post]
func (h *TeamHandler) AddComment(c *gin.Context) {
	// 从请求中获取表单数据
	var addCommentRequest struct {
		MemberID        uint   `json:"member_id" binding:"required"`         // 团队成员ID
		CommentedUserID uint   `json:"commented_user_id" binding:"required"` // 被评论的用户ID
		Comment         string `json:"comment" binding:"required"`           // 评论内容
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

// @Summary 获取评论列表
// @Description 根据团队成员ID获取评论列表
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path uint true "团队成员ID"
// @Success 200 {object} interface{} "评论列表"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /teams/comments/{id} [get]
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
