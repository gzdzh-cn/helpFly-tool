// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KvStoreDao is the data access object for the table kv_store.
type KvStoreDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KvStoreColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KvStoreColumns defines and stores column names for the table kv_store.
type KvStoreColumns struct {
	ItemKey string //
	Value   string //
}

// kvStoreColumns holds the columns for the table kv_store.
var kvStoreColumns = KvStoreColumns{
	ItemKey: "item_key",
	Value:   "value",
}

// NewKvStoreDao creates and returns a new DAO object for table data access.
func NewKvStoreDao(handlers ...gdb.ModelHandler) *KvStoreDao {
	return &KvStoreDao{
		group:    "default",
		table:    "kv_store",
		columns:  kvStoreColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KvStoreDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KvStoreDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KvStoreDao) Columns() KvStoreColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KvStoreDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KvStoreDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KvStoreDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
