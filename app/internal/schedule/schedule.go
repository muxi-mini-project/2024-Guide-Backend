package schedule

import (
	"app/model/daily"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
)

func SetupDailyTaskReminder(c *cron.Cron, db *gorm.DB) {
	c.AddFunc("@daily", func() {
		users, err := daily.GetAllUsers(db)
		if err != nil {
			log.Printf("获取用户列表时出错: %v\n", err)
			return
		}

		for _, u := range users {
			// 获取用户的每日打卡任务
			task, err := daily.GetDailyTask(db, u.ID)
			if err != nil {
				log.Printf("获取用户 %s 的每日打卡任务时出错: %v\n", u.Email, err)
				continue
			}

			if task != nil {
				// 发送系统通知
				sendSystemNotification("每日任务提醒", fmt.Sprintf("请完成每日打卡任务: %s", task.TaskName))
			}
		}
	})

	c.Start()
}

func sendSystemNotification(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		log.Printf("发送系统通知失败: %v\n", err)
		return
	}

	log.Printf("已发送系统通知: %s\n", message)
}
