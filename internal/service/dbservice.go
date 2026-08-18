// Package service 提供数据库服务接口，包括任务管理、消费记录生成等功能
package service

import (
	"fmt"
	"helpfly/internal/service/db"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
)

// DbService 数据库服务结构体，提供基础的数据库操作接口
// 该服务通过 Wails 框架暴露给前端调用
type DbService struct{}

// GetData 从数据库中获取测试数据
// 返回值: 如果成功返回数据字符串，失败返回错误信息
func (g *DbService) GetData() any {
	data, err := db.GetKV(gctx.New(), "testkey")
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// SetData 保存数据到数据库
// 参数:
//   - param: 要保存的字符串数据
//
// 返回值: 保存成功后返回保存的数据，失败返回错误
func (g *DbService) SetData(param string) any {
	err := db.SetKV(gctx.New(), "testkey", []byte(param))
	if err != nil {
		return err
	}
	data, err := db.GetKV(gctx.New(), "testkey")
	if err != nil {
		return err
	}
	return string(data)
}

// AddConsumption 新增一条消费记录（用于前端账单表格打通后端演示）
func (g *DbService) AddConsumption(record db.Consumption) any {
	if record.Key == "" {
		record.Key = fmt.Sprintf("C%v", time.Now().UnixNano())
	}
	if err := db.SaveConsumption(gctx.New(), record); err != nil {
		return err
	}
	return record
}

// UpdateConsumption 根据 Key 更新一条消费记录（用于前端账单表格打通后端演示）
func (g *DbService) UpdateConsumption(record db.Consumption) any {
	if err := db.UpdateConsumption(gctx.New(), record); err != nil {
		return err
	}
	return record
}
