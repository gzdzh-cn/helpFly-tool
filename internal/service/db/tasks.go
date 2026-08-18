package db

import (
	"context"
	"fmt"

	"helpfly/internal/dao"
	"helpfly/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
)

type taskRow struct {
	TaskKey              string   `orm:"task_key"`
	TaskID               string   `orm:"task_id"`
	TaskName             string   `orm:"task_name"`
	AccountName          string   `orm:"account_name"`
	Account              string   `orm:"account"`
	CardNumber           string   `orm:"card_number"`
	Bank                 string   `orm:"bank"`
	DateRangeStart       string   `orm:"date_range_start"`
	DateRangeEnd         string   `orm:"date_range_end"`
	DepositType          string   `orm:"deposit_type"`
	Currency             string   `orm:"currency"`
	TransactionDateStart string   `orm:"transaction_date_start"`
	TransactionDateEnd   string   `orm:"transaction_date_end"`
	CashExchange         string   `orm:"cash_exchange"`
	OpeningBalance       *float64 `orm:"opening_balance"`
	AmountType           string   `orm:"amount_type"`
	OrderTime            string   `orm:"order_time"`
	TransactionCount     *int     `orm:"transaction_count"`
	Status               int      `orm:"status"`
	QRCodeURL            string   `orm:"qr_code_url"`
	StartSerialNumber    string   `orm:"start_serial_number"`
	CreatedAt            string   `orm:"created_at"`
}

type taskOptionRow struct {
	TaskID    string `orm:"task_id"`
	GroupName string `orm:"group_name"`
	Position  int    `orm:"position"`
	Value     string `orm:"value"`
}

func GetTask(ctx context.Context, key string) (Task, error) {
	if err := ready(); err != nil {
		return Task{}, err
	}
	var row taskRow
	err := dao.Tasks.Ctx(normalizeContext(ctx)).Where(dao.Tasks.Columns().TaskKey, key).Scan(&row)
	if isNoRows(err) {
		return Task{}, notFoundError("任务不存在")
	}
	if err != nil {
		return Task{}, fmt.Errorf("查询任务失败: %w", err)
	}
	options, err := loadTaskOptions(normalizeContext(ctx), []string{row.TaskID})
	if err != nil {
		return Task{}, err
	}
	return taskFromRow(row, options[row.TaskID]), nil
}

func TaskExists(ctx context.Context, taskID string) (bool, error) {
	if err := ready(); err != nil {
		return false, err
	}
	exists, err := dao.Tasks.Ctx(normalizeContext(ctx)).Where(dao.Tasks.Columns().TaskId, taskID).Exist()
	if err != nil {
		return false, fmt.Errorf("检查任务是否存在失败: %w", err)
	}
	return exists, nil
}

func SaveTask(ctx context.Context, task Task) error {
	if err := ready(); err != nil {
		return err
	}
	ctx = normalizeContext(ctx)
	if task.Key == "" {
		task.Key = task.TaskID
	}
	if task.TaskID == "" {
		task.TaskID = task.Key
	}
	return dao.Tasks.Transaction(ctx, func(_ context.Context, tx gdb.TX) error {
		data := do.Tasks{
			TaskKey:              task.Key,
			TaskId:               task.TaskID,
			TaskName:             task.TaskName,
			AccountName:          task.AccountName,
			Account:              task.Account,
			CardNumber:           task.CardNumber,
			Bank:                 task.Bank,
			DateRangeStart:       valueAt(task.DateRange, 0),
			DateRangeEnd:         valueAt(task.DateRange, 1),
			DepositType:          task.DepositType,
			Currency:             task.Currency,
			TransactionDateStart: valueAt(task.TransactionDateRange, 0),
			TransactionDateEnd:   valueAt(task.TransactionDateRange, 1),
			CashExchange:         task.CashExchange,
			OpeningBalance:       task.OpeningBalance,
			AmountType:           task.AmountType,
			OrderTime:            task.OrderTime,
			TransactionCount:     task.TransactionCount,
			Status:               task.Status,
			QrCodeUrl:            task.QRCodeURL,
			StartSerialNumber:    task.StartSerialNumber,
			CreatedAt:            formatDBTime(task.CreatedAt),
		}
		if _, err := tx.Model(dao.Tasks.Table()).Data(data).OnConflict(dao.Tasks.Columns().TaskKey).Save(); err != nil {
			return fmt.Errorf("保存任务失败: %w", err)
		}
		if _, err := tx.Model(dao.TaskChannels.Table()).Where(dao.TaskChannels.Columns().TaskId, task.TaskID).Delete(); err != nil {
			return fmt.Errorf("清理任务渠道失败: %w", err)
		}
		if err := insertTaskOptions(tx, dao.TaskChannels.Table(), dao.TaskChannels.Columns().TaskId, task.TaskID, task.Channels); err != nil {
			return err
		}
		if _, err := tx.Model(dao.TaskSummaries.Table()).Where(dao.TaskSummaries.Columns().TaskId, task.TaskID).Delete(); err != nil {
			return fmt.Errorf("清理任务摘要失败: %w", err)
		}
		return insertTaskOptions(tx, dao.TaskSummaries.Table(), dao.TaskSummaries.Columns().TaskId, task.TaskID, task.Summaries)
	})
}

func ListTasks(ctx context.Context) ([]Task, error) {
	return listTasks(ctx, nil)
}

func SearchTasks(ctx context.Context, keyword string) ([]Task, error) {
	return listTasks(ctx, &keyword)
}

func listTasks(ctx context.Context, keyword *string) ([]Task, error) {
	if err := ready(); err != nil {
		return nil, err
	}
	ctx = normalizeContext(ctx)
	model := dao.Tasks.Ctx(ctx).OrderDesc(dao.Tasks.Columns().CreatedAt)
	if keyword != nil {
		pattern := "%" + *keyword + "%"
		where := model.Builder().
			WhereLike(dao.Tasks.Columns().TaskId, pattern).
			WhereOrLike(dao.Tasks.Columns().TaskName, pattern).
			WhereOrLike(dao.Tasks.Columns().AccountName, pattern).
			WhereOrLike(dao.Tasks.Columns().Account, pattern).
			WhereOrLike(dao.Tasks.Columns().Bank, pattern).
			WhereOrLike(dao.Tasks.Columns().Currency, pattern)
		model = model.Where(where)
	}
	var rows []taskRow
	if err := model.Scan(&rows); err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %w", err)
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.TaskID)
	}
	options, err := loadTaskOptions(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make([]Task, 0, len(rows))
	for _, row := range rows {
		result = append(result, taskFromRow(row, options[row.TaskID]))
	}
	return result, nil
}

func UpdateTaskStatus(ctx context.Context, taskID string, status int) error {
	if err := ready(); err != nil {
		return err
	}
	_, err := dao.Tasks.Ctx(normalizeContext(ctx)).
		Data(do.Tasks{Status: status}).
		Where(dao.Tasks.Columns().TaskId, taskID).Update()
	if err != nil {
		return fmt.Errorf("更新任务状态失败: %w", err)
	}
	return nil
}

func DeleteTasks(ctx context.Context, keys []string) error {
	if err := ready(); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	_, err := dao.Tasks.Ctx(normalizeContext(ctx)).WhereIn(dao.Tasks.Columns().TaskKey, keys).Delete()
	if err != nil {
		return fmt.Errorf("删除任务失败: %w", err)
	}
	return nil
}

func loadTaskOptions(ctx context.Context, taskIDs []string) (map[string]taskOptions, error) {
	result := make(map[string]taskOptions, len(taskIDs))
	if len(taskIDs) == 0 {
		return result, nil
	}
	for _, id := range taskIDs {
		result[id] = taskOptions{}
	}
	var channels []taskOptionRow
	if err := dao.TaskChannels.Ctx(ctx).WhereIn(dao.TaskChannels.Columns().TaskId, taskIDs).
		OrderAsc(dao.TaskChannels.Columns().TaskId).OrderAsc(dao.TaskChannels.Columns().Position).Scan(&channels); err != nil {
		return nil, fmt.Errorf("查询任务渠道失败: %w", err)
	}
	for _, row := range channels {
		options := result[row.TaskID]
		options.Channels = append(options.Channels, row.Value)
		result[row.TaskID] = options
	}
	var summaries []taskOptionRow
	if err := dao.TaskSummaries.Ctx(ctx).WhereIn(dao.TaskSummaries.Columns().TaskId, taskIDs).
		OrderAsc(dao.TaskSummaries.Columns().TaskId).OrderAsc(dao.TaskSummaries.Columns().Position).Scan(&summaries); err != nil {
		return nil, fmt.Errorf("查询任务摘要失败: %w", err)
	}
	for _, row := range summaries {
		options := result[row.TaskID]
		options.Summaries = append(options.Summaries, row.Value)
		result[row.TaskID] = options
	}
	return result, nil
}

type taskOptions struct {
	Channels  []string
	Summaries []string
}

func insertTaskOptions(tx gdb.TX, table, column, taskID string, values []string) error {
	if len(values) == 0 {
		return nil
	}
	data := make(gdb.List, 0, len(values))
	for position, value := range values {
		data = append(data, gdb.Map{column: taskID, "position": position, "value": value})
	}
	if _, err := tx.Model(table).Data(data).Batch(len(data)).Insert(); err != nil {
		return fmt.Errorf("批量保存任务选项失败: %w", err)
	}
	return nil
}

func taskFromRow(row taskRow, options taskOptions) Task {
	return Task{
		Key:                  row.TaskKey,
		TaskID:               row.TaskID,
		TaskName:             row.TaskName,
		AccountName:          row.AccountName,
		Account:              row.Account,
		CardNumber:           row.CardNumber,
		Bank:                 row.Bank,
		DateRange:            []string{row.DateRangeStart, row.DateRangeEnd},
		DepositType:          row.DepositType,
		Currency:             row.Currency,
		TransactionDateRange: []string{row.TransactionDateStart, row.TransactionDateEnd},
		CashExchange:         row.CashExchange,
		Channels:             options.Channels,
		Summaries:            options.Summaries,
		OpeningBalance:       row.OpeningBalance,
		AmountType:           row.AmountType,
		OrderTime:            row.OrderTime,
		TransactionCount:     row.TransactionCount,
		Status:               row.Status,
		QRCodeURL:            row.QRCodeURL,
		StartSerialNumber:    row.StartSerialNumber,
		CreatedAt:            parseDBTime(row.CreatedAt),
	}
}

func valueAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}
