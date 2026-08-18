// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ConsumptionsDao is the data access object for the table consumptions.
type ConsumptionsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ConsumptionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ConsumptionsColumns defines and stores column names for the table consumptions.
type ConsumptionsColumns struct {
	RecordKey             string //
	TaskId                string //
	TradeDate             string //
	Account               string //
	StorageType           string //
	SerialNumber          string //
	Currency              string //
	CashOrRemit           string //
	Summary               string //
	Region                string //
	IncomeOrExpenseAmount string //
	Balance               string //
	Channel               string //
}

// consumptionsColumns holds the columns for the table consumptions.
var consumptionsColumns = ConsumptionsColumns{
	RecordKey:             "record_key",
	TaskId:                "task_id",
	TradeDate:             "trade_date",
	Account:               "account",
	StorageType:           "storage_type",
	SerialNumber:          "serial_number",
	Currency:              "currency",
	CashOrRemit:           "cash_or_remit",
	Summary:               "summary",
	Region:                "region",
	IncomeOrExpenseAmount: "income_or_expense_amount",
	Balance:               "balance",
	Channel:               "channel",
}

// NewConsumptionsDao creates and returns a new DAO object for table data access.
func NewConsumptionsDao(handlers ...gdb.ModelHandler) *ConsumptionsDao {
	return &ConsumptionsDao{
		group:    "default",
		table:    "consumptions",
		columns:  consumptionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ConsumptionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ConsumptionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ConsumptionsDao) Columns() ConsumptionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ConsumptionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ConsumptionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ConsumptionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
