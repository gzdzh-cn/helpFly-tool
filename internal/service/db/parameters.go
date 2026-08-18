package db

import (
	"context"
	"fmt"

	"helpfly/internal/dao"
	"helpfly/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
)

type parametersRow struct {
	ID            int    `orm:"id"`
	ExportPath    string `orm:"export_path"`
	AddWatermark  int    `orm:"add_watermark"`
	WatermarkPath string `orm:"watermark_path"`
}

type parameterOptionRow struct {
	GroupName string `orm:"group_name"`
	Position  int    `orm:"position"`
	Value     string `orm:"value"`
}

func GetParameters(ctx context.Context) (SystemParameters, error) {
	if err := ready(); err != nil {
		return SystemParameters{}, err
	}
	ctx = normalizeContext(ctx)
	var row parametersRow
	if err := dao.SystemParameters.Ctx(ctx).Where(dao.SystemParameters.Columns().Id, 1).Scan(&row); err != nil {
		if isNoRows(err) {
			return SystemParameters{}, notFoundError("系统参数不存在")
		}
		return SystemParameters{}, fmt.Errorf("查询系统参数失败: %w", err)
	}
	var options []parameterOptionRow
	if err := dao.SystemParameterOptions.Ctx(ctx).
		OrderAsc(dao.SystemParameterOptions.Columns().GroupName).
		OrderAsc(dao.SystemParameterOptions.Columns().Position).Scan(&options); err != nil {
		return SystemParameters{}, fmt.Errorf("查询系统参数选项失败: %w", err)
	}
	result := SystemParameters{
		ExportPath:    row.ExportPath,
		AddWatermark:  row.AddWatermark != 0,
		WatermarkPath: row.WatermarkPath,
	}
	for _, option := range options {
		switch option.GroupName {
		case "banks":
			result.Banks = append(result.Banks, option.Value)
		case "depositTypes":
			result.DepositTypes = append(result.DepositTypes, option.Value)
		case "cashExchanges":
			result.CashExchanges = append(result.CashExchanges, option.Value)
		case "currencies":
			result.Currencies = append(result.Currencies, option.Value)
		case "channels":
			result.Channels = append(result.Channels, option.Value)
		case "summaries":
			result.Summaries = append(result.Summaries, option.Value)
		case "regions":
			result.Regions = append(result.Regions, option.Value)
		}
	}
	return result, nil
}

func SaveParameters(ctx context.Context, params SystemParameters) error {
	if err := ready(); err != nil {
		return err
	}
	ctx = normalizeContext(ctx)
	return dao.SystemParameters.Transaction(ctx, func(_ context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemParameters.Table()).Data(do.SystemParameters{
			Id:            1,
			ExportPath:    params.ExportPath,
			AddWatermark:  boolToInt(params.AddWatermark),
			WatermarkPath: params.WatermarkPath,
		}).OnConflict(dao.SystemParameters.Columns().Id).Save(); err != nil {
			return fmt.Errorf("保存系统参数失败: %w", err)
		}
		if _, err := tx.Model(dao.SystemParameterOptions.Table()).
			WhereNot(dao.SystemParameterOptions.Columns().GroupName, "").Delete(); err != nil {
			return fmt.Errorf("清理系统参数选项失败: %w", err)
		}
		data := make(gdb.List, 0)
		for _, group := range parameterOptions(params) {
			for position, value := range group.values {
				data = append(data, gdb.Map{
					dao.SystemParameterOptions.Columns().GroupName: group.group,
					dao.SystemParameterOptions.Columns().Position:  position,
					dao.SystemParameterOptions.Columns().Value:     value,
				})
			}
		}
		if len(data) > 0 {
			if _, err := tx.Model(dao.SystemParameterOptions.Table()).Data(data).Batch(len(data)).Insert(); err != nil {
				return fmt.Errorf("保存系统参数选项失败: %w", err)
			}
		}
		return nil
	})
}

func parameterOptions(params SystemParameters) []struct {
	group  string
	values []string
} {
	return []struct {
		group  string
		values []string
	}{
		{"banks", params.Banks},
		{"depositTypes", params.DepositTypes},
		{"cashExchanges", params.CashExchanges},
		{"currencies", params.Currencies},
		{"channels", params.Channels},
		{"summaries", params.Summaries},
		{"regions", params.Regions},
	}
}
