// Package service 提供参数设置服务接口
package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"helpfly/internal/service/db"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SystemParameters 系统参数结构体
// 用于存储系统配置参数，包括各种选项列表
type SystemParameters struct {
	Banks         []string `json:"banks"`         // 银行选项列表
	DepositTypes  []string `json:"depositTypes"`  // 储种选项列表
	CashExchanges []string `json:"cashExchanges"` // 钞汇选项列表
	Currencies    []string `json:"currencies"`    // 币种选项列表
	Channels      []string `json:"channels"`      // 渠道选项列表
	Summaries     []string `json:"summaries"`     // 摘要选项列表
	Regions       []string `json:"regions"`       // 地区选项列表
	ExportPath    string   `json:"exportPath"`    // PDF导出路径（绝对路径）
	AddWatermark  bool     `json:"addWatermark"`  // 导出PDF时是否添加水印
	WatermarkPath string   `json:"watermarkPath"` // PDF水印图片路径（PNG）
}

// GetSystemParameters 获取系统参数
// 返回值:
//   - SystemParameters: 系统参数
//   - error: 如果获取失败返回错误
func (g *DbService) GetSystemParameters() (SystemParameters, error) {
	// 尝试从数据库获取参数
	data, err := db.GetParameters(gctx.New())
	if err == nil {
		return systemParametersFromDB(data), nil
	}
	if gerror.Code(err) != gcode.CodeNotFound {
		return SystemParameters{}, err
	}

	// 如果数据库中没有参数或解析失败，返回空参数
	return SystemParameters{}, nil
}

// SaveSystemParameters 保存系统参数
// 参数:
//   - params: 系统参数
//
// 返回值:
//   - string: 成功返回"保存成功"
//   - error: 如果保存失败返回错误
func (g *DbService) SaveSystemParameters(params SystemParameters) (string, error) {
	// 验证参数
	for _, field := range []struct {
		value   []string
		message string
	}{
		{params.Banks, "银行选项不能为空"}, {params.DepositTypes, "储种选项不能为空"},
		{params.CashExchanges, "钞汇选项不能为空"}, {params.Currencies, "币种选项不能为空"},
		{params.Channels, "渠道选项不能为空"}, {params.Summaries, "摘要选项不能为空"},
		{params.Regions, "地区选项不能为空"},
	} {
		if err := validateOptionList(field.value, field.message); err != nil {
			return "", err
		}
	}
	// 验证导出路径（如果为空，使用默认路径）
	if params.ExportPath == "" {
		defaultPath, err := g.GetDefaultExportPath()
		if err == nil {
			params.ExportPath = defaultPath
		}
	}
	// 水印图片为可选项；保留旧版文字水印的兼容逻辑。
	if params.WatermarkPath != "" {
		if strings.ToLower(filepath.Ext(params.WatermarkPath)) != ".png" {
			return "", fmt.Errorf("水印图片仅支持 PNG 格式")
		}
		if params.AddWatermark {
			info, err := os.Stat(params.WatermarkPath)
			if err != nil {
				return "", fmt.Errorf("水印图片不存在: %v", err)
			}
			if info.IsDir() {
				return "", fmt.Errorf("水印路径不是图片文件")
			}
		}
	}

	// 保存到数据库
	if err := db.SaveParameters(gctx.New(), systemParametersToDB(params)); err != nil {
		return "", fmt.Errorf("保存参数失败: %v", err)
	}

	return "保存成功", nil
}

func validateOptionList(values []string, message string) error {
	if err := g.Validator().Data(values).Rules("required|min-length:1").Messages(message).Run(gctx.New()); err != nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, message)
	}
	return nil
}

func systemParametersToDB(params SystemParameters) db.SystemParameters {
	return db.SystemParameters{
		Banks: params.Banks, DepositTypes: params.DepositTypes, CashExchanges: params.CashExchanges,
		Currencies: params.Currencies, Channels: params.Channels, Summaries: params.Summaries,
		Regions: params.Regions, ExportPath: params.ExportPath, AddWatermark: params.AddWatermark,
		WatermarkPath: params.WatermarkPath,
	}
}

func systemParametersFromDB(params db.SystemParameters) SystemParameters {
	return SystemParameters{
		Banks: params.Banks, DepositTypes: params.DepositTypes, CashExchanges: params.CashExchanges,
		Currencies: params.Currencies, Channels: params.Channels, Summaries: params.Summaries,
		Regions: params.Regions, ExportPath: params.ExportPath, AddWatermark: params.AddWatermark,
		WatermarkPath: params.WatermarkPath,
	}
}

// SelectWatermarkImage 打开 PNG 图片选择对话框。
func (g *DbService) SelectWatermarkImage() (string, error) {
	app := application.Get()
	if app == nil {
		return "", fmt.Errorf("应用尚未初始化")
	}
	parentWindow := app.Window.Current()
	if parentWindow == nil {
		return "", fmt.Errorf("主窗口尚未准备完成")
	}

	dialog := app.Dialog.OpenFile().
		CanChooseFiles(true).
		CanChooseDirectories(false).
		CanCreateDirectories(false).
		SetTitle("选择 PNG 水印图片").
		SetMessage("请选择 PNG 格式的水印图片").
		AddFilter("PNG 图片", "*.png;*.PNG").
		AttachToWindow(parentWindow)

	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("选择水印图片失败: %v", err)
	}
	result = strings.TrimSpace(result)
	if result == "" {
		return "", fmt.Errorf("未选择水印图片")
	}
	if strings.ToLower(filepath.Ext(result)) != ".png" {
		return "", fmt.Errorf("水印图片仅支持 PNG 格式")
	}

	info, err := os.Stat(result)
	if err != nil {
		return "", fmt.Errorf("水印图片不存在: %v", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("选择的路径不是图片文件")
	}

	return filepath.Abs(result)
}
