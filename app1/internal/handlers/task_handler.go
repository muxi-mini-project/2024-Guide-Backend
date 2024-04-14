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

// @Summary 创建任务
// @Description 创建新任务
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body createTaskRequest true "任务信息"
// @Success 200 {object} gin.H{"message": "任务创建成功"}
// @Failure 400 {object} gin.H{"error": "错误消息"}
// @Failure 500 {object} gin.H{"error": "服务器内部错误"}
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
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

	if err := h.taskService.CreateTask(createTaskRequest.UserID, createTaskRequest.Title, createTaskRequest.Description, createTaskRequest.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务创建成功"})
}

// @Summary 删除任务
// @Description 根据任务ID删除任务
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "任务ID"
// @Success 200 {object} gin.H{"message": "任务删除成功"}
// @Failure 400 {object} gin.H{"error": "错误消息"}
// @Failure 500 {object} gin.H{"error": "服务器内部错误"}
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.taskService.DeleteTask(uint(taskID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}

// @Summary 更新任务
// @Description 根据任务ID更新任务信息
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "任务ID"
// @Param request body updateTaskRequest true "更新后的任务信息"
// @Success 200 {object} gin.H{"message": "任务更新成功"}
// @Failure 400 {object} gin.H{"error": "错误消息"}
// @Failure 500 {object} gin.H{"error": "服务器内部错误"}
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

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

	if err := h.taskService.UpdateTask(uint(taskID), updateTaskRequest.Title, updateTaskRequest.Description, updateTaskRequest.Points, updateTaskRequest.Completed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务更新成功"})
}

// @Summary 获取个人任务列表
// @Description 根据用户ID获取个人任务列表
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} []models.Task "个人任务列表"
// @Failure 400 {object} gin.H{"error": "错误消息"}
// @Failure 500 {object} gin.H{"error": "服务器内部错误"}
// @Router /tasks/personal/{id} [get]
func (h *TaskHandler) GetPersonalTasks(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tasks, err := h.taskService.GetPersonalTasks(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary 获取团队任务列表
// @Description 根据团队ID获取团队任务列表
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "团队ID"
// @Success 200 {object} []models.Task "团队任务列表"
// @Failure 400 {object} gin.H{"error": "错误消息"}
// @Failure 500 {object} gin.H{"error": "服务器内部错误"}
// @Router /tasks/team/{id} [get]
func (h *TaskHandler) GetTeamTasks(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	tasks, err := h.taskService.GetTeamTasks(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary 获取随机冒险任务
// @Description 获取随机的冒险任务
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} models.Task "随机冒险任务"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"}
// @Router /tasks/random_adventure [get]
func (h *TaskHandler) GetRandomAdventureTask(c *gin.Context) {
	task, err := h.taskService.GetRandomAdventureTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary 创建组合任务
// @Description 创建一个组合任务，该任务由多个子任务组成
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body createCombinationTaskRequest true "组合任务信息"
// @Success 200 {object} gin.H{"message": "组合任务创建成功"} "成功消息"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /tasks/combination [post]
func (h *TaskHandler) CreateCombinationTask(c *gin.Context) {
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

	if err := h.taskService.CreateCombinationTask(createCombinationTaskRequest.UserID, createCombinationTaskRequest.Title, createCombinationTaskRequest.Description, createCombinationTaskRequest.SubTasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "组合任务创建成功"})
}

// @Summary 将任务标记为完成
// @Description 将任务标记为已完成状态
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "任务ID"
// @Success 200 {object} gin.H{"message": "任务标记为完成成功"} "成功消息"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /tasks/{id}/complete [put]
func (h *TaskHandler) MarkTaskAsCompleted(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.taskService.MarkTaskAsCompleted(uint(taskID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务标记为完成成功"})
}

// 继续添加其他API文档化注释...
// @Summary 删除已完成任务
// @Description 删除系统中已完成的任务记录
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"message": "已完成任务删除成功"} "成功消息"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /tasks/completed [delete]
func (h *TaskHandler) DeleteCompletedTasks(c *gin.Context) {
	if err := h.taskService.DeleteCompletedTasks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已完成任务删除成功"})
}

// @Summary 完成团队任务
// @Description 标记指定团队任务为已完成状态，并指定任务级别
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "团队任务ID"
// @Param level query int true "任务级别"
// @Success 200 {object} gin.H{"message": "团队任务完成成功"} "成功消息"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /tasks/team/{id}/complete [put]
func (h *TaskHandler) CompleteTeamTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	levelStr := c.Query("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level"})
		return
	}

	if err := h.taskService.CompleteTeamTask(uint(taskID), level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "团队任务完成成功"})
}

// @Summary 计算用户完成任务的百分比
// @Description 根据用户ID计算完成任务的百分比
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} gin.H{"percentage": float64} "完成任务百分比"
// @Failure 400 {object} gin.H{"error": "错误消息"} "请求错误"
// @Failure 500 {object} gin.H{"error": "服务器内部错误"} "服务器内部错误"
// @Router /tasks/completion/{id} [get]
func (h *TaskHandler) CalculateCompletionPercentage(c *gin.Context) {
	userIDStr := c.Param("id")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	percentage, err := h.taskService.CalculateCompletionPercentage(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"percentage": percentage})
}
