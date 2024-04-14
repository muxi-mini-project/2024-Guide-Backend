package services

import (
	"app/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type TaskService struct {
    db *gorm.DB
}

// NewTaskService 创建一个新的任务服务实例
func NewTaskService(db *gorm.DB) *TaskService {
    return &TaskService{db: db}
}

func (s *TaskService) CreateTask(userID uint, title, description string, points int) error {
    // 创建任务实例
    newTask := models.Task{
        UserID:      userID,
        Title:       title,
        Description: description,
        Points:      points,
        Completed:   false, // 默认任务未完成
    }

    // 将任务保存到数据库
    if err := s.db.Create(&newTask).Error; err != nil {
        return err
    }

    return nil
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
func (s *TaskService) GetRandomAdventureTask() (models.Task, error) {
    var task models.Task
    if err := s.db.Where("task_type = ?", "adventure").Order("RANDOM()").First(&task).Error; err != nil {
        return models.Task{}, err
    }
    return task, nil
}

func (s *TaskService) CreateCombinationTask(userID uint, title, description string, subTasks []models.SubTask) error {
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

// GetUserByID 方法用于根据用户ID从数据库中获取用户信息
func (s *TaskService) GetUserByID(userID uint) (*models.User, error) {
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

// AddExperience 方法用于给特定用户添加经验值
func (s *TaskService) AddExperience(userID uint, experience int) error {
    // 根据用户ID从数据库中获取用户信息
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return err
    }

    // 更新用户经验值
    user.Experience += experience

    // 将更新后的用户信息保存到数据库
    if err := s.db.Save(&user).Error; err != nil {
        return err
    }

    fmt.Printf("Added %d experience points to user %s\n", experience, user.Username)
    return nil
}
func (s *TaskService) CalculateCompletionPercentage(userID uint) (float64, error) {
    // 调用任务服务获取用户的所有任务
    allTasks, err := s.GetPersonalTasks(userID)
    if err != nil {
        return 0, err
    }

    // 初始化完成任务数和总任务数
    completedTasks := 0
    totalTasks := len(allTasks)

    // 计算完成任务数
    for _, task := range allTasks {
        if task.Completed {
            completedTasks++
        }
    }

    // 计算完成任务的百分比
    percentage := float64(completedTasks) / float64(totalTasks) * 100

    return percentage, nil
}
