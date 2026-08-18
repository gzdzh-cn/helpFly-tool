// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SystemParametersDao is the data access object for the table system_parameters.
type SystemParametersDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  SystemParametersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// SystemParametersColumns defines and stores column names for the table system_parameters.
type SystemParametersColumns struct {
	Id            string //
	ExportPath    string //
	AddWatermark  string //
	WatermarkPath string //
}

// systemParametersColumns holds the columns for the table system_parameters.
var systemParametersColumns = SystemParametersColumns{
	Id:            "id",
	ExportPath:    "export_path",
	AddWatermark:  "add_watermark",
	WatermarkPath: "watermark_path",
}

// NewSystemParametersDao creates and returns a new DAO object for table data access.
func NewSystemParametersDao(handlers ...gdb.ModelHandler) *SystemParametersDao {
	return &SystemParametersDao{
		group:    "default",
		table:    "system_parameters",
		columns:  systemParametersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SystemParametersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SystemParametersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SystemParametersDao) Columns() SystemParametersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SystemParametersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SystemParametersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SystemParametersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
