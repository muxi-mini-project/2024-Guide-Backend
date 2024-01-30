package random

import (
	"fmt"
	"math/rand"
	"time"
)

var UserTaskMap = make(map[string][]Task)

type Task struct {
	Name     string `json:"name"`
	Points   int    `json:"points"`
	Complete bool   `json:"complete"`
}

func AddUserTask(email string, task Task) {
	UserTaskMap[email] = append(UserTaskMap[email], task)
}

func GetRandomUserTasks(email string, count int) ([]Task, error) {
	userTasks, exists := UserTaskMap[email]
	if !exists {
		return nil, fmt.Errorf("用户 %s 不存在", email)
	}

	// 随机获取 count 个任务
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userTasks), func(i, j int) {
		userTasks[i], userTasks[j] = userTasks[j], userTasks[i]
	})

	if count > len(userTasks) {
		count = len(userTasks)
	}

	return userTasks[:count], nil
}
