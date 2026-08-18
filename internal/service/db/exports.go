package db

import (
	"context"
	"fmt"
	"strings"

	"helpfly/internal/dao"
	"helpfly/internal/model/do"
)

type exportRow struct {
	RecordKey string `orm:"record_key"`
	TaskID    string `orm:"task_id"`
	FilePath  string `orm:"file_path"`
	CreatedAt string `orm:"created_at"`
}

func SaveExport(ctx context.Context, record Export) error {
	if err := ready(); err != nil {
		return err
	}
	_, err := dao.Exports.Ctx(normalizeContext(ctx)).Data(do.Exports{
		RecordKey: record.Key,
		TaskId:    record.TaskID,
		FilePath:  record.FilePath,
		CreatedAt: formatDBTime(record.CreatedAt),
	}).OnConflict(dao.Exports.Columns().RecordKey).Save()
	if err != nil {
		return fmt.Errorf("保存导出记录失败: %w", err)
	}
	return nil
}

func GetExport(ctx context.Context, key string) (Export, error) {
	if err := ready(); err != nil {
		return Export{}, err
	}
	var row exportRow
	err := dao.Exports.Ctx(normalizeContext(ctx)).Where(dao.Exports.Columns().RecordKey, key).Scan(&row)
	if isNoRows(err) {
		return Export{}, notFoundError("导出记录不存在")
	}
	if err != nil {
		return Export{}, fmt.Errorf("查询导出记录失败: %w", err)
	}
	return exportFromRow(row), nil
}

func ListExports(ctx context.Context) ([]Export, error) {
	if err := ready(); err != nil {
		return nil, err
	}
	var rows []exportRow
	if err := dao.Exports.Ctx(normalizeContext(ctx)).OrderDesc(dao.Exports.Columns().CreatedAt).Scan(&rows); err != nil {
		return nil, fmt.Errorf("查询导出记录失败: %w", err)
	}
	return exportRowsFromRows(rows), nil
}

func SearchExports(ctx context.Context, keyword string) ([]Export, error) {
	if err := ready(); err != nil {
		return nil, err
	}
	pattern := "%" + strings.TrimSpace(keyword) + "%"
	model := dao.Exports.Ctx(normalizeContext(ctx)).OrderDesc(dao.Exports.Columns().CreatedAt)
	model = model.Where(model.Builder().
		WhereLike(dao.Exports.Columns().TaskId, pattern).
		WhereOrLike(dao.Exports.Columns().FilePath, pattern))
	var rows []exportRow
	if err := model.Scan(&rows); err != nil {
		return nil, fmt.Errorf("搜索导出记录失败: %w", err)
	}
	return exportRowsFromRows(rows), nil
}

func DeleteExports(ctx context.Context, keys []string) error {
	if err := ready(); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	if _, err := dao.Exports.Ctx(normalizeContext(ctx)).WhereIn(dao.Exports.Columns().RecordKey, keys).Delete(); err != nil {
		return fmt.Errorf("删除导出记录失败: %w", err)
	}
	return nil
}

func exportRowsFromRows(rows []exportRow) []Export {
	result := make([]Export, 0, len(rows))
	for _, row := range rows {
		result = append(result, exportFromRow(row))
	}
	return result
}

func exportFromRow(row exportRow) Export {
	return Export{Key: row.RecordKey, TaskID: row.TaskID, FilePath: row.FilePath, CreatedAt: parseDBTime(row.CreatedAt)}
}
