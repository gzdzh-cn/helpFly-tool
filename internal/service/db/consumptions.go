package db

import (
	"context"
	"fmt"
	"strings"

	"helpfly/internal/dao"
	"helpfly/internal/model/do"
)

type consumptionRow struct {
	RecordKey             string  `orm:"record_key"`
	TaskID                string  `orm:"task_id"`
	TradeDate             string  `orm:"trade_date"`
	Account               string  `orm:"account"`
	StorageType           string  `orm:"storage_type"`
	SerialNumber          string  `orm:"serial_number"`
	Currency              string  `orm:"currency"`
	CashOrRemit           string  `orm:"cash_or_remit"`
	Summary               string  `orm:"summary"`
	Region                string  `orm:"region"`
	IncomeOrExpenseAmount float64 `orm:"income_or_expense_amount"`
	Balance               float64 `orm:"balance"`
	Channel               string  `orm:"channel"`
}

func SaveConsumption(ctx context.Context, record Consumption) error {
	if err := ready(); err != nil {
		return err
	}
	_, err := dao.Consumptions.Ctx(normalizeContext(ctx)).Data(do.Consumptions{
		RecordKey:             record.Key,
		TaskId:                record.TaskID,
		TradeDate:             record.TradeDate,
		Account:               record.Account,
		StorageType:           record.StorageType,
		SerialNumber:          record.SerialNumber,
		Currency:              record.Currency,
		CashOrRemit:           record.CashOrRemit,
		Summary:               record.Summary,
		Region:                record.Region,
		IncomeOrExpenseAmount: record.IncomeOrExpenseAmount,
		Balance:               record.Balance,
		Channel:               record.Channel,
	}).Insert()
	if err != nil {
		return fmt.Errorf("保存消费记录失败: %w", err)
	}
	return nil
}

func GetConsumption(ctx context.Context, key string) (Consumption, error) {
	if err := ready(); err != nil {
		return Consumption{}, err
	}
	var row consumptionRow
	err := dao.Consumptions.Ctx(normalizeContext(ctx)).Where(dao.Consumptions.Columns().RecordKey, key).Scan(&row)
	if isNoRows(err) {
		return Consumption{}, notFoundError("消费记录不存在")
	}
	if err != nil {
		return Consumption{}, fmt.Errorf("查询消费记录失败: %w", err)
	}
	return consumptionFromRow(row), nil
}

// UpdateConsumption 根据记录 Key 更新消费记录的非主键字段（用于编辑）
func UpdateConsumption(ctx context.Context, record Consumption) error {
	if err := ready(); err != nil {
		return err
	}
	_, err := dao.Consumptions.Ctx(normalizeContext(ctx)).
		Data(do.Consumptions{
			TradeDate:             record.TradeDate,
			Account:               record.Account,
			StorageType:           record.StorageType,
			SerialNumber:          record.SerialNumber,
			Currency:              record.Currency,
			CashOrRemit:           record.CashOrRemit,
			Summary:               record.Summary,
			Region:                record.Region,
			IncomeOrExpenseAmount: record.IncomeOrExpenseAmount,
			Balance:               record.Balance,
			Channel:               record.Channel,
	 }).
		Where(dao.Consumptions.Columns().RecordKey, record.Key).
		Update()
	if err != nil {
		return fmt.Errorf("更新消费记录失败: %w", err)
	}
	return nil
}

func ListConsumptions(ctx context.Context, taskID string) ([]Consumption, error) {
	if err := ready(); err != nil {
		return nil, err
	}
	model := dao.Consumptions.Ctx(normalizeContext(ctx)).OrderAsc(dao.Consumptions.Columns().SerialNumber)
	if strings.TrimSpace(taskID) != "" {
		model = model.Where(dao.Consumptions.Columns().TaskId, taskID)
	}
	var rows []consumptionRow
	if err := model.Scan(&rows); err != nil {
		return nil, fmt.Errorf("查询消费记录失败: %w", err)
	}
	return consumptionRowsFromRows(rows), nil
}

func SearchConsumptions(ctx context.Context, taskID, keyword string) ([]Consumption, error) {
	if err := ready(); err != nil {
		return nil, err
	}
	pattern := "%" + strings.TrimSpace(keyword) + "%"
	model := dao.Consumptions.Ctx(normalizeContext(ctx)).OrderAsc(dao.Consumptions.Columns().SerialNumber)
	where := model.Builder().
		WhereLike(dao.Consumptions.Columns().TradeDate, pattern).
		WhereOrLike(dao.Consumptions.Columns().Account, pattern).
		WhereOrLike(dao.Consumptions.Columns().StorageType, pattern).
		WhereOrLike(dao.Consumptions.Columns().SerialNumber, pattern).
		WhereOrLike(dao.Consumptions.Columns().Currency, pattern).
		WhereOrLike(dao.Consumptions.Columns().CashOrRemit, pattern).
		WhereOrLike(dao.Consumptions.Columns().Summary, pattern).
		WhereOrLike(dao.Consumptions.Columns().Region, pattern).
		WhereOr("CAST(income_or_expense_amount AS TEXT) LIKE ?", pattern).
		WhereOr("CAST(balance AS TEXT) LIKE ?", pattern).
		WhereOrLike(dao.Consumptions.Columns().Channel, pattern)
	model = model.Where(where)
	if strings.TrimSpace(taskID) != "" {
		model = model.Where(dao.Consumptions.Columns().TaskId, taskID)
	}
	var rows []consumptionRow
	if err := model.Scan(&rows); err != nil {
		return nil, fmt.Errorf("搜索消费记录失败: %w", err)
	}
	return consumptionRowsFromRows(rows), nil
}

func DeleteConsumptions(ctx context.Context, keys []string) error {
	if err := ready(); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	if _, err := dao.Consumptions.Ctx(normalizeContext(ctx)).WhereIn(dao.Consumptions.Columns().RecordKey, keys).Delete(); err != nil {
		return fmt.Errorf("删除消费记录失败: %w", err)
	}
	return nil
}

func DeleteConsumptionsByTaskID(ctx context.Context, taskID string) error {
	if err := ready(); err != nil {
		return err
	}
	if _, err := dao.Consumptions.Ctx(normalizeContext(ctx)).Where(dao.Consumptions.Columns().TaskId, taskID).Delete(); err != nil {
		return fmt.Errorf("删除任务消费记录失败: %w", err)
	}
	return nil
}

func consumptionRowsFromRows(rows []consumptionRow) []Consumption {
	result := make([]Consumption, 0, len(rows))
	for _, row := range rows {
		result = append(result, consumptionFromRow(row))
	}
	return result
}

func consumptionFromRow(row consumptionRow) Consumption {
	return Consumption{
		Key:                   row.RecordKey,
		TaskID:                row.TaskID,
		TradeDate:             row.TradeDate,
		Account:               row.Account,
		StorageType:           row.StorageType,
		SerialNumber:          row.SerialNumber,
		Currency:              row.Currency,
		CashOrRemit:           row.CashOrRemit,
		Summary:               row.Summary,
		Region:                row.Region,
		IncomeOrExpenseAmount: row.IncomeOrExpenseAmount,
		Balance:               row.Balance,
		Channel:               row.Channel,
	}
}
