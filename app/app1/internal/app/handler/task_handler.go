package handlers

import (
	models "app/internal/app/model"
	services "app/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var taskService = services.NewTaskService(db)

// CreateTaskHandler 创建任务处理函数
func CreateTaskHandler(taskService *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		err := taskService.CreateTask(task.UserID, task.Title, task.Description, task.Points)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

// GetTaskHandler 根据ID获取任务处理函数
func GetTaskHandler(taskService *services.TaskService) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 32)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
            return
        }

        task, err := taskService.GetTaskByID(uint(id))
        if err != nil {
            if err.Error() == "task not found" {
                c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
            }
            return
        }

        c.JSON(http.StatusOK, task)
    }
}


// GetRandomDailyTaskHandler 获取每日打卡任务处理函数
func GetRandomDailyTaskHandler(c *gin.Context) {
	// 从URL路径中解析userID
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	taskService := services.TaskService{DB: db}

	// 调用服务获取任务
	task, err := taskService.GetRandomDailyTask(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching task"})
		return
	}

	// 检查是否找到了任务
	if task == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found for the user"})
		return
	}

	// 成功找到任务，返回给客户端
	c.JSON(http.StatusOK, task)
}

// MarkTaskCompletedHandler 打卡任务完成处理函数
func MarkTaskCompletedHandler(c *gin.Context) {
	// 从请求路径中提取 taskID
	taskIDStr := c.Param("taskID")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 调用服务层方法来标记任务为完成
	err = taskService.MarkTaskAsCompleted(uint(taskID))
	if err != nil {
		// 如果更新时遇到错误，向客户端返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark task as completed"})
		return
	}

	// 如果成功，向客户端返回成功信息
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed successfully"})
}

// GetRandomAdventureTaskHandler 获取随机冒险任务处理函数
func GetRandomAdventureTaskHandler(c *gin.Context) {
	task, err := taskService.GetRandomAdventureTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateCombinationTaskHandler 创建组合任务处理函数
func CreateCombinationTaskHandler(c *gin.Context) {
	var combinationTask models.CombinationTask
	if err := c.BindJSON(&combinationTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 确保使用 taskService 实例调用 CreateCombinationTask 方法
	err := taskService.CreateCombinationTask(combinationTask.UserID, combinationTask.Title, combinationTask.Description, combinationTask.SubTasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, combinationTask)
}

// DeleteCompletedTasksHandler 删除已完成任务处理函数
func DeleteCompletedTasksHandler(c *gin.Context) {
	// 使用 taskService 实例调用 DeleteCompletedTasks 方法
	err := taskService.DeleteCompletedTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Completed tasks deleted"})
}

// CompleteTeamTaskHandler 完成团队任务处理函数
func CompleteTeamTaskHandler(c *gin.Context) {
	// 获取 taskID
	taskID := c.Param("id")
	id, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 假设从请求中获取 level 参数
	levelParam := c.Query("level")
	level, err := strconv.Atoi(levelParam) // 将 level 参数从字符串转换为 int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level"})
		return
	}

	// 假设 taskService 是之前初始化好的 *services.TaskService 实例
	err = taskService.CompleteTeamTask(uint(id), level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team task completed"})
}
