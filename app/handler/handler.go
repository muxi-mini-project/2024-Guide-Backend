package handler

import (
	"app/common/usercommon"
	"app/database"
	"app/internal/auth"
	"app/internal/random"
	"app/model/daily"
	"app/model/task"
	"app/model/team"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/register", func(c *gin.Context) {
		var request struct {
			Email           string `json:"email"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirmPassword"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := auth.RegisterWithEmailVerification(request.Email, request.Password, request.ConfirmPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
	})

	router.POST("/login", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := auth.LoginUser(request.Email, request.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
	})

	router.POST("/forgot-password", func(c *gin.Context) {
		var request struct {
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := auth.ForgotPassword(request.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "验证码邮件已发送"})
	})

	router.POST("/reset-password", func(c *gin.Context) {
		var request struct {
			Email           string `json:"email"`
			NewPassword     string `json:"newPassword"`
			ConfirmPassword string `json:"confirmPassword"`
			UserInputCode   string `json:"userInputCode"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := auth.ResetPassword(request.Email, request.NewPassword, request.ConfirmPassword, request.UserInputCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
	})
	router.POST("/insert-task", func(c *gin.Context) {
		var request struct {
			Email       string    `json:"email"`
			TaskName    string    `json:"taskname"`
			Description string    `json:"description"`
			Status      int       `json:"status"`
			Type        string    `json:"thetype"`
			Deadline    time.Time `json:"deadline"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := task.InsertTask(database.Db, request.Email, request.TaskName, request.Description, request.Status, request.Type, request.Deadline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "任务添加成功"})
	})

	router.POST("/complete-task", func(c *gin.Context) {
		var request struct {
			TaskID uint `json:"taskID"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := task.CompleteTask(database.Db, request.TaskID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "任务已完成"})
	})

	router.GET("/display-tasks", func(c *gin.Context) {
		tasks := task.DisplayTasks(database.Db)
		c.JSON(http.StatusOK, tasks)
	})

	router.POST("/create-team", func(c *gin.Context) {
		var request struct {
			Creator  string `json:"creator"`
			TeamName string `json:"teamName"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		teamService := &team.TeamService{
			Team: &database.Team{
				Creator: request.Creator,
				Name:    request.TeamName,
			},
		}

		err := teamService.CreateTeam(database.Db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "团队创建成功"})
	})

	router.POST("/join-team", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
			TeamID   uint   `json:"teamID"`
			InvCode  string `json:"invCode"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		teamService := &team.TeamService{
			Team: &database.Team{}, // 需要根据 teamID 从数据库中获取实际的 Team 对象
		}

		err := teamService.JoinTeam(database.Db, request.Username, request.InvCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "成功加入团队"})
	})

	router.GET("/get-team-members", func(c *gin.Context) {
		teamID, err := strconv.Atoi(c.Query("teamID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamID"})
			return
		}

		members, err := team.GetTeamMembers(database.Db, uint(teamID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, members)
	})

	router.POST("/create-combine-task", func(c *gin.Context) {
		var request struct {
			UserID uint   `json:"userID"`
			Name   string `json:"name"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := task.CreateCombineTask(database.Db, request.UserID, request.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "组合任务创建成功"})
	})

	router.GET("/get-combine-tasks", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Query("userID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
			return
		}

		tasks, err := task.GetCombineTasks(database.Db, uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	})

	router.GET("/get-random-tasks", func(c *gin.Context) {
		email := c.Query("email")
		count, err := strconv.Atoi(c.Query("count"))
		if err != nil || count <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}

		tasks, err := random.GetRandomUserTasks(email, count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	})

	router.GET("/get-random-tasks", func(c *gin.Context) {
		email := c.Query("email")
		count, err := strconv.Atoi(c.Query("count"))
		if err != nil || count <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}

		tasks, err := random.GetRandomUserTasks(email, count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	})

	router.POST("/create-combine-task", func(c *gin.Context) {
		var request struct {
			UserID uint   `json:"userID"`
			Name   string `json:"name"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := task.CreateCombineTask(db, request.UserID, request.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "组合任务创建成功"})
	})

	router.GET("/get-combine-tasks", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Query("userID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
			return
		}

		tasks, err := task.GetCombineTasks(db, uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	})
	router.POST("/set-daily-task", func(c *gin.Context) {
		var request struct {
			UserID   uint   `json:"userID"`
			TaskName string `json:"taskName"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := daily.SetDailyTask(db, request.UserID, request.TaskName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "每日打卡任务设置成功"})
	})

	router.GET("/get-daily-task", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Query("userID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
			return
		}

		task, err := daily.GetDailyTask(db, uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, task)
	})
	router.POST("/set-task-score", func(c *gin.Context) {
		var request struct {
			TaskID uint `json:"taskID"`
			Score  int  `json:"score"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := database.SetTaskScore(db, request.TaskID, request.Score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "任务分数设置成功"})
	})
	router.GET("/get-task-score", func(c *gin.Context) {
		taskID, err := strconv.Atoi(c.Query("taskID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid taskID"})
			return
		}

		score, err := database.GetTaskScore(db, uint(taskID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"score": score})
	})

	router.POST("/send-code", func(c *gin.Context) {
		var request struct {
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		configPath := "path/to/your/config.json" // Replace with the actual path to your SMTP config file
		code := usercommon.GenerateVerificationCode()
		err := usercommon.SendVerificationCode(request.Email, code, configPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
	})

	return router
}
