// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PausePointsDao is the data access object for the table pause_points.
type PausePointsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PausePointsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PausePointsColumns defines and stores column names for the table pause_points.
type PausePointsColumns struct {
	TaskId           string //
	LastSerialNumber string //
	CurrentBalance   string //
	CurrentProgress  string //
	Percent          string //
	PausedAt         string //
}

// pausePointsColumns holds the columns for the table pause_points.
var pausePointsColumns = PausePointsColumns{
	TaskId:           "task_id",
	LastSerialNumber: "last_serial_number",
	CurrentBalance:   "current_balance",
	CurrentProgress:  "current_progress",
	Percent:          "percent",
	PausedAt:         "paused_at",
}

// NewPausePointsDao creates and returns a new DAO object for table data access.
func NewPausePointsDao(handlers ...gdb.ModelHandler) *PausePointsDao {
	return &PausePointsDao{
		group:    "default",
		table:    "pause_points",
		columns:  pausePointsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PausePointsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PausePointsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PausePointsDao) Columns() PausePointsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PausePointsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PausePointsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *PausePointsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
