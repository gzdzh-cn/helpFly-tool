// Package service 提供数据库服务接口，包括任务管理、消费记录生成等功能
package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"helpfly/internal/service/db"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// TaskPayload 任务数据载荷结构体
// 包含任务的所有配置参数和业务数据
type TaskPayload struct {
	TaskName             string   `json:"taskName"`             // 任务名称
	AccountName          string   `json:"accountName"`          // 户名
	Account              string   `json:"account"`              // 账号（表格中“账号”列）
	CardNumber           string   `json:"cardNumber"`           // 卡号（PDF 顶部银行信息行使用）
	Bank                 string   `json:"bank"`                 // 银行名称
	DateRange            []string `json:"dateRange"`            // 起止日期范围 [开始日期, 结束日期]
	DepositType          string   `json:"depositType"`          // 储种（活期/定期）
	Currency             string   `json:"currency"`             // 币种（人民币/美元）
	TransactionDateRange []string `json:"transactionDateRange"` // 交易日期范围 [开始日期, 结束日期]
	CashExchange         string   `json:"cashExchange"`         // 钞汇类型（钞/汇）
	Channels             []string `json:"channels"`             // 渠道列表（多选）
	Summaries            []string `json:"summaries"`            // 摘要列表（多选）
	OpeningBalance       *float64 `json:"openingBalance"`       // 起始余额
	AmountType           string   `json:"amountType"`           // 收入/支出金额规则（规则一/规则二）
	OrderTime            string   `json:"orderTime"`            // 下单时间
	TransactionCount     *int     `json:"transactionCount"`     // 交易笔数（必须大于0）
	Status               int      `json:"status"`               // 任务状态：0=未开始, 1=进行中, 2=已完成, 3=已暂停, 4=已取消
	QRCodeURL            string   `json:"qrCodeURL"`            // 二维码 URL
	StartSerialNumber    string   `json:"startSerialNumber"`    // 序号起始值，用于控制消费列表序号（字符串）
}

// TaskRecord 任务记录结构体
// 包含任务的元数据（Key、TaskID、创建时间）和任务载荷数据
type TaskRecord struct {
	Key         string    `json:"key"`       // 数据库存储键，通常与TaskID相同
	TaskID      string    `json:"taskId"`    // 任务唯一标识符，格式：TASK + 6位数字
	CreatedAt   time.Time `json:"createdAt"` // 任务创建时间
	TaskPayload           // 嵌入任务载荷数据
}

// validateTaskPayload 校验任务载荷的必填字段
// 与前端 add.vue 表单中的必填项保持一致，CreateTask / UpdateTask 均复用
func validateTaskPayload(payload TaskPayload) error {
	for _, field := range []struct {
		value   string
		message string
	}{
		{payload.TaskName, "任务名称不能为空"}, {payload.AccountName, "户名不能为空"},
		{payload.Account, "账号不能为空"}, {payload.CardNumber, "卡号不能为空"},
		{payload.Bank, "银行不能为空"}, {payload.DepositType, "储种不能为空"},
		{payload.Currency, "币种不能为空"}, {payload.CashExchange, "钞汇类型不能为空"},
		{payload.AmountType, "收入/支出规则不能为空"}, {payload.OrderTime, "下单时间不能为空"},
		{payload.StartSerialNumber, "序号不能为空"}, {payload.QRCodeURL, "二维码URL不能为空"},
	} {
		if err := validateRequired(field.value, field.message); err != nil {
			return err
		}
	}
	if len(payload.DateRange) != 2 {
		return invalidParameter("起止日期不能为空")
	}
	if len(payload.TransactionDateRange) != 2 {
		return invalidParameter("交易日期不能为空")
	}
	if len(payload.Channels) == 0 {
		return invalidParameter("渠道不能为空")
	}
	if len(payload.Summaries) == 0 {
		return invalidParameter("摘要不能为空")
	}
	if payload.OpeningBalance == nil {
		return invalidParameter("起始余额不能为空")
	}
	if payload.TransactionCount == nil || *payload.TransactionCount <= 0 {
		return invalidParameter("交易笔数必须大于0")
	}
	return nil
}

func invalidParameter(message string) error {
	return gerror.NewCode(gcode.CodeInvalidParameter, message)
}

func validateRequired(value, message string) error {
	if err := g.Validator().Data(strings.TrimSpace(value)).Rules("required").Messages(message).Run(gctx.New()); err != nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, message)
	}
	return nil
}

// generateTaskID 生成唯一的任务ID
// 格式：TASK + 6位数字（如：TASK123456）
// 生成策略：
//  1. 使用当前时间戳的后6位作为基础数字
//  2. 如果冲突，递增数字重试（最多100次）
//  3. 如果仍冲突，使用随机数生成（最多100次）
//
// 返回值:
//   - string: 生成的任务ID
//   - error: 如果生成失败返回错误
func (g *DbService) generateTaskID() (string, error) {
	ctx := gctx.New()
	maxRetries := 100
	baseNum := int(time.Now().Unix() % 1000000) // 使用秒级时间戳的后6位作为基础

	for i := 0; i < maxRetries; i++ {
		// 生成6位数字，如果不足6位则前面补0
		num := (baseNum + i) % 1000000
		taskID := fmt.Sprintf("TASK%06d", num)

		// 检查ID是否已存在
		exists, err := db.TaskExists(ctx, taskID)
		if err != nil {
			return "", err
		}
		if !exists {
			return taskID, nil
		}
	}

	// 如果100次重试都失败，使用随机数
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < maxRetries; i++ {
		num := rng.Intn(1000000)
		taskID := fmt.Sprintf("TASK%06d", num)
		exists, err := db.TaskExists(ctx, taskID)
		if err != nil {
			return "", err
		}
		if !exists {
			return taskID, nil
		}
	}

	return "", fmt.Errorf("无法生成唯一的任务ID，请稍后重试")
}

// CreateTask 创建新任务并持久化到数据库
// 参数:
//   - payload: 任务数据载荷
//
// 返回值:
//   - TaskRecord: 创建的任务记录（包含生成的TaskID和Key）
//   - error: 如果创建失败返回错误
//
// 注意: 任务状态默认为0（未开始）
func (g *DbService) CreateTask(payload TaskPayload) (TaskRecord, error) {
	// 基本必填校验（与前端保持一致）
	if err := validateTaskPayload(payload); err != nil {
		return TaskRecord{}, err
	}

	taskID, err := g.generateTaskID()
	if err != nil {
		return TaskRecord{}, err
	}
	// status 默认为 0（Go int 类型的零值）
	record := TaskRecord{
		Key:         taskID,
		TaskID:      taskID,
		TaskPayload: payload,
		CreatedAt:   time.Now(),
	}

	if err := db.SaveTask(gctx.New(), taskToDB(record)); err != nil {
		return TaskRecord{}, err
	}
	return record, nil
}

// UpdateTask 更新已存在的任务记录
// 参数:
//   - record: 包含更新数据的任务记录（必须包含Key字段）
//
// 返回值:
//   - TaskRecord: 更新后的任务记录
//   - error: 如果更新失败返回错误
//
// 注意: 如果record中TaskID或CreatedAt为空，会从现有记录中继承
func (g *DbService) UpdateTask(record TaskRecord) (TaskRecord, error) {
	if record.Key == "" {
		return TaskRecord{}, errors.New("task key required")
	}

	existing, err := db.GetTask(gctx.New(), record.Key)
	if err == nil {
		if record.TaskID == "" {
			record.TaskID = existing.TaskID
		}
		if record.CreatedAt.IsZero() {
			record.CreatedAt = existing.CreatedAt
		}
	} else if gerror.Code(err) != gcode.CodeNotFound {
		return TaskRecord{}, err
	}
	if record.TaskID == "" {
		record.TaskID = record.Key
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}

	// 更新任务时同样校验必填字段
	if err := validateTaskPayload(record.TaskPayload); err != nil {
		return TaskRecord{}, err
	}

	if err := db.SaveTask(gctx.New(), taskToDB(record)); err != nil {
		return TaskRecord{}, err
	}
	return record, nil
}

// ListTasks 获取所有任务记录列表
// 返回值:
//   - []TaskRecord: 任务记录列表，按创建时间倒序排列（最新的在前）
//   - error: 如果获取失败返回错误
func (g *DbService) ListTasks() ([]TaskRecord, error) {
	raw, err := db.ListTasks(gctx.New())
	if err != nil {
		return nil, err
	}
	records := make([]TaskRecord, 0, len(raw))
	for _, entry := range raw {
		records = append(records, taskFromDB(entry))
	}
	return records, nil
}

// SearchTasks 根据关键词筛选任务记录
// 关键词会在 taskId、taskName、accountName、account、bank、currency 等字段中匹配
func (g *DbService) SearchTasks(keyword string) ([]TaskRecord, error) {
	trimmed := strings.TrimSpace(keyword)
	if trimmed == "" {
		return g.ListTasks()
	}

	raw, err := db.SearchTasks(gctx.New(), trimmed)
	if err != nil {
		return nil, err
	}
	records := make([]TaskRecord, 0, len(raw))
	for _, entry := range raw {
		records = append(records, taskFromDB(entry))
	}
	return records, nil
}

// DeleteTasks 批量删除任务
// 参数:
//   - keys: 要删除的任务Key列表
//
// 返回值:
//   - string: 成功返回"success"
//   - error: 如果删除失败返回错误
//
// 注意: 删除任务时会自动取消正在运行的任务并清理进度信息
func (g *DbService) DeleteTasks(keys []string) (string, error) {
	ctx := gctx.New()
	for _, key := range keys {
		if key == "" {
			continue
		}

		// Get task to find taskId
		task, err := db.GetTask(ctx, key)
		if err == nil {
			globalTaskManager.cancelTask(task.TaskID)
			globalTaskManager.removeTask(task.TaskID)
		} else if gerror.Code(err) != gcode.CodeNotFound {
			return "", err
		}
	}
	if err := db.DeleteTasks(ctx, keys); err != nil {
		return "", err
	}

	return "success", nil
}

func taskToDB(record TaskRecord) db.Task {
	return db.Task{
		Key: record.Key, TaskID: record.TaskID, TaskName: record.TaskName,
		AccountName: record.AccountName, Account: record.Account, CardNumber: record.CardNumber,
		Bank: record.Bank, DateRange: record.DateRange, DepositType: record.DepositType,
		Currency: record.Currency, TransactionDateRange: record.TransactionDateRange,
		CashExchange: record.CashExchange, Channels: record.Channels, Summaries: record.Summaries,
		OpeningBalance: record.OpeningBalance, AmountType: record.AmountType, OrderTime: record.OrderTime,
		TransactionCount: record.TransactionCount, Status: record.Status, QRCodeURL: record.QRCodeURL,
		StartSerialNumber: record.StartSerialNumber, CreatedAt: record.CreatedAt,
	}
}

func taskFromDB(record db.Task) TaskRecord {
	return TaskRecord{
		Key: record.Key, TaskID: record.TaskID, CreatedAt: record.CreatedAt,
		TaskPayload: TaskPayload{
			TaskName: record.TaskName, AccountName: record.AccountName, Account: record.Account,
			CardNumber: record.CardNumber, Bank: record.Bank, DateRange: record.DateRange,
			DepositType: record.DepositType, Currency: record.Currency,
			TransactionDateRange: record.TransactionDateRange, CashExchange: record.CashExchange,
			Channels: record.Channels, Summaries: record.Summaries, OpeningBalance: record.OpeningBalance,
			AmountType: record.AmountType, OrderTime: record.OrderTime, TransactionCount: record.TransactionCount,
			Status: record.Status, QRCodeURL: record.QRCodeURL, StartSerialNumber: record.StartSerialNumber,
		},
	}
}

type PromiseResult struct {
	Result string `json:"result"`
}

func (g *DbService) TestPromiose(key string) (*PromiseResult, error) {
	return nil, fmt.Errorf("任务不存在1")
}
