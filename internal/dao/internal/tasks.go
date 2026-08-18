// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TasksDao is the data access object for the table tasks.
type TasksDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TasksColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TasksColumns defines and stores column names for the table tasks.
type TasksColumns struct {
	TaskKey              string //
	TaskId               string //
	TaskName             string //
	AccountName          string //
	Account              string //
	CardNumber           string //
	Bank                 string //
	DateRangeStart       string //
	DateRangeEnd         string //
	DepositType          string //
	Currency             string //
	TransactionDateStart string //
	TransactionDateEnd   string //
	CashExchange         string //
	OpeningBalance       string //
	AmountType           string //
	OrderTime            string //
	TransactionCount     string //
	Status               string //
	QrCodeUrl            string //
	StartSerialNumber    string //
	CreatedAt            string //
}

// tasksColumns holds the columns for the table tasks.
var tasksColumns = TasksColumns{
	TaskKey:              "task_key",
	TaskId:               "task_id",
	TaskName:             "task_name",
	AccountName:          "account_name",
	Account:              "account",
	CardNumber:           "card_number",
	Bank:                 "bank",
	DateRangeStart:       "date_range_start",
	DateRangeEnd:         "date_range_end",
	DepositType:          "deposit_type",
	Currency:             "currency",
	TransactionDateStart: "transaction_date_start",
	TransactionDateEnd:   "transaction_date_end",
	CashExchange:         "cash_exchange",
	OpeningBalance:       "opening_balance",
	AmountType:           "amount_type",
	OrderTime:            "order_time",
	TransactionCount:     "transaction_count",
	Status:               "status",
	QrCodeUrl:            "qr_code_url",
	StartSerialNumber:    "start_serial_number",
	CreatedAt:            "created_at",
}

// NewTasksDao creates and returns a new DAO object for table data access.
func NewTasksDao(handlers ...gdb.ModelHandler) *TasksDao {
	return &TasksDao{
		group:    "default",
		table:    "tasks",
		columns:  tasksColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TasksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TasksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TasksDao) Columns() TasksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TasksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TasksDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TasksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
