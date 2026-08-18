package db

import (
	"context"
	"fmt"

	"helpfly/internal/dao"
	"helpfly/internal/model/do"
)

type pausePointRow struct {
	TaskID           string  `orm:"task_id"`
	LastSerialNumber int     `orm:"last_serial_number"`
	CurrentBalance   float64 `orm:"current_balance"`
	CurrentProgress  int     `orm:"current_progress"`
	Percent          float64 `orm:"percent"`
	PausedAt         string  `orm:"paused_at"`
}

func GetPausePoint(ctx context.Context, taskID string) (PausePoint, error) {
	if err := ready(); err != nil {
		return PausePoint{}, err
	}
	var row pausePointRow
	err := dao.PausePoints.Ctx(normalizeContext(ctx)).Where(dao.PausePoints.Columns().TaskId, taskID).Scan(&row)
	if isNoRows(err) {
		return PausePoint{}, notFoundError("暂停点不存在")
	}
	if err != nil {
		return PausePoint{}, fmt.Errorf("查询暂停点失败: %w", err)
	}
	return PausePoint{
		TaskID:           row.TaskID,
		LastSerialNumber: row.LastSerialNumber,
		CurrentBalance:   row.CurrentBalance,
		CurrentProgress:  row.CurrentProgress,
		Percent:          row.Percent,
		PausedAt:         row.PausedAt,
	}, nil
}

func SavePausePoint(ctx context.Context, point PausePoint) error {
	if err := ready(); err != nil {
		return err
	}
	_, err := dao.PausePoints.Ctx(normalizeContext(ctx)).Data(do.PausePoints{
		TaskId:           point.TaskID,
		LastSerialNumber: point.LastSerialNumber,
		CurrentBalance:   point.CurrentBalance,
		CurrentProgress:  point.CurrentProgress,
		Percent:          point.Percent,
		PausedAt:         point.PausedAt,
	}).OnConflict(dao.PausePoints.Columns().TaskId).Save()
	if err != nil {
		return fmt.Errorf("保存暂停点失败: %w", err)
	}
	return nil
}

func DeletePausePoint(ctx context.Context, taskID string) error {
	if err := ready(); err != nil {
		return err
	}
	if _, err := dao.PausePoints.Ctx(normalizeContext(ctx)).Where(dao.PausePoints.Columns().TaskId, taskID).Delete(); err != nil {
		return fmt.Errorf("删除暂停点失败: %w", err)
	}
	return nil
}
