package services

import (
	models "app/internal/app/model"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) UpgradeUser(userID uint) error {
	// 获取用户的经验值和等级
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 计算阈值
	threshold := int(user.Level)*4*2 - 3*2

	// 根据经验值判断是否满足升级条件
	if user.SelfImprovementExp >= threshold &&
		user.WorkExp >= threshold &&
		user.HabitExp >= threshold &&
		user.TodoExp >= threshold {
		// 如果满足升级条件，则将用户等级加一，并扣除相应的经验值
		user.Level++
		user.SelfImprovementExp -= threshold
		user.WorkExp -= threshold
		user.HabitExp -= threshold
		user.TodoExp -= threshold
		// 保存用户信息
		if err := s.db.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *AuthService) GetDailyTask() ([]models.Task, error) {
	// 在此处编写逻辑以获取每日任务
	// 例如，从数据库中查询所有每日任务并返回

	var dailyTasks []models.Task
	if err := s.db.Where("is_daily = ?", true).Find(&dailyTasks).Error; err != nil {
		return nil, err
	}

	return dailyTasks, nil
}
func (s *AuthService) GetDailyTaskByUserID(userID uint) ([]models.Task, error) {
	// 在数据库中查询特定用户的每日任务并返回
	var dailyTasks []models.Task
	if err := s.db.Where("user_id = ? AND is_daily = ?", userID, true).Find(&dailyTasks).Error; err != nil {
		return nil, err
	}

	return dailyTasks, nil
}

func (s *AuthService) ConvertPointsToExperience(userID uint, points int) error {
	// 根据任务积分计算经验值并更新用户经验值
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}
	user.Experience += points
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetUserExperienceAndLevel(userID uint64) (*models.User, error) {
	// 根据用户ID查询用户信息
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
