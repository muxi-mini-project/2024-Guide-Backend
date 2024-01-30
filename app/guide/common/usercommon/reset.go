package usercommon

import (
	"fmt"
	"sync"
)

type ResetCodeRepo struct {
	Codes map[string]*ResetCode
	mu    sync.Mutex
}

// CreateResetCode 创建重置码
func (r *ResetCodeRepo) CreateResetCode(code *ResetCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 假设代码是唯一的
	if _, exists := r.Codes[code.Email]; exists {
		return fmt.Errorf("重置码已存在")
	}

	r.Codes[code.Email] = code
	return nil
}

// GetResetCodeByEmail 通过邮箱获取重置码
func (r *ResetCodeRepo) GetResetCodeByEmail(email string) (*ResetCode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	code, exists := r.Codes[email]
	if !exists {
		return nil, fmt.Errorf("找不到对应的重置码")
	}

	return code, nil
}

// DeleteResetCodeByEmail 通过邮箱删除重置码
func (r *ResetCodeRepo) DeleteResetCodeByEmail(email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Codes[email]; exists {
		delete(r.Codes, email)
		return nil
	}

	return fmt.Errorf("找不到要删除的重置码")
}
