// Package service 提供数据库服务接口，包括任务管理、消费记录生成等功能
package service

import (
	"context"
	"sync"
)

// TaskProgress 任务进度信息结构体
// 用于跟踪任务执行的实时进度状态
type TaskProgress struct {
	TaskID     string  `json:"taskId"`          // 任务ID，唯一标识符
	Current    int     `json:"current"`         // 当前已完成的记录数
	Total      int     `json:"total"`           // 总记录数
	Percent    float64 `json:"percent"`         // 完成百分比 (0-100)
	IsRunning  bool    `json:"isRunning"`       // 是否正在运行
	IsCanceled bool    `json:"isCanceled"`      // 是否已取消
	Error      string  `json:"error,omitempty"` // 错误信息（如果有）
}

// taskManager 任务管理器，管理所有正在运行的任务
// 使用读写锁保证并发安全，支持多任务并发执行
type taskManager struct {
	mu          sync.RWMutex                  // 读写锁，保护并发访问
	progresses  map[string]*TaskProgress      // 任务进度映射表，key为taskID
	cancelFuncs map[string]context.CancelFunc // 任务取消函数映射表，key为taskID
}

// globalTaskManager 全局任务管理器实例
// 所有任务进度和取消操作都通过此实例进行管理
var globalTaskManager = &taskManager{
	progresses:  make(map[string]*TaskProgress),
	cancelFuncs: make(map[string]context.CancelFunc),
}

// getProgress 获取指定任务的进度信息
// 参数:
//   - taskID: 任务ID
//
// 返回值: 任务进度指针，如果不存在返回nil
func (tm *taskManager) getProgress(taskID string) *TaskProgress {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.progresses[taskID]
}

// setProgress 设置或更新任务进度信息
// 参数:
//   - taskID: 任务ID
//   - progress: 任务进度信息指针
func (tm *taskManager) setProgress(taskID string, progress *TaskProgress) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.progresses[taskID] = progress
}

// setCancelFunc 设置任务的取消函数
// 参数:
//   - taskID: 任务ID
//   - cancel: context取消函数
func (tm *taskManager) setCancelFunc(taskID string, cancel context.CancelFunc) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.cancelFuncs[taskID] = cancel
}

// cancelTask 取消指定任务的执行
// 参数:
//   - taskID: 任务ID
//
// 返回值: 如果任务存在且成功取消返回true，否则返回false
func (tm *taskManager) cancelTask(taskID string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	cancel, exists := tm.cancelFuncs[taskID]
	if exists && cancel != nil {
		cancel()
		delete(tm.cancelFuncs, taskID)
		if progress, ok := tm.progresses[taskID]; ok {
			progress.IsCanceled = true
			progress.IsRunning = false
		}
		return true
	}
	return false
}

// removeTask 移除任务的所有信息（进度和取消函数）
// 用于任务完成后清理资源，防止内存泄漏
// 参数:
//   - taskID: 任务ID
func (tm *taskManager) removeTask(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.progresses, taskID)
	delete(tm.cancelFuncs, taskID)
}

// deleteManager 删除进度管理器，管理所有正在删除的任务
// 使用读写锁保证并发安全，支持多任务并发删除
type deleteManager struct {
	mu          sync.RWMutex                  // 读写锁，保护并发访问
	progresses  map[string]*TaskProgress      // 删除进度映射表，key为taskID
	cancelFuncs map[string]context.CancelFunc // 删除取消函数映射表，key为taskID
}

// globalDeleteManager 全局删除进度管理器实例
// 所有删除进度和取消操作都通过此实例进行管理
var globalDeleteManager = &deleteManager{
	progresses:  make(map[string]*TaskProgress),
	cancelFuncs: make(map[string]context.CancelFunc),
}

// getProgress 获取指定任务的删除进度信息
// 参数:
//   - taskID: 任务ID
//
// 返回值: 删除进度指针，如果不存在返回nil
func (dm *deleteManager) getProgress(taskID string) *TaskProgress {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.progresses[taskID]
}

// setProgress 设置或更新删除进度信息
// 参数:
//   - taskID: 任务ID
//   - progress: 删除进度信息指针
func (dm *deleteManager) setProgress(taskID string, progress *TaskProgress) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.progresses[taskID] = progress
}

// setCancelFunc 设置删除任务的取消函数
// 参数:
//   - taskID: 任务ID
//   - cancel: context取消函数
func (dm *deleteManager) setCancelFunc(taskID string, cancel context.CancelFunc) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.cancelFuncs[taskID] = cancel
}

// cancelDelete 取消指定任务的删除操作
// 参数:
//   - taskID: 任务ID
//
// 返回值: 如果任务存在且成功取消返回true，否则返回false
func (dm *deleteManager) cancelDelete(taskID string) bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	cancel, exists := dm.cancelFuncs[taskID]
	if exists && cancel != nil {
		cancel()
		delete(dm.cancelFuncs, taskID)
		if progress, ok := dm.progresses[taskID]; ok {
			progress.IsCanceled = true
			progress.IsRunning = false
		}
		return true
	}
	return false
}

// removeDelete 移除删除任务的所有信息（进度和取消函数）
// 用于删除完成后清理资源，防止内存泄漏
// 参数:
//   - taskID: 任务ID
func (dm *deleteManager) removeDelete(taskID string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	delete(dm.progresses, taskID)
	delete(dm.cancelFuncs, taskID)
}
