package handlers

import (
	"app/internal/models"
	"app/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}
func (h *TaskHandler) CreateTask(c *gin.Context) {
	// 从请求中获取表单数据
	var createTaskRequest struct {
		UserID      uint   `json:"user_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Points      int    `json:"points" binding:"required"`
	}
	if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用任务服务创建任务
	if err := h.taskService.CreateTask(createTaskRequest.UserID, createTaskRequest.Title, createTaskRequest.Description, createTaskRequest.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建任务成功
	c.JSON(http.StatusOK, gin.H{"message": "task created successfully"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// 从URL参数中获取任务ID
	taskIDStr := c.Param("id")

	// 将任务ID字符串转换为uint类型
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 调用任务服务删除任务
	if err := h.taskService.DeleteTask(uint(taskID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除任务成功
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// 从URL参数中获取任务ID
	taskIDStr := c.Param("id")

	// 将任务ID字符串转换为uint类型
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 从请求中获取表单数据
	var updateTaskRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Points      int    `json:"points"`
		Completed   bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&updateTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用任务服务更新任务
	if err := h.taskService.UpdateTask(uint(taskID), updateTaskRequest.Title, updateTaskRequest.Description, updateTaskRequest.Points, updateTaskRequest.Completed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新任务成功
	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

func (h *TaskHandler) GetPersonalTasks(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIDStr := c.Param("id")

	// 将用户ID字符串转换为uint类型
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 调用任务服务获取个人任务列表
	tasks, err := h.taskService.GetPersonalTasks(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回个人任务列表
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTeamTasks(c *gin.Context) {
	// 从URL参数中获取团队ID
	teamIDStr := c.Param("id")

	// 将团队ID字符串转换为uint类型
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// 调用任务服务获取团队任务列表
	tasks, err := h.taskService.GetTeamTasks(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回团队任务列表
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetRandomAdventureTask(c *gin.Context) {
	// 调用任务服务获取随机冒险任务
	task, err := h.taskService.GetRandomAdventureTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回随机冒险任务
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) CreateCombinationTask(c *gin.Context) {
	// 从请求中获取表单数据
	var createCombinationTaskRequest struct {
		UserID      uint             `json:"user_id" binding:"required"`
		Title       string           `json:"title" binding:"required"`
		Description string           `json:"description"`
		SubTasks    []models.SubTask `json:"sub_tasks" binding:"required"`
	}
	if err := c.ShouldBindJSON(&createCombinationTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用任务服务创建组合任务
	if err := h.taskService.CreateCombinationTask(createCombinationTaskRequest.UserID, createCombinationTaskRequest.Title, createCombinationTaskRequest.Description, createCombinationTaskRequest.SubTasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建组合任务成功
	c.JSON(http.StatusOK, gin.H{"message": "combination task created successfully"})
}

func (h *TaskHandler) MarkTaskAsCompleted(c *gin.Context) {
	// 从URL参数中获取任务ID
	taskIDStr := c.Param("id")

	// 将任务ID字符串转换为uint类型
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	// 调用任务服务将任务标记为完成
	if err := h.taskService.MarkTaskAsCompleted(uint(taskID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 任务标记为完成成功
	c.JSON(http.StatusOK, gin.H{"message": "task marked as completed successfully"})
}
func (h *TaskHandler) DeleteCompletedTasks(c *gin.Context) {
	// 调用任务服务定期删除已完成任务
	if err := h.taskService.DeleteCompletedTasks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 已完成任务删除成功
	c.JSON(http.StatusOK, gin.H{"message": "completed tasks deleted successfully"})
}

// handlers/team_task_handler.go

func (h *TaskHandler) CompleteTeamTask(c *gin.Context) {
	// 从URL参数中获取团队任务ID
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	// 从请求中获取团队任务级别
	levelStr := c.Query("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level"})
		return
	}

	// 调用团队任务服务完成团队任务
	if err := h.taskService.CompleteTeamTask(uint(taskID), level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 团队任务完成成功
	c.JSON(http.StatusOK, gin.H{"message": "team task completed successfully"})
}
