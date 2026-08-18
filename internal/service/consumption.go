// Package service 提供数据库服务接口，包括任务管理、消费记录生成等功能
package service

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"helpfly/internal/service/db"

	gf "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/xuri/excelize/v2"
)

// TaskPausePoint 任务暂停点结构体
// 用于存储任务暂停时的状态信息，支持后续恢复任务
type TaskPausePoint struct {
	TaskID           string  `json:"taskId"`           // 任务ID
	LastSerialNumber int     `json:"lastSerialNumber"` // 最后生成的序号（用于保持序列号连续性）
	CurrentBalance   float64 `json:"currentBalance"`   // 当前余额（用于保持余额连续性）
	CurrentProgress  int     `json:"currentProgress"`  // 当前进度（已完成的记录数，从1开始）
	Percent          float64 `json:"percent"`          // 完成百分比
	PausedAt         string  `json:"pausedAt"`         // 暂停时间
}

// ConsumptionRecord 消费记录结构体
// 表示一条银行交易消费记录，包含交易的所有详细信息
type ConsumptionRecord struct {
	Key                   string  `json:"key"`                   // 数据库存储键，唯一标识符
	TaskID                string  `json:"taskId"`                // 关联的任务ID
	TradeDate             string  `json:"tradeDate"`             // 交易日期时间（格式：2006-01-02 15:04:05）
	Account               string  `json:"account"`               // 账号
	StorageType           string  `json:"storageType"`           // 储种（活期/定期）
	SerialNumber          string  `json:"serialNumber"`          // 序号（允许任意字符串）
	Currency              string  `json:"currency"`              // 币种（人民币/美元）
	CashOrRemit           string  `json:"cashOrRemit"`           // 钞汇类型（钞/汇）
	Summary               string  `json:"summary"`               // 摘要（消费/工资/银联消费等）
	Region                string  `json:"region"`                // 地区（北京/上海/广州等）
	IncomeOrExpenseAmount float64 `json:"incomeOrExpenseAmount"` // 收入/支出金额（正数为收入，负数为支出）
	Balance               float64 `json:"balance"`               // 余额（累加计算）
	Channel               string  `json:"channel"`               // 渠道（快捷支付/网上银行等）
}

// ExportRecord 导出记录结构体
// 表示一条PDF导出记录
type ExportRecord struct {
	Key       string    `json:"key"`       // 数据库存储键，唯一标识符
	TaskID    string    `json:"taskId"`    // 关联的任务ID
	FilePath  string    `json:"filePath"`  // PDF文件保存路径
	CreatedAt time.Time `json:"createdAt"` // 导出时间
}

// 验证
func (g *DbService) validateTask(task TaskRecord) error {
	if task.TransactionCount == nil || *task.TransactionCount <= 0 {
		return fmt.Errorf("交易笔数未设置，请先填写大于0的交易笔数")
	}
	if len(task.TransactionDateRange) != 2 {
		return fmt.Errorf("交易日期区间无效")
	}
	return nil
}

// 验证任务是否正在运行
func (g *DbService) isTaskRunning(taskID string) bool {
	existingProgress := globalTaskManager.getProgress(taskID)
	if existingProgress != nil && existingProgress.IsRunning {
		return true
	}
	return false
}

// 解析暂停点数据
func pausePointFromDB(point db.PausePoint) *TaskPausePoint {
	return &TaskPausePoint{TaskID: point.TaskID, LastSerialNumber: point.LastSerialNumber,
		CurrentBalance: point.CurrentBalance, CurrentProgress: point.CurrentProgress,
		Percent: point.Percent, PausedAt: point.PausedAt}
}

func pausePointToDB(point TaskPausePoint) db.PausePoint {
	return db.PausePoint{TaskID: point.TaskID, LastSerialNumber: point.LastSerialNumber,
		CurrentBalance: point.CurrentBalance, CurrentProgress: point.CurrentProgress,
		Percent: point.Percent, PausedAt: point.PausedAt}
}

func consumptionToDB(record ConsumptionRecord) db.Consumption {
	return db.Consumption{Key: record.Key, TaskID: record.TaskID, TradeDate: record.TradeDate,
		Account: record.Account, StorageType: record.StorageType, SerialNumber: record.SerialNumber,
		Currency: record.Currency, CashOrRemit: record.CashOrRemit, Summary: record.Summary,
		Region: record.Region, IncomeOrExpenseAmount: record.IncomeOrExpenseAmount,
		Balance: record.Balance, Channel: record.Channel}
}

func consumptionFromDB(record db.Consumption) ConsumptionRecord {
	return ConsumptionRecord{Key: record.Key, TaskID: record.TaskID, TradeDate: record.TradeDate,
		Account: record.Account, StorageType: record.StorageType, SerialNumber: record.SerialNumber,
		Currency: record.Currency, CashOrRemit: record.CashOrRemit, Summary: record.Summary,
		Region: record.Region, IncomeOrExpenseAmount: record.IncomeOrExpenseAmount,
		Balance: record.Balance, Channel: record.Channel}
}

func exportToDB(record ExportRecord) db.Export {
	return db.Export{Key: record.Key, TaskID: record.TaskID, FilePath: record.FilePath, CreatedAt: record.CreatedAt}
}

func exportFromDB(record db.Export) ExportRecord {
	return ExportRecord{Key: record.Key, TaskID: record.TaskID, FilePath: record.FilePath, CreatedAt: record.CreatedAt}
}

// GenerateConsumptions 根据任务参数生成消费记录
// 该方法异步执行，实时更新任务进度，支持并发执行多个任务
// 参数:
//   - task: 任务记录，包含生成消费记录所需的所有参数
//   - doType: 操作类型，"pause"表示暂停任务，"resume"表示恢复任务，空字符串表示正常启动任务
//
// 返回值:
//   - string: 成功返回"任务已启动"、"任务已暂停"或"任务已恢复"
//   - error: 如果参数无效或任务已在运行返回错误
//
// 生成逻辑:
//   - 交易日期：在TransactionDateRange范围内随机生成
//   - 渠道：从Channels列表中随机选择一个
//   - 摘要：从Summaries列表中随机选择一个
//   - 金额：根据AmountType规则生成（规则一：0-1000正数，规则二：-1000-0负数，默认：-1000到1000随机）
//   - 余额：从OpeningBalance开始，每笔交易累加
func (g *DbService) GenerateConsumptions(task TaskRecord, doType string) (string, error) {

	if err := g.validateTask(task); err != nil {
		return "", err
	}
	var pausePoint *TaskPausePoint
	currentProgress := 0
	percent := 0.0

	// 初始化进度信息
	progress := &TaskProgress{
		TaskID:     task.TaskID,
		Current:    currentProgress,
		Total:      *task.TransactionCount,
		Percent:    percent,
		IsRunning:  false,
		IsCanceled: false,
	}

	// 根据操作类型执行不同的处理
	switch doType {
	case "start":
		// Check if task is already running
		// 如果任务状态是未开始（status=0），即使任务管理器中有进度信息，也允许启动
		// 这样可以处理删除消费记录后重新启动的情况
		if task.Status != 0 && g.isTaskRunning(task.TaskID) {
			return "", fmt.Errorf("任务正在执行中")
		}
		// 如果任务状态是未开始，清理可能存在的残留进度信息
		if task.Status == 0 {
			globalTaskManager.removeTask(task.TaskID)
			// 清理可能存在的暂停点
			_ = db.DeletePausePoint(gctx.New(), task.TaskID)
		}
		progress.IsRunning = true
		task.Status = 1 // 1=进行中
	case "pause":
		// 暂停操作 - 取消正在运行的任务
		if globalTaskManager.cancelTask(task.TaskID) {
			gf.Log().Info(gctx.New(), "任务开始暂停pause")
			// 等待一小段时间，确保暂停点已保存
			time.Sleep(200 * time.Millisecond)

			// 从暂停点获取进度信息，更新到任务管理器
			pausePointData, err := db.GetPausePoint(gctx.New(), task.TaskID)
			if err == nil {
				pausePoint := pausePointFromDB(pausePointData)
				if pausePoint != nil {
					// 更新进度信息到任务管理器，以便前端查询
					progress.Current = pausePoint.CurrentProgress
					progress.Percent = pausePoint.Percent
					progress.IsRunning = false
					globalTaskManager.setProgress(task.TaskID, progress)
				}
			}

			task.Status = 3 // 3=已暂停
			// 更新任务状态
			_ = db.UpdateTaskStatus(gctx.New(), task.TaskID, task.Status)
			return "任务已暂停", nil
		}
		return "", fmt.Errorf("任务不存在或未在运行")
	case "resume":

		var pausePointData db.PausePoint
		var err error
		maxRetries := 20                     // 增加重试次数
		retryDelay := 100 * time.Millisecond // 增加重试间隔

		for range maxRetries {
			pausePointData, err = db.GetPausePoint(gctx.New(), task.TaskID)
			if err == nil {
				gf.Log().Info(gctx.New(), "找到任务暂停点，开始恢复")
				break
			}
			// 如果暂停点不存在，等待一小段时间后重试
			time.Sleep(retryDelay)
		}

		if err != nil {
			gf.Log().Warning(gctx.New(), "任务没有暂停点，无法恢复", err)
			return "", fmt.Errorf("任务没有暂停点，无法恢复")
		}

		// 验证暂停点数据是否有效
		// var pausePoint TaskPausePoint
		pausePoint = pausePointFromDB(pausePointData)
		// 加载暂停点信息，设置初始进度
		if pausePoint.CurrentProgress > 0 {
			currentProgress = pausePoint.CurrentProgress
			percent = pausePoint.Percent
		}
		task.Status = 1 // 1=进行中
		progress.Current = currentProgress
		progress.Percent = percent
		progress.IsRunning = true

	default:
		// 正常启动 - 检查任务是否已在运行
		existingProgress := globalTaskManager.getProgress(task.TaskID)
		if existingProgress != nil && existingProgress.IsRunning {
			return "", fmt.Errorf("任务正在执行中")
		}
		// 清除可能存在的暂停点（全新开始）
		_ = db.DeletePausePoint(gctx.New(), task.TaskID)
	}

	globalTaskManager.setProgress(task.TaskID, progress)

	// pause 操作不需要启动新的 goroutine，因为只是取消现有任务
	if doType == "pause" {
		return "任务已暂停", nil
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	globalTaskManager.setCancelFunc(task.TaskID, cancel)

	// Run generation in goroutine
	go func() {
		// defer 确保即使发生 panic 也能执行清理工作
		defer func() {
			// 延迟清理任务管理器中的进度信息（保留一段时间以便前端查询）
			go func() {
				time.Sleep(5 * time.Minute)
				globalTaskManager.removeTask(task.TaskID)
			}()
		}()

		err := g.generateConsumptionsWithProgress(ctx, task, progress, pausePoint)

		// 根据错误类型决定如何处理
		if err != nil {
			progress.Error = err.Error()
			progress.IsRunning = false
			globalTaskManager.setProgress(task.TaskID, progress)

			// 如果是暂停操作，不删除暂停点，保持暂停状态
			if err.Error() == "任务已暂停" {
				// 任务被暂停，保留暂停点，更新任务状态为已暂停
				task.Status = 3 // 3=已暂停
				_ = db.UpdateTaskStatus(gctx.New(), task.TaskID, task.Status)
				// 注意：暂停时不删除暂停点，以便后续恢复
			} else {
				// 其他错误（如任务失败），删除暂停点，更新状态为已取消
				_ = db.DeletePausePoint(gctx.New(), task.TaskID)
				task.Status = 4 // 4=已取消（表示出错）
				_ = db.UpdateTaskStatus(gctx.New(), task.TaskID, task.Status)
			}
		} else {
			// 任务成功完成，删除暂停点，更新状态为已完成
			progress.IsRunning = false
			globalTaskManager.setProgress(task.TaskID, progress)
			_ = db.DeletePausePoint(gctx.New(), task.TaskID)
			task.Status = 2 // 2=已完成
			_ = db.UpdateTaskStatus(gctx.New(), task.TaskID, task.Status)
		}
	}()

	// 更新任务状态
	_ = db.UpdateTaskStatus(gctx.New(), task.TaskID, task.Status)

	if doType == "resume" {
		return "任务已恢复", nil
	}
	return "任务已启动", nil
}

// generateConsumptionsWithProgress 生成消费记录并实时更新进度
// 这是实际执行生成逻辑的内部方法，在goroutine中运行
// 参数:
//   - ctx: 上下文，用于支持任务取消
//   - task: 任务记录
//   - progress: 进度信息指针，会实时更新
//
// 返回值:
//   - error: 如果生成失败或任务被取消返回错误
func (g *DbService) generateConsumptionsWithProgress(ctx context.Context, task TaskRecord, progress *TaskProgress, pausePoint *TaskPausePoint) error {
	// 初始化变量
	var startIndex = 0
	var resumeFromPause = false
	var balance = 0.0
	var lastAmount = 0.0 // 用于规则二的递减逻辑

	// 如果传入了暂停点，说明是从暂停点恢复
	if pausePoint != nil {
		resumeFromPause = true
		// CurrentProgress 表示已完成的记录数，startIndex 应该是下一个要处理的索引
		// 因为索引从0开始，所以 startIndex = CurrentProgress
		startIndex = pausePoint.CurrentProgress
		balance = pausePoint.CurrentBalance
		// 规则二从暂停点恢复时，根据当前余额和剩余记录数重新计算初始金额
		if task.AmountType == "规则二" && balance > 0 {
			remainingCount := *task.TransactionCount - pausePoint.CurrentProgress
			if remainingCount > 0 {
				// 使用当前余额的某个比例作为初始值，确保能完成剩余记录
				lastAmount = balance * 0.8 / float64(remainingCount)
			}
		}
	}

	// Parse date range
	startDate, err := time.Parse("2006-01-02", task.TransactionDateRange[0])
	if err != nil {
		return fmt.Errorf("解析开始日期失败: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", task.TransactionDateRange[1])
	if err != nil {
		return fmt.Errorf("解析结束日期失败: %v", err)
	}
	if startDate.After(endDate) {
		return fmt.Errorf("开始日期不能晚于结束日期")
	}

	// Calculate date range in seconds
	dateRangeSeconds := int(endDate.Sub(startDate).Seconds())
	if dateRangeSeconds <= 0 {
		dateRangeSeconds = 1
	}

	if task.OpeningBalance != nil {
		balance = *task.OpeningBalance
	}

	// 如果是从暂停点恢复，使用暂停点的余额
	if resumeFromPause {
		balance = pausePoint.CurrentBalance
	}

	// Go 1.20 起默认使用全局随机源，无需手动 Seed
	// rand.Seed(time.Now().UnixNano())

	// 从系统参数读取地区选项（任务字段中没有地区选项）
	params, err := g.GetSystemParameters()
	var regions []string
	if err == nil && len(params.Regions) > 0 {
		regions = params.Regions
	}
	if len(regions) == 0 {
		return fmt.Errorf("地区选项为空，请先在系统参数中配置地区选项")
	}

	for i := startIndex; i < *task.TransactionCount; i++ {
		// Check for cancellation - 在循环开始处检查，确保能及时响应暂停
		select {
		case <-ctx.Done():
			gf.Log().Info(gctx.New(), "任务开始被暂停")
			// 保存暂停点信息
			// i 是当前要处理的索引，但还没有处理
			// 如果 i == startIndex，说明还没有处理任何新记录
			// 如果 i > startIndex，说明已经处理了从 startIndex 到 i-1 的记录（共 i - startIndex 条）
			// 计算当前进度（序号固定，不再用于恢复）
			var currentProgress int

			if i == startIndex {
				// 还没有处理任何新记录
				if resumeFromPause && pausePoint != nil {
					// 从暂停点恢复但还没处理新记录，使用暂停点的原始值
					currentProgress = pausePoint.CurrentProgress
				} else {
					// 新任务但还没处理任何记录
					currentProgress = 0
				}
			} else {
				// 已经处理了从 startIndex 到 i-1 的记录（共 i - startIndex 条）
				if resumeFromPause && pausePoint != nil {
					// 从暂停点恢复：当前进度 = 暂停点进度 + (i - startIndex)
					currentProgress = pausePoint.CurrentProgress + (i - startIndex)
				} else {
					// 新任务：当前进度 = 已完成的记录数 = i（因为索引从0开始）
					currentProgress = i
				}
			}

			// 计算完成百分比
			percent := float64(currentProgress) / float64(progress.Total) * 100

			pausePoint := TaskPausePoint{
				TaskID: task.TaskID,
				// 序号固定不再递增，这里保持 0 即可
				LastSerialNumber: 0,
				CurrentBalance:   balance,
				CurrentProgress:  currentProgress,
				Percent:          percent,
				PausedAt:         time.Now().Format("2006-01-02 15:04:05"),
			}

			if err := db.SavePausePoint(gctx.New(), pausePointToDB(pausePoint)); err != nil {
				return fmt.Errorf("保存暂停点失败: %v", err)
			}
			return fmt.Errorf("任务已暂停")
		default:
		}

		// Random date within range
		randomSeconds := rand.Intn(dateRangeSeconds)
		tradeDate := startDate.Add(time.Duration(randomSeconds) * time.Second)

		// Random channel (从任务字段的Channels中随机选择)
		channel := task.CashExchange
		if len(task.Channels) > 0 {
			channel = task.Channels[rand.Intn(len(task.Channels))]
		}

		// Random summary (从任务字段的Summaries中随机选择)
		summary := "消费"
		if len(task.Summaries) > 0 {
			summary = task.Summaries[rand.Intn(len(task.Summaries))]
		}

		// Random amount (positive or negative)
		amount := rand.Float64()*2000 - 1000 // 默认范围：-1000 至 1000
		switch task.AmountType {
		case "规则一":
			amount = rand.Float64() * 1000 // 0 至 1000
		case "规则二":
			// 规则二：金额递减且保证为正数（支出的绝对值递减）
			if i == startIndex {
				// 第一次生成，从余额的某个比例开始（例如余额的80%）
				if balance > 0 {
					lastAmount = balance * 0.8
				} else {
					// 如果余额为0或负数，使用固定初始值
					lastAmount = 1000.0
				}
			} else {
				// 递减：每次减少10%-30%
				decreaseRate := 0.1 + rand.Float64()*0.2 // 10% 到 30%
				lastAmount = lastAmount * (1 - decreaseRate)
			}
			// 确保金额为正数且不超过当前余额
			if lastAmount <= 0 {
				// 金额变为负数或0，停止任务
				return fmt.Errorf("规则二：金额已递减至0或负数，任务停止")
			}
			if lastAmount > balance {
				// 如果金额超过余额，使用余额
				lastAmount = balance
			}
			// 规则二是支出，所以是负数
			amount = -lastAmount
		}

		balance += amount

		// 规则二：如果余额变成负数，停止任务
		if task.AmountType == "规则二" && balance < 0 {
			return fmt.Errorf("规则二：余额已变为负数，任务停止")
		}

		// 计算序号：使用任务中配置的字符串，不再在后端递增
		serialNumber := task.StartSerialNumber
		if serialNumber == "" {
			serialNumber = "0000"
		}

		record := ConsumptionRecord{
			Key:                   fmt.Sprintf("CONS%v_%s", time.Now().UnixNano(), serialNumber),
			TaskID:                task.TaskID,
			TradeDate:             tradeDate.Format("2006-01-02 15:04:05"),
			Account:               task.Account,
			StorageType:           task.DepositType,
			SerialNumber:          serialNumber,
			Currency:              task.Currency,
			CashOrRemit:           task.CashExchange,
			Summary:               summary,
			Region:                regions[rand.Intn(len(regions))],
			IncomeOrExpenseAmount: amount,
			Balance:               balance,
			Channel:               channel,
		}

		// Save record
		if err := db.SaveConsumption(gctx.New(), consumptionToDB(record)); err != nil {
			return fmt.Errorf("保存失败: %v", err)
		}

		// 更新进度
		// 当前进度 = 已完成的记录数（已保存的记录数）
		if resumeFromPause {
			// 从暂停点恢复：当前进度 = 暂停点进度 + 本次已处理的记录数
			// 本次已处理的记录数 = i - startIndex + 1（从 startIndex 到 i，共 i - startIndex + 1 条）
			progress.Current = pausePoint.CurrentProgress + (i - startIndex + 1)
		} else {
			// 新任务：当前进度 = i + 1（因为索引从0开始，已完成 i+1 条）
			progress.Current = i + 1
		}
		progress.Percent = float64(progress.Current) / float64(progress.Total) * 100
		globalTaskManager.setProgress(task.TaskID, progress)

		// Small delay to simulate processing
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// GetTaskProgress 获取指定任务的当前进度
// 参数:
//   - taskID: 任务ID
//
// 返回值:
//   - TaskProgress: 任务进度信息
//   - error: 如果任务不存在返回错误
func (g *DbService) GetTaskProgress(taskID string) (*TaskProgress, error) {
	// 首先尝试从任务管理器获取进度（正在运行的任务）
	progress := globalTaskManager.getProgress(taskID)
	if progress != nil {
		return progress, nil
	}

	// 获取任务信息
	taskData, err := db.GetTask(gctx.New(), taskID)
	if err != nil {
		return nil, fmt.Errorf("任务不存在")
	}

	task := taskFromDB(taskData)

	if task.TransactionCount == nil {
		return nil, fmt.Errorf("任务交易笔数无效")
	}

	// 如果任务管理器中没有进度信息，尝试从暂停点获取（暂停的任务）
	pausePointData, err := db.GetPausePoint(gctx.New(), taskID)
	if err == nil {
		pausePoint := pausePointFromDB(pausePointData)
		if pausePoint != nil {
			// 从暂停点构建进度信息
			progress = &TaskProgress{
				TaskID:     taskID,
				Current:    pausePoint.CurrentProgress,
				Total:      *task.TransactionCount,
				Percent:    pausePoint.Percent,
				IsRunning:  false,
				IsCanceled: false,
			}
			return progress, nil
		}
	}

	// 如果任务已完成（status=2），从消费记录数量计算进度
	if task.Status == 2 {
		consumptions, err := g.ListConsumptions(taskID)
		if err == nil {
			currentProgress := len(consumptions)
			percent := float64(currentProgress) / float64(*task.TransactionCount) * 100
			if percent > 100 {
				percent = 100
			}
			progress = &TaskProgress{
				TaskID:     taskID,
				Current:    currentProgress,
				Total:      *task.TransactionCount,
				Percent:    percent,
				IsRunning:  false,
				IsCanceled: false,
			}
			return progress, nil
		}
		// 如果无法获取消费记录，假设已完成
		progress = &TaskProgress{
			TaskID:     taskID,
			Current:    *task.TransactionCount,
			Total:      *task.TransactionCount,
			Percent:    100,
			IsRunning:  false,
			IsCanceled: false,
		}
		return progress, nil
	}

	return nil, fmt.Errorf("任务不存在")
}

// CancelTask 取消正在运行的任务
// 参数:
//   - taskID: 任务ID
//
// 返回值:
//   - string: 成功返回"任务已取消"
//   - error: 如果任务不存在或未在运行返回错误
//
// 注意: 取消后会延迟1分钟清理进度信息，以便前端获取最终状态
func (g *DbService) CancelTask(taskID string) (string, error) {
	if globalTaskManager.cancelTask(taskID) {
		// Clean up after a delay to allow frontend to fetch final progress
		go func() {
			time.Sleep(1 * time.Minute)
			globalTaskManager.removeTask(taskID)
		}()
		return "任务已取消", nil
	}
	return "", fmt.Errorf("任务不存在或未在运行")
}

// GetDeleteProgress 获取指定任务的删除进度
// 参数:
//   - taskID: 任务ID
//
// 返回值:
//   - TaskProgress: 删除进度信息
//   - error: 如果任务不存在返回错误
func (g *DbService) GetDeleteProgress(taskID string) (*TaskProgress, error) {
	// 从删除管理器获取进度（正在删除的任务）
	progress := globalDeleteManager.getProgress(taskID)
	if progress != nil {
		return progress, nil
	}

	return nil, fmt.Errorf("删除进度不存在")
}

// ListConsumptions 获取消费记录列表
// 参数:
//   - taskID: 任务ID，如果为空字符串则返回所有任务的消费记录
//
// 返回值:
//   - []ConsumptionRecord: 消费记录列表，按序号升序排列
//   - error: 如果获取失败返回错误
func (g *DbService) ListConsumptions(taskID string) ([]ConsumptionRecord, error) {
	raw, err := db.ListConsumptions(gctx.New(), taskID)
	if err != nil {
		return nil, err
	}
	records := make([]ConsumptionRecord, 0, len(raw))
	for _, entry := range raw {
		records = append(records, consumptionFromDB(entry))
	}
	return records, nil
}

// SearchConsumptions 根据任务ID与关键词检索消费记录
// 参数:
//   - taskID: 任务ID
//   - keyword: 搜索关键词
//
// 返回值:
//   - []ConsumptionRecord: 匹配的消费记录列表
//   - error: 如果获取失败返回错误
func (g *DbService) SearchConsumptions(taskID, keyword string) ([]ConsumptionRecord, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return g.ListConsumptions(taskID)
	}

	raw, err := db.SearchConsumptions(gctx.New(), taskID, keyword)
	if err != nil {
		return nil, err
	}

	records := make([]ConsumptionRecord, 0, len(raw))
	for _, entry := range raw {
		records = append(records, consumptionFromDB(entry))
	}
	return records, nil
}

// DeleteConsumptions 批量删除消费记录
// 参数:
//   - keys: 要删除的消费记录Key列表
//
// 返回值:
//   - string: 成功返回"success"
//   - error: 如果删除失败返回错误
func (g *DbService) DeleteConsumptions(keys []string) (string, error) {
	if err := db.DeleteConsumptions(gctx.New(), keys); err != nil {
		return "", fmt.Errorf("删除消费记录失败: %v", err)
	}
	return "success", nil
}

// DeleteConsumptionsByTaskID 删除指定任务的所有消费记录
// 该方法异步执行，实时更新删除进度，支持并发执行多个删除任务
// 参数:
//   - taskID: 任务ID
//
// 返回值:
//   - string: 成功返回"删除已启动"
//   - error: 如果参数无效或删除已在运行返回错误
func (g *DbService) DeleteConsumptionsByTaskID(taskID string) (string, error) {
	if taskID == "" {
		return "", fmt.Errorf("taskID 不能为空")
	}

	// 检查是否已经在删除中
	existingProgress := globalDeleteManager.getProgress(taskID)
	if existingProgress != nil && existingProgress.IsRunning {
		return "", fmt.Errorf("删除操作正在执行中")
	}

	// 获取要删除的记录列表
	consumptions, err := g.ListConsumptions(taskID)
	if err != nil {
		return "", fmt.Errorf("获取消费记录失败: %v", err)
	}

	total := len(consumptions)
	if total == 0 {
		return "success", nil // 没有记录需要删除，直接返回成功
	}

	// 初始化删除进度信息
	progress := &TaskProgress{
		TaskID:     taskID,
		Current:    0,
		Total:      total,
		Percent:    0,
		IsRunning:  true,
		IsCanceled: false,
	}

	globalDeleteManager.setProgress(taskID, progress)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	globalDeleteManager.setCancelFunc(taskID, cancel)

	// Run deletion in goroutine
	go func() {
		// defer 确保即使发生 panic 也能执行清理工作
		defer func() {
			// 延迟清理删除管理器中的进度信息（保留一段时间以便前端查询）
			go func() {
				time.Sleep(5 * time.Minute)
				globalDeleteManager.removeDelete(taskID)
			}()
		}()

		err := g.deleteConsumptionsWithProgress(ctx, taskID, consumptions, progress)

		// 根据错误类型决定如何处理
		if err != nil {
			progress.Error = err.Error()
			progress.IsRunning = false
			globalDeleteManager.setProgress(taskID, progress)

			// 如果是取消操作，不清理进度信息
			if err.Error() == "删除已取消" {
				// 删除被取消，保留进度信息以便前端查询
			} else {
				// 其他错误（如删除失败），清理进度信息
				go func() {
					time.Sleep(1 * time.Minute)
					globalDeleteManager.removeDelete(taskID)
				}()
			}
		} else {
			// 删除成功完成
			progress.IsRunning = false
			progress.Current = total
			progress.Percent = 100
			globalDeleteManager.setProgress(taskID, progress)

			// 删除完成后，更新任务状态为未开始（status=0）
			_ = db.UpdateTaskStatus(gctx.New(), taskID, 0)
		}
	}()

	return "删除已启动", nil
}

// deleteConsumptionsWithProgress 删除消费记录并实时更新进度
// 这是实际执行删除逻辑的内部方法，在goroutine中运行
// 参数:
//   - ctx: 上下文，用于支持删除取消
//   - taskID: 任务ID
//   - consumptions: 要删除的消费记录列表
//   - progress: 进度信息指针，会实时更新
//
// 返回值:
//   - error: 如果删除失败或操作被取消返回错误
func (g *DbService) deleteConsumptionsWithProgress(ctx context.Context, taskID string, consumptions []ConsumptionRecord, progress *TaskProgress) error {
	const batchSize = 100
	for offset := 0; offset < len(consumptions); offset += batchSize {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("删除已取消")
		default:
		}

		end := offset + batchSize
		if end > len(consumptions) {
			end = len(consumptions)
		}
		keys := make([]string, 0, end-offset)
		for _, record := range consumptions[offset:end] {
			keys = append(keys, record.Key)
		}
		if err := db.DeleteConsumptions(ctx, keys); err != nil {
			return fmt.Errorf("删除失败: %v", err)
		}

		// 更新进度
		progress.Current = end
		progress.Percent = float64(progress.Current) / float64(progress.Total) * 100
		globalDeleteManager.setProgress(taskID, progress)

		// Small delay to simulate processing
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// consumptionMatchesKeyword 判断消费记录的任意字段是否包含关键词
func consumptionMatchesKeyword(record ConsumptionRecord, keywordLower string) bool {
	matchField := func(value any) bool {
		if value == nil {
			return false
		}
		text := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", value)))
		if text == "" {
			return false
		}
		return strings.Contains(text, keywordLower)
	}

	return matchField(record.Summary) ||
		matchField(record.Region) ||
		matchField(record.Channel) ||
		matchField(record.TradeDate) ||
		matchField(record.Account) ||
		matchField(record.StorageType) ||
		matchField(record.Currency) ||
		matchField(record.CashOrRemit) ||
		matchField(record.SerialNumber) ||
		matchField(record.IncomeOrExpenseAmount) ||
		matchField(record.Balance)
}

// ExportConsumptionsToPDF 将消费记录导出为 PDF 并保存到数据库
// 参数:
//   - taskID: 任务ID，如果为空则导出所有任务的消费记录
//   - keys: 要导出的消费记录Key列表，如果为空则导出所有记录
//
// 返回值:
//   - string: 导出记录的唯一标识符
//   - error: 如果导出失败返回错误
func (g *DbService) ExportConsumptionsToPDF(taskID string, keys []string) (string, error) {

	// 获取要导出的消费记录
	var records []ConsumptionRecord
	var err error

	if len(keys) > 0 {
		// 如果指定了 keys，获取对应的记录
		allRecords, err := g.ListConsumptions(taskID)
		if err != nil {
			return "", fmt.Errorf("获取消费记录失败: %v", err)
		}
		keyMap := make(map[string]bool)
		for _, key := range keys {
			keyMap[key] = true
		}
		for _, record := range allRecords {
			if keyMap[record.Key] {
				records = append(records, record)
			}
		}
	} else {
		// 否则获取所有记录
		records, err = g.ListConsumptions(taskID)
		if err != nil {
			return "", fmt.Errorf("获取消费记录失败: %v", err)
		}
	}

	if len(records) == 0 {
		return "", fmt.Errorf("没有可导出的消费记录")
	}

	// 获取任务信息
	var task TaskRecord
	if taskID != "" {
		taskData, err := db.GetTask(gctx.New(), taskID)
		if err == nil {
			task = taskFromDB(taskData)
		}
	}

	// 从系统参数读取自定义 PNG 水印图片
	params, err := g.GetSystemParameters()
	if err != nil {
		params = SystemParameters{}
	}
	watermarkPath := strings.TrimSpace(params.WatermarkPath)
	addWatermark := params.AddWatermark
	// 保留原有文字水印；如果配置了 PNG，则在文字水印之外再绘制图片水印。
	addTextWatermark := addWatermark
	if addWatermark && watermarkPath != "" {
		if strings.ToLower(filepath.Ext(watermarkPath)) != ".png" {
			return "", fmt.Errorf("水印图片仅支持 PNG 格式")
		}
		if info, statErr := os.Stat(watermarkPath); statErr != nil {
			return "", fmt.Errorf("水印图片不存在: %v", statErr)
		} else if info.IsDir() {
			return "", fmt.Errorf("水印路径不是图片文件")
		}
	}

	// 使用 gofpdf 创建 PDF（横向A4）
	pdf := gofpdf.New("L", "mm", "A4", "")

	const ptToMM = 25.4 / 72.0
	// 顶部空白
	topMargin := 12 * ptToMM //
	sideMargin := 15.0       // A4 常用左右边距约 15mm
	bottomMargin := 6.0      // 底部边距，减小以增加表格使用面积
	pdf.SetMargins(sideMargin, topMargin, sideMargin)
	// 设置自动分页的底部边距（较小值以增加表格可用空间）
	pdf.SetAutoPageBreak(true, bottomMargin)

	// 添加中文字体支持
	// 优先使用项目本地字体，然后尝试系统字体路径
	fontPaths := []struct {
		path     string
		fontName string
	}{
		{"./resource/font/SimHei.ttf", "simhei"}, // 项目本地字体 - 黑体（优先）
		{"./resource/font/SimSun.ttf", "simsun"}, // 项目本地字体 - 宋体
		{"./resource/font/SimKai.ttf", "simkai"}, // 项目本地字体 - 楷体
	}

	// 根据操作系统添加系统字体路径
	homeDir, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		// Windows 系统字体路径
		fontPaths = append(fontPaths,
			struct {
				path     string
				fontName string
			}{"C:/Windows/Fonts/simhei.ttf", "simhei"},
			struct {
				path     string
				fontName string
			}{"C:/Windows/Fonts/simsun.ttc", "simsun"},
			struct {
				path     string
				fontName string
			}{"C:/Windows/Fonts/simsun.ttf", "simsun"},
			struct {
				path     string
				fontName string
			}{"C:/Windows/Fonts/msyh.ttc", "microsoftyahei"},
			struct {
				path     string
				fontName string
			}{"C:/Windows/Fonts/simfang.ttf", "simfang"},
		)
	case "darwin":
		// macOS 系统字体路径
		if homeDir != "" {
			fontPaths = append(fontPaths,
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, "Library/Fonts/SimHei.ttf"), "simhei"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, "Library/Fonts/SimSun.ttf"), "simsun"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, "Library/Fonts/SimKai.ttf"), "simkai"},
			)
		}
		// macOS 系统级字体目录
		fontPaths = append(fontPaths,
			struct {
				path     string
				fontName string
			}{"/Library/Fonts/SimHei.ttf", "simhei"},
			struct {
				path     string
				fontName string
			}{"/Library/Fonts/SimSun.ttf", "simsun"},
			struct {
				path     string
				fontName string
			}{"/Library/Fonts/SimKai.ttf", "simkai"},
			// macOS 可能有的中文字体
			struct {
				path     string
				fontName string
			}{"/System/Library/Fonts/Supplemental/STHeiti Light.ttc", "simhei"},
			struct {
				path     string
				fontName string
			}{"/System/Library/Fonts/Supplemental/STSong.ttc", "simsun"},
		)
	case "linux":
		// Linux 系统字体路径
		if homeDir != "" {
			fontPaths = append(fontPaths,
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, ".fonts/SimHei.ttf"), "simhei"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, ".fonts/SimSun.ttf"), "simsun"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, ".fonts/SimKai.ttf"), "simkai"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, ".local/share/fonts/SimHei.ttf"), "simhei"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, ".local/share/fonts/SimSun.ttf"), "simsun"},
				struct {
					path     string
					fontName string
				}{filepath.Join(homeDir, ".local/share/fonts/SimKai.ttf"), "simkai"},
			)
		}
		// Linux 系统级字体目录
		fontPaths = append(fontPaths,
			struct {
				path     string
				fontName string
			}{"/usr/share/fonts/truetype/simhei/SimHei.ttf", "simhei"},
			struct {
				path     string
				fontName string
			}{"/usr/share/fonts/truetype/simsun/SimSun.ttf", "simsun"},
			struct {
				path     string
				fontName string
			}{"/usr/share/fonts/truetype/simkai/SimKai.ttf", "simkai"},
			struct {
				path     string
				fontName string
			}{"/usr/share/fonts/SimHei.ttf", "simhei"},
			struct {
				path     string
				fontName string
			}{"/usr/share/fonts/SimSun.ttf", "simsun"},
			struct {
				path     string
				fontName string
			}{"/usr/share/fonts/SimKai.ttf", "simkai"},
		)
	}

	var fontPath string
	var fontName string
	for _, font := range fontPaths {
		// 检查字体文件是否存在
		testPath := font.path
		// 如果是相对路径，尝试转换为绝对路径
		if !filepath.IsAbs(font.path) {
			if absPath, err := filepath.Abs(font.path); err == nil {
				testPath = absPath
			}
		}

		if _, err := os.Stat(testPath); err == nil {
			fontPath = testPath
			fontName = font.fontName
			// AddUTF8Font 接收的是字体目录下的文件名，而不是绝对路径。
			// 如果直接传入 /Volumes/...，gofpdf 内部再次使用 path.Join 时会丢失开头的斜杠。
			pdf.SetFontLocation(filepath.Dir(fontPath))
			fontFileName := filepath.Base(fontPath)
			// 添加 UTF-8 字体（gofpdf 会自动处理 TTF 和 TTC 文件）
			pdf.AddUTF8Font(fontName, "", fontFileName)
			// 尝试添加粗体版本（使用相同的字体文件）
			// 如果字体文件本身不支持粗体，gofpdf 会使用普通字体来模拟粗体
			pdf.AddUTF8Font(fontName, "B", fontFileName)
			// 找到可用字体后跳出循环
			break
		}
	}

	// 如果找到字体文件，使用中文字体
	if fontPath == "" || fontName == "" {
		return "", fmt.Errorf("未找到中文字体，无法生成PDF")
	}

	// 添加第一页（必须在字体加载之后）
	pdf.AddPage()

	// 生成每页唯一码的函数（基于页码+时间戳+随机数）
	generatePageUniqueCode := func(pageNum int) string {
		// 使用页码、当前时间戳和随机数生成唯一码
		rand.Seed(time.Now().UnixNano() + int64(pageNum))
		randomPart := rand.Intn(65536) // 0-65535 的随机数
		// 格式：2位大写字母 + 10位十六进制（大写），共12位，类似 "BA09A8D42049"
		// 前2位：基于页码和随机数的字母组合
		letter1 := byte('A' + (pageNum*7+randomPart)%26)
		letter2 := byte('A' + (pageNum*13+randomPart*3)%26)
		// 后10位：时间戳和随机数的十六进制
		hexPart := fmt.Sprintf("%010X", (time.Now().UnixNano()%10000000000)+int64(randomPart))
		code := fmt.Sprintf("%c%c%s", letter1, letter2, hexPart)
		return code
	}

	// 绘制水印函数：在每一页的最底层绘制水印（完全不影响PDF状态和分页）
	drawWatermark := func(uniqueCode string) {
		if !addTextWatermark {
			return
		}

		// 保存当前所有状态
		currentX := pdf.GetX()
		currentY := pdf.GetY()
		currentTextR, currentTextG, currentTextB := pdf.GetTextColor()
		currentAutoPageBreak, currentBreakMargin := pdf.GetAutoPageBreak()

		// 获取页面尺寸和当前页码
		pageWidth, pageHeight := pdf.GetPageSize()
		currentPage := pdf.PageNo()

		// 临时禁用自动分页，防止水印绘制时触发分页
		pdf.SetAutoPageBreak(false, 0)

		// 水印文字使用当前任务选择的银行；没有任务银行时兼容旧版默认文案。
		watermarkBank := strings.TrimSpace(task.Bank)
		if watermarkBank == "" {
			watermarkBank = "中国工商银行"
		}
		watermarkText := fmt.Sprintf("%s %s %s %s", watermarkBank, uniqueCode, time.Now().Format("2006-01-02 15:04:05"), uniqueCode)

		// 设置水印字体和颜色（淡的水印色，可见但不遮挡内容）
		pdf.SetFont(fontName, "", 20)
		// 使用更浅的灰色来模拟透明度效果（值越大越透明）
		pdf.SetTextColor(230, 230, 230) // 更透明的水印颜色（模拟透明度）

		// 计算水印文字宽度和高度
		textWidth := pdf.GetStringWidth(watermarkText)
		textHeight := 20.0 * ptToMM

		// 行间距（单位：mm），继续减小
		lineSpacing := 33.0

		// 计算需要绘制的列数，确保覆盖整个页面
		// 考虑30度旋转后，需要计算对角线长度
		diagonalLength := math.Sqrt(pageWidth*pageWidth + pageHeight*pageHeight)
		// 计算旋转后的有效宽度（30度旋转后，宽度会相应变化）
		rotatedWidth := textWidth*math.Cos(30.0*math.Pi/180.0) + textHeight*math.Sin(30.0*math.Pi/180.0)

		// 计算需要的列数（行数固定为4行）
		cols := int(math.Ceil((diagonalLength + rotatedWidth) / (textWidth * 1.5)))

		// 从页面左侧开始绘制水印
		// 第一行y距离顶部50%的高度开始，再下降100mm，x起点从页面左侧最边缘开始（x=0）
		// 每行的x起点偏移100mm
		baseStartX := 12.0           // 第一行从x=10开始（页面左侧最边缘）
		rowOffset := 35.0            // 每行x起点向右偏移100mm
		startY := pageHeight*0.5 + 8 // 第一行y距离顶部50%的高度，再下降100mm

		// 绘制多行多列水印
		// 使用Text方法绘制，完全不影响PDF的当前位置和分页
		// 限制每页最多只绘制4行水印
		maxRowsPerPage := 4
		for row := 0; row < maxRowsPerPage; row++ {
			// 计算当前行的起始X坐标：第一行30mm，第二行80mm，第三行130mm...
			rowStartX := baseStartX + float64(row)*rowOffset
			for col := 0; col <= cols*2; col++ {
				// 计算当前水印位置（从左侧开始，30度方向重复）
				x := rowStartX + float64(col)*textWidth*1.5
				y := startY + float64(row)*lineSpacing

				// 应用旋转变换（顺时针30度，尾部高头部低）
				// 每个水印文字都使用独立的变换，确保不影响其他内容
				pdf.TransformBegin()
				pdf.TransformRotate(30, x, y)
				// 使用Text方法绘制水印，不改变PDF的实际当前位置
				pdf.Text(x, y, watermarkText)
				pdf.TransformEnd()
			}
		}

		// 完全恢复所有状态（确保不影响分页和后续绘制）
		pdf.SetXY(currentX, currentY)
		pdf.SetFont(fontName, "", 7) // 恢复为默认表格字体
		pdf.SetTextColor(currentTextR, currentTextG, currentTextB)
		pdf.SetAutoPageBreak(currentAutoPageBreak, currentBreakMargin)

		// 确保页码没有改变（防止意外翻页）
		if pdf.PageNo() != currentPage {
			// 如果页码改变了，说明有问题，需要修复
			// 但正常情况下不应该发生
		}
	}

	// 在第一页立即绘制水印（最底层，不遮挡任何内容）
	// 生成第一页的唯一码
	firstPageCode := generatePageUniqueCode(1)
	drawWatermark(firstPageCode)

	// 绘制底部图片水印：在每页底部居中绘制用户选择的 PNG 图片
	drawBottomImageWatermark := func() {
		if !addWatermark || watermarkPath == "" {
			return
		}
		// 获取页面尺寸
		pageWidth, pageHeight := pdf.GetPageSize()

		// 保存当前状态
		currentX := pdf.GetX()
		currentY := pdf.GetY()
		currentAutoPageBreak, currentBreakMargin := pdf.GetAutoPageBreak()
		currentPage := pdf.PageNo()
		currentTextR, currentTextG, currentTextB := pdf.GetTextColor()
		currentDrawR, currentDrawG, currentDrawB := pdf.GetDrawColor()
		currentLineWidth := pdf.GetLineWidth()
		// 保存当前字体大小（GetFontSize 返回字体大小和单位）
		currentFontSize, _ := pdf.GetFontSize()

		// 禁止自动分页，避免绘制过程中影响正文
		pdf.SetAutoPageBreak(false, 0)

		// 使用参数中保存的 PNG 图片作为底部水印
		imgPath := watermarkPath
		opt := gofpdf.ImageOptions{ImageType: "PNG"}
		info := pdf.RegisterImageOptions(imgPath, opt)
		if info != nil {
			imgWd, imgHt := info.Extent()
			// 固定图片宽度为 80mm，高度按原始比例自适应
			targetWd := 55.0
			scale := targetWd / imgWd
			w := targetWd
			h := imgHt * scale

			// 底部居中放置，距离页面底边留出 12mm，避免压住页脚
			x := (pageWidth - w) / 2
			y := pageHeight - h - 12.0

			// 直接绘制图片，保留 PNG 自带的透明背景。
			pdf.ImageOptions(imgPath, x, y, w, h, false, opt, 0, "")

		}

		// 恢复状态
		pdf.SetFont(fontName, "", currentFontSize)
		pdf.SetDrawColor(currentDrawR, currentDrawG, currentDrawB)
		pdf.SetTextColor(currentTextR, currentTextG, currentTextB)
		pdf.SetLineWidth(currentLineWidth)
		pdf.SetXY(currentX, currentY)
		pdf.SetAutoPageBreak(currentAutoPageBreak, currentBreakMargin)

		if pdf.PageNo() != currentPage {
			// 理论上不会发生
		}
	}

	// 在第一页绘制底部图片水印（作为背景层，不遮挡内容，参考drawWatermark的用法）
	drawBottomImageWatermark()

	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()

	// 先准备表格的总宽度与起始 X，用于对齐标题、银行信息行和表格
	colWidths := []float64{
		22, // 交易日期
		34, // 账号
		15, // 储种
		14, // 序号
		15, // 币种
		12, // 钞汇
		18, // 摘要
		15, // 地区
		36, // 收入/支出
		36, // 余额
		21, // 渠道
	}
	effectiveWidth := pageWidth - leftMargin - rightMargin
	totalWidth := 0.0
	for _, w := range colWidths {
		totalWidth += w
	}
	tableStartX := leftMargin
	if totalWidth < effectiveWidth {
		tableStartX = leftMargin + (effectiveWidth-totalWidth)/2
	}

	// 表头文本
	headers := []string{"交易日期", "账号", "储种", "序号", "币种", "钞汇", "摘要", "地区", "收入/支出", "余额", "渠道"}
	headerHeight := 6.0

	// 每一页的页眉（标题 + 银行信息行 + 表头），第一页和新页都会调用
	drawPageHeader := func() {
		// 1. 顶部标题
		title := "消费记录列表"
		if bankName := strings.TrimSpace(task.Bank); bankName != "" {
			// 使用当前导出任务中选择的银行，避免非工商银行回退成通用标题。
			title = bankName + "借记账户历史明细（电子版）"
		}
		pdf.SetY(topMargin)
		pdf.SetFont(fontName, "", 16)
		titleWidth := pdf.GetStringWidth(title)
		pdf.SetX((pageWidth - titleWidth) / 2)
		titleLineHeight := 16 * ptToMM
		// 标题不使用填充，让水印可见
		pdf.Cell(0, titleLineHeight, title)
		pdf.Ln(20 * ptToMM)

		// 2. 银行信息行（卡号 / 户名 / 起始日期 / 二维码），与表格左右对齐
		pdf.SetFont(fontName, "B", 12)

		// 当前 Y 作为银行信息行的 Y
		infoY := pdf.GetY()
		pdf.SetY(infoY)
		pdf.SetX(tableStartX)

		// 布局参数
		smallLineHeight := 8.0
		cardWidth := 53.0
		nameWidth := 66.0
		gapBetween := 5.0

		// 组装文案
		cardLabel := "卡号"
		// 顶部“卡号”优先使用任务配置的 CardNumber，表格中的账号仍然使用 Account
		cardValue := task.CardNumber
		if cardValue == "" && task.Account != "" {
			cardValue = task.Account
		}
		if cardValue == "" && len(records) > 0 {
			cardValue = records[0].Account
		}
		nameLabel := "户名:"
		nameValue := task.AccountName
		dateLabel := "起始日期:"
		dateValue := ""
		if len(task.DateRange) >= 2 {
			dateValue = task.DateRange[0] + "-" + task.DateRange[1]
		} else if len(records) > 0 {
			dates := make([]string, 0)
			for _, r := range records {
				if len(r.TradeDate) >= 10 {
					dates = append(dates, r.TradeDate[:10])
				}
			}
			if len(dates) > 0 {
				sort.Strings(dates)
				dateValue = dates[0] + "-" + dates[len(dates)-1]
			}
		}

		// 左侧区域：卡号 + 户名（不使用填充，让水印可见）
		cardText := cardLabel + " " + cardValue
		pdf.SetXY(tableStartX, infoY)
		pdf.CellFormat(cardWidth, smallLineHeight, cardText, "", 0, "L", false, 0, "")

		nameStartX := tableStartX + cardWidth + gapBetween
		pdf.SetXY(nameStartX, infoY)
		nameText := nameLabel + " " + nameValue
		pdf.CellFormat(nameWidth, smallLineHeight, nameText, "", 0, "L", false, 0, "")

		// 右侧区域：起始日期 + 二维码
		rightAreaStartX := tableStartX + cardWidth + gapBetween + nameWidth
		rightAreaWidth := totalWidth - (cardWidth + gapBetween + nameWidth)
		dateAreaWidth := rightAreaWidth / 2
		qrAreaWidth := rightAreaWidth - dateAreaWidth

		// 起始日期（不使用填充，让水印可见）
		pdf.SetXY(rightAreaStartX, infoY)
		dateText := dateLabel + " " + dateValue
		pdf.CellFormat(dateAreaWidth, smallLineHeight, dateText, "", 0, "L", false, 0, "")

		// 二维码区域
		qrSize := 16.0
		if qrSize > qrAreaWidth {
			qrSize = qrAreaWidth
		}
		qrX := tableStartX + totalWidth - qrSize
		qrY := infoY + smallLineHeight - qrSize - 3*ptToMM
		qrContent := task.QRCodeURL
		if qrContent == "" {
			qrContent = taskID
		}
		if qrContent == "" {
			qrContent = cardValue
		}
		if qrContent != "" {
			tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("qr_%d.png", time.Now().UnixNano()))
			if err := qrcode.WriteFile(qrContent, qrcode.Medium, 256, tmpFile); err == nil {
				defer os.Remove(tmpFile)
				opt := gofpdf.ImageOptions{ImageType: "PNG"}
				pdf.RegisterImageOptions(tmpFile, opt)
				pdf.ImageOptions(tmpFile, qrX, qrY, qrSize, qrSize, false, opt, 0, "")
			}
		}

		// 二维码左侧两行小字说明（不使用填充，让水印可见）
		pdf.SetFont(fontName, "", 7)
		textLineHeight := smallLineHeight / 3
		totalTextHeight := 2 * textLineHeight
		textStartY := qrY + (qrSize-totalTextHeight)/2
		textAreaRightX := qrX
		textAreaLeftX := rightAreaStartX + dateAreaWidth
		textAreaWidth := textAreaRightX - textAreaLeftX
		if textAreaWidth < 10 {
			textAreaWidth = 10
		}
		pdf.SetXY(textAreaRightX-textAreaWidth, textStartY)
		pdf.CellFormat(textAreaWidth, textLineHeight, "请扫描二维码", "", 0, "R", false, 0, "")
		pdf.SetXY(textAreaRightX-textAreaWidth, textStartY+textLineHeight)
		pdf.CellFormat(textAreaWidth, textLineHeight, "识别明细真伪", "", 0, "R", false, 0, "")

		// 银行信息行下方预留一点距离，然后画表头
		pdf.SetY(infoY + smallLineHeight - 3*ptToMM)

		// 3. 表头（不使用填充，让水印可见，但文字颜色足够深）
		pdf.SetFont(fontName, "", 7)
		pdf.SetTextColor(0, 0, 0) // 确保文字是黑色，即使有水印也能看清
		pdf.SetX(tableStartX)
		for i, header := range headers {
			pdf.CellFormat(colWidths[i], headerHeight, header, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(headerHeight)
	}

	// 先画第一页的页眉（drawPageHeader 内部会先绘制水印，然后绘制页眉内容）
	drawPageHeader()

	// 填充数据（每页最多 25 条）
	const rowsPerPage = 25
	rowHeight := 6.4 // 每行高度约 6mm
	rowsOnPage := 0
	pageIncomeTotal := 0.0
	pageExpenseTotal := 0.0
	pageIndex := 1
	totalPages := int(math.Ceil(float64(len(records)) / float64(rowsPerPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	// 设置表格内容字体约 8pt（不使用填充，让水印可见，但文字颜色足够深）
	pdf.SetFont(fontName, "", 7)
	pdf.SetTextColor(0, 0, 0) // 确保文字是黑色，即使有水印也能看清

	printFooter := func(rows int, income, expense float64, pageIdx int) {
		if rows == 0 {
			return
		}
		expense = math.Abs(expense)
		pdf.SetFont(fontName, "", 7)
		// 略增与表格之间的间距，让 footer 第一行不要太贴表格
		pdf.SetY(pdf.GetY() + 2)

		leftText := fmt.Sprintf("本页支出算术合计: %.2f", expense)
		centerText := fmt.Sprintf("本页收入算术合计: %.2f", income)
		rightText := fmt.Sprintf("下单时间: %s", time.Now().Format("2006-01-02 15:04:05"))

		// 第一行高度由 6 缩小为 4（不使用填充，让水印可见）
		lineH1 := 4.0
		pdf.SetX(tableStartX)
		pdf.CellFormat(totalWidth/3, lineH1, leftText, "", 0, "L", false, 0, "")
		pdf.CellFormat(totalWidth/3, lineH1, centerText, "", 0, "C", false, 0, "")
		pdf.CellFormat(totalWidth/3, lineH1, rightText, "", 0, "R", false, 0, "")
		pdf.Ln(lineH1 + 2)

		leftBottom := fmt.Sprintf("本页交易笔数: %d", rows)
		rightBottom := fmt.Sprintf("共 %d 页 第 %d 页", totalPages, pageIdx)

		// 第二行高度也改为 4，并减少末尾额外空白（不使用填充，让水印可见）
		lineH2 := 4.0
		pdf.SetX(tableStartX)
		pdf.CellFormat(totalWidth/2, lineH2, leftBottom, "", 0, "L", false, 0, "")
		pdf.CellFormat(totalWidth/2, lineH2, rightBottom, "", 0, "R", false, 0, "")
		pdf.Ln(3)
	}

	for _, record := range records {
		// 如果本页已满 25 条，则换页并重绘页眉（标题 + 银行信息行 + 表头）
		if rowsOnPage >= rowsPerPage {
			printFooter(rowsOnPage, pageIncomeTotal, pageExpenseTotal, pageIndex)
			// 换页
			pdf.AddPage()
			// 生成当前页的唯一码
			pageCode := generatePageUniqueCode(pageIndex + 1)
			// 在新页立即绘制水印（最底层，在所有内容之前，不遮挡任何内容）
			drawWatermark(pageCode)
			// 在新页绘制底部图片水印（作为背景层，不遮挡内容，参考drawWatermark的用法）
			drawBottomImageWatermark()
			// 然后绘制页眉内容
			drawPageHeader()
			rowsOnPage = 0
			pageIncomeTotal = 0
			pageExpenseTotal = 0
			pageIndex++
		}

		// 格式化金额
		amountStr := fmt.Sprintf("%.2f", record.IncomeOrExpenseAmount)
		if record.IncomeOrExpenseAmount >= 0 {
			amountStr = "+" + amountStr
		}
		balanceStr := fmt.Sprintf("%.2f", record.Balance)

		// 设置表格行位置（居中）
		pdf.SetX(tableStartX)

		// 填充行数据（统一行高）
		tradeDateText := record.TradeDate
		rowData := []string{
			tradeDateText,
			record.Account,
			record.StorageType,
			record.SerialNumber,
			record.Currency,
			record.CashOrRemit,
			record.Summary,
			record.Region,
			amountStr,
			balanceStr,
			record.Channel,
		}
		for i, data := range rowData {
			if i == 0 {
				// 交易日期列：不做省略号截断，拆成日期/时间两行
				tradeData := data
				// 交易日期：拆成日期/时间两行，如果有空格（日期与时间）
				parts := strings.SplitN(tradeData, " ", 2)
				dateLine := parts[0]
				timeLine := ""
				if len(parts) > 1 {
					timeLine = parts[1]
				}
				cellLeft := pdf.GetX()
				cellTop := pdf.GetY()
				// 不使用填充，让水印可见，但文字颜色足够深
				pdf.MultiCell(colWidths[i], rowHeight/2, dateLine, "LR", "C", false)
				if timeLine != "" {
					pdf.SetXY(cellLeft, cellTop+rowHeight/2)
					pdf.MultiCell(colWidths[i], rowHeight/2, timeLine, "LRB", "C", false)
				} else {
					// 如果没有时间，只画一个单元格边框（单行居中）
					pdf.SetXY(cellLeft, cellTop)
					pdf.CellFormat(colWidths[i], rowHeight, dateLine, "1", 0, "C", false, 0, "")
				}
				pdf.SetXY(cellLeft+colWidths[i], cellTop)
			} else {
				// 其他列：不使用填充，让水印可见，但文字颜色足够深
				pdf.CellFormat(colWidths[i], rowHeight, data, "1", 0, "C", false, 0, "")
			}
		}

		value := record.IncomeOrExpenseAmount
		if value >= 0 {
			pageIncomeTotal += value
		} else {
			pageExpenseTotal += value
		}

		rowsOnPage++

		// 移动到下一行
		pdf.Ln(rowHeight)
	}

	printFooter(rowsOnPage, pageIncomeTotal, pageExpenseTotal, pageIndex)
	// 注意：最后一页的底部图片水印已经在 drawPageHeader() 中绘制了（作为背景层，不遮挡内容）

	// 生成 PDF 字节数据
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return "", fmt.Errorf("生成 PDF 失败: %v", err)
	}
	pdfData := buf.Bytes()

	// 获取导出路径
	exportDir, err := g.GetExportPath()
	if err != nil {
		return "", fmt.Errorf("获取导出路径失败: %v", err)
	}

	// 生成文件名（使用任务ID作为标题）
	fileName := fmt.Sprintf("%s_%s.pdf", taskID, time.Now().Format("20060102_150405"))

	// 生成完整的文件路径（绝对路径）
	filePath := filepath.Join(exportDir, fileName)

	// 确保目录存在
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("创建导出目录失败: %v", err)
	}

	// 保存PDF文件到文件系统
	if err := os.WriteFile(filePath, pdfData, 0644); err != nil {
		return "", fmt.Errorf("保存PDF文件失败: %v", err)
	}

	// 生成导出记录的唯一标识符
	exportKey := fmt.Sprintf("EXPORT%v%06d", time.Now().UnixNano(), len(records))

	// 创建导出记录（保存绝对路径）
	exportRecord := ExportRecord{
		Key:       exportKey,
		TaskID:    taskID,
		FilePath:  filePath, // 保存绝对路径
		CreatedAt: time.Now(),
	}

	// 保存导出记录
	if err := db.SaveExport(gctx.New(), exportToDB(exportRecord)); err != nil {
		// 如果保存记录失败，删除已创建的文件
		os.Remove(filePath)
		return "", fmt.Errorf("保存导出记录失败: %v", err)
	}

	return exportKey, nil
}

// ExportConsumptionsToExcel 将消费记录导出为 Excel 并保存到数据库
// 参数:
//   - taskID: 任务ID，如果为空则导出所有任务的消费记录
//   - keys: 要导出的消费记录Key列表，如果为空则导出所有记录
//
// 返回值:
//   - string: 导出记录的唯一标识符
//   - error: 如果导出失败返回错误
func (g *DbService) ExportConsumptionsToExcel(taskID string, keys []string) (string, error) {
	// 获取要导出的消费记录
	var records []ConsumptionRecord
	var err error

	if len(keys) > 0 {
		// 如果指定了 keys，获取对应的记录
		allRecords, err := g.ListConsumptions(taskID)
		if err != nil {
			return "", fmt.Errorf("获取消费记录失败: %v", err)
		}
		keyMap := make(map[string]bool)
		for _, key := range keys {
			keyMap[key] = true
		}
		for _, record := range allRecords {
			if keyMap[record.Key] {
				records = append(records, record)
			}
		}
	} else {
		// 否则获取所有记录
		records, err = g.ListConsumptions(taskID)
		if err != nil {
			return "", fmt.Errorf("获取消费记录失败: %v", err)
		}
	}

	if len(records) == 0 {
		return "", fmt.Errorf("没有可导出的消费记录")
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			gf.Log().Error(gctx.New(), "关闭 Excel 文件失败:", err)
		}
	}()

	// 设置工作表名称
	sheetName := "消费记录"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("创建工作表失败: %v", err)
	}
	f.SetActiveSheet(index)

	// 删除默认的 Sheet1
	f.DeleteSheet("Sheet1")

	// 设置表头
	headers := []string{"任务ID", "交易日期", "账号", "储种", "序号", "币种", "钞汇", "摘要", "地区", "收入/支出金额", "余额", "渠道"}
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F0F0F0"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建表头样式失败: %v", err)
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return "", fmt.Errorf("设置表头失败: %v", err)
		}
		if err := f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return "", fmt.Errorf("设置表头样式失败: %v", err)
		}
	}

	// 填充数据
	for rowIndex, record := range records {
		row := rowIndex + 2 // 从第2行开始（第1行是表头）
		rowData := []interface{}{
			record.TaskID,
			record.TradeDate,
			record.Account,
			record.StorageType,
			record.SerialNumber,
			record.Currency,
			record.CashOrRemit,
			record.Summary,
			record.Region,
			record.IncomeOrExpenseAmount,
			record.Balance,
			record.Channel,
		}

		for colIndex, value := range rowData {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, row)
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return "", fmt.Errorf("设置单元格值失败: %v", err)
			}
		}
	}

	// 设置列宽
	colWidths := []float64{15, 20, 15, 10, 8, 10, 8, 15, 10, 15, 12, 15}
	for i, width := range colWidths {
		col := fmt.Sprintf("%c", 'A'+i)
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return "", fmt.Errorf("设置列宽失败: %v", err)
		}
	}

	// 获取导出路径
	exportDir, err := g.GetExportPath()
	if err != nil {
		return "", fmt.Errorf("获取导出路径失败: %v", err)
	}

	// 生成文件名（使用任务ID作为标题）
	fileName := fmt.Sprintf("%s_%s.xlsx", taskID, time.Now().Format("20060102_150405"))

	// 生成完整的文件路径（绝对路径）
	filePath := filepath.Join(exportDir, fileName)

	// 确保目录存在
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("创建导出目录失败: %v", err)
	}

	// 保存 Excel 文件到文件系统
	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("保存Excel文件失败: %v", err)
	}

	// 生成导出记录的唯一标识符
	exportKey := fmt.Sprintf("EXPORT%v%06d", time.Now().UnixNano(), len(records))

	// 创建导出记录（保存绝对路径）
	exportRecord := ExportRecord{
		Key:       exportKey,
		TaskID:    taskID,
		FilePath:  filePath, // 保存绝对路径
		CreatedAt: time.Now(),
	}

	// 保存导出记录
	if err := db.SaveExport(gctx.New(), exportToDB(exportRecord)); err != nil {
		// 如果保存记录失败，删除已创建的文件
		os.Remove(filePath)
		return "", fmt.Errorf("保存导出记录失败: %v", err)
	}

	return exportKey, nil
}

// ListExports 获取所有导出记录
// 返回值:
//   - []ExportRecord: 导出记录列表，按创建时间倒序排列（最新的在前）
//   - error: 如果获取失败返回错误
func (g *DbService) ListExports() ([]ExportRecord, error) {
	raw, err := db.ListExports(gctx.New())
	if err != nil {
		return nil, err
	}

	records := make([]ExportRecord, 0, len(raw))
	for _, entry := range raw {
		records = append(records, exportFromDB(entry))
	}

	return records, nil
}

// SearchExports 根据关键词搜索导出记录
// 参数:
//   - keyword: 搜索关键词
//
// 返回值:
//   - []ExportRecord: 匹配的导出记录列表
//   - error: 如果获取失败返回错误
func (g *DbService) SearchExports(keyword string) ([]ExportRecord, error) {
	raw, err := db.SearchExports(gctx.New(), keyword)
	if err != nil {
		return nil, err
	}
	records := make([]ExportRecord, 0, len(raw))
	for _, entry := range raw {
		records = append(records, exportFromDB(entry))
	}
	return records, nil
}

// DeleteExports 批量删除导出记录
// 参数:
//   - keys: 要删除的导出记录Key列表
//
// 返回值:
//   - string: 成功返回"success"
//   - error: 如果删除失败返回错误
func (g *DbService) DeleteExports(keys []string) (string, error) {
	for _, key := range keys {
		// 获取导出记录，删除关联的PDF文件
		exportData, err := db.GetExport(gctx.New(), key)
		if err == nil {
			// 删除PDF文件（从文件系统删除）
			if exportData.FilePath != "" {
				os.Remove(exportData.FilePath)
			}
		}
		// 删除导出记录
		if err := db.DeleteExports(gctx.New(), []string{key}); err != nil {
			return "", fmt.Errorf("删除导出记录失败: %v", err)
		}
	}
	return "success", nil
}

// GetDefaultExportPath 获取默认导出路径
// 返回值:
//   - string: 默认导出路径（系统下载目录下的 bill_help/export/）
//   - error: 如果获取失败返回错误
func (g *DbService) GetDefaultExportPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户主目录失败: %v", err)
	}

	// 根据操作系统确定下载目录
	var downloadDir string
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") || os.Getenv("OS") == "":
		// Windows
		downloadDir = filepath.Join(homeDir, "Downloads")
	case os.Getenv("HOME") != "":
		// Linux/Mac
		downloadDir = filepath.Join(homeDir, "Downloads")
	default:
		downloadDir = homeDir
	}

	// 创建导出目录
	exportDir := filepath.Join(downloadDir, "bill_help", "export")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("创建导出目录失败: %v", err)
	}

	// 返回绝对路径
	absPath, err := filepath.Abs(exportDir)
	if err != nil {
		return exportDir, nil // 如果获取绝对路径失败，返回相对路径
	}

	return absPath, nil
}

// SelectExportPath 打开文件夹选择对话框
// 参数:
//   - defaultPath: 默认路径，如果为空则使用系统默认导出路径
//
// 返回值:
//   - string: 用户选择的文件夹路径（绝对路径）
//   - error: 如果选择失败返回错误
func (g *DbService) SelectExportPath(defaultPath string) (string, error) {
	app := application.Get()
	if app == nil {
		return "", fmt.Errorf("应用尚未初始化")
	}

	// 如果没有提供默认路径，使用系统默认导出路径
	if defaultPath == "" {
		var err error
		defaultPath, err = g.GetDefaultExportPath()
		if err != nil {
			defaultPath = ""
		}
	}

	// 使用 Wails 原生对话框，并绑定到当前窗口。
	// 这样 macOS 上会以主窗口为父窗口显示，避免选择器被 Finder 或其他窗口遮挡；
	// Windows/Linux 也复用同一套原生实现，保持跨平台行为一致。
	parentWindow := app.Window.Current()
	if parentWindow == nil {
		return "", fmt.Errorf("主窗口尚未准备完成")
	}

	dialog := app.Dialog.OpenFile().
		CanChooseFiles(false).
		CanChooseDirectories(true).
		CanCreateDirectories(true).
		SetTitle("选择导出路径").
		SetMessage("请选择导出文件夹").
		SetDirectory(defaultPath)
	dialog.AttachToWindow(parentWindow)

	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("选择文件夹失败: %v", err)
	}
	result = strings.TrimSpace(result)

	if result == "" {
		return "", fmt.Errorf("未选择文件夹")
	}

	// 验证路径是否存在，如果不存在则创建
	if info, err := os.Stat(result); os.IsNotExist(err) {
		// 路径不存在，尝试创建
		if err := os.MkdirAll(result, 0755); err != nil {
			return "", fmt.Errorf("创建目录失败: %v", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("路径无效: %v", err)
	} else if !info.IsDir() {
		return "", fmt.Errorf("路径不是文件夹: %s", result)
	}

	// 返回绝对路径
	absPath, err := filepath.Abs(result)
	if err != nil {
		return result, nil // 如果获取绝对路径失败，返回原始路径
	}

	return absPath, nil
}

// OpenFolder 打开指定文件夹，如果是文件路径则选中该文件
// 参数:
//   - path: 文件夹路径或文件路径
//
// 返回值:
//   - string: 成功返回"success"
//   - error: 如果打开失败返回错误
func (g *DbService) OpenFolder(path string) (string, error) {
	// 检查路径是否存在
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("路径不存在: %v", err)
	}

	// 使用系统命令打开文件夹，如果是文件则选中该文件
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		if info.IsDir() {
			// Windows: 如果是文件夹，直接打开
			cmd = exec.Command("explorer", path)
		} else {
			// Windows: 如果是文件，打开文件夹并选中该文件
			// 使用 /select, 参数来选中文件
			cmd = exec.Command("explorer", "/select,", path)
		}
	case "darwin":
		if info.IsDir() {
			// macOS: 如果是文件夹，直接打开
			cmd = exec.Command("open", path)
		} else {
			// macOS: 如果是文件，打开文件夹并选中该文件
			// 使用 -R 参数来显示文件并选中
			cmd = exec.Command("open", "-R", path)
		}
	default:
		// Linux: 使用 xdg-open 打开文件夹
		// 如果是文件，获取其目录
		dirPath := path
		if !info.IsDir() {
			dirPath = filepath.Dir(path)
		}
		cmd = exec.Command("xdg-open", dirPath)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("打开文件夹失败: %v", err)
	}

	return "success", nil
}

// GetExportPath 获取当前设置的导出路径
// 返回值:
//   - string: 导出路径（绝对路径）
//   - error: 如果获取失败返回错误
func (g *DbService) GetExportPath() (string, error) {
	// 从系统参数中获取导出路径
	params, err := g.GetSystemParameters()
	if err == nil && params.ExportPath != "" {
		// 返回绝对路径
		absPath, err := filepath.Abs(params.ExportPath)
		if err != nil {
			return params.ExportPath, nil // 如果获取绝对路径失败，返回原始路径
		}
		return absPath, nil
	}

	// 如果系统参数中没有设置，返回默认路径
	return g.GetDefaultExportPath()
}
