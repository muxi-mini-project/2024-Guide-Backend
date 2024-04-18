package services

import (
	models "app/internal/app/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type TaskService struct {
	db *gorm.DB
	DB *gorm.DB
}

// NewTaskService 创建一个新的任务服务实例

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	result := s.DB.First(&task, id) // GORM 使用 First 方法查找第一个匹配的记录
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("task not found")
		}
		return nil, result.Error
	}

	return &task, nil
}

func (s *TaskService) CreateTask(userID uint, title string, description string, points int) error {
	// 参数验证
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if points < 0 {
		return errors.New("points must be non-negative")
	}

	// 创建 Task 实例
	newTask := models.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Points:      points,
		Completed:   false, // 默认任务未完成
	}

	// 将新任务保存到数据库
	if err := s.DB.Create(&newTask).Error; err != nil {
		return err // 如果数据库操作出错，则返回错误
	}

	return nil // 如果没有错误发生，返回 nil
}

func (s *TaskService) DeleteTask(taskID uint) error {
	// 查询要删除的任务
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 从数据库中删除任务
	if err := s.db.Delete(&task).Error; err != nil {
		return err
	}

	return nil
}

func (s *TaskService) UpdateTask(taskID uint, title, description string, points int, completed bool) error {
	// 查询要更新的任务
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 更新任务属性
	task.Title = title
	task.Description = description
	task.Points = points
	task.Completed = completed

	// 保存更新后的任务到数据库
	if err := s.db.Save(&task).Error; err != nil {
		return err
	}

	return nil
}

func (s *TaskService) GetPersonalTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetTeamTasks(teamID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Where("team_id = ?", teamID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetRandomAdventureTask() (*models.Task, error) {
	var task models.Task
	// 示例：从数据库中随机获取一个冒险类型的任务
	// 这里的查询条件和逻辑需要根据你的具体情况调整
	err := s.DB.Where("type = ?", "adventure").Order("RAND()").First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) CreateCombinationTask(userID uint, title string, description string, subTasks []models.SubTask) error {
	// 创建组合任务实例
	newCombinationTask := models.Task{
		// 设置组合任务的属性
		TaskType:    "combination",
		UserID:      userID,
		Title:       title,
		Description: description,
		// 其他属性根据实际情况设置
	}

	// 将组合任务保存到数据库
	if err := s.db.Create(&newCombinationTask).Error; err != nil {
		return err
	}

	// 保存组合任务的ID以及所有子任务的关联关系到数据库
	for _, subTask := range subTasks {
		// 创建子任务关联关系实例
		taskSubTask := models.SubTask{
			TaskID:      newCombinationTask.ID,
			Title:       subTask.Title,
			Description: subTask.Description,
			Points:      subTask.Points,
			Completed:   subTask.Completed,
		}

		// 将子任务关联关系保存到数据库
		if err := s.db.Create(&taskSubTask).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) MarkTaskAsCompleted(taskID uint) error {
	// 查询要标记为已完成的任务
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 将任务标记为已完成
	task.Completed = true

	// 保存更新后的任务到数据库
	if err := s.db.Save(&task).Error; err != nil {
		return err
	}

	return nil
}

func (s *TaskService) DeleteCompletedTasks() error {
	// 查询所有已完成的任务
	var tasks []models.Task
	if err := s.db.Where("completed = ?", true).Find(&tasks).Error; err != nil {
		return err
	}

	// 逐个删除已完成的任务
	for _, task := range tasks {
		if err := s.db.Delete(&task).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) CompleteTeamTask(taskID uint, level int) error {
	// 查询团队任务
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 分配经验值给贡献者
	for memberID, percentage := range task.Contributors {
		// 计算贡献的经验值
		experience := int(float64(task.Points) * percentage)
		// 此处省略了对用户添加经验值的逻辑，需要根据实际情况实现
		fmt.Printf("Adding %d experience points to user with ID %d\n", experience, memberID)
	}

	// 标记团队任务为已完成
	task.Completed = true
	if err := s.db.Save(&task).Error; err != nil {
		return err
	}

	return nil
}

func (s *TaskService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *TaskService) AddExperience(userID uint, experience int) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	user.Experience += experience

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	fmt.Printf("Added %d experience points to user %s\n", experience, user.Username)
	return nil
}

func (s *TaskService) CalculateCompletionPercentage(userID uint) (float64, error) {
	allTasks, err := s.GetPersonalTasks(userID)
	if err != nil {
		return 0, err
	}

	completedTasks := 0
	totalTasks := len(allTasks)

	for _, task := range allTasks {
		if task.Completed {
			completedTasks++
		}
	}

	percentage := float64(completedTasks) / float64(totalTasks) * 100

	return percentage, nil
}

// GetRandomDailyTask 为指定用户随机返回一个任务
func (ts *TaskService) GetRandomDailyTask(userID uint) (*models.Task, error) {
	var tasks []models.Task
	var task models.Task

	// 从数据库中获取该用户的所有任务
	result := ts.DB.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(tasks) == 0 {
		// 如果用户没有任务，返回nil
		return nil, nil
	}

	// 随机选择一个任务
	rand.Seed(time.Now().UnixNano())
	task = tasks[rand.Intn(len(tasks))]

	return &task, nil
}
