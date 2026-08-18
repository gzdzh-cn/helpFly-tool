// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// KvStore is the golang structure of table kv_store for DAO operations like Where/Data.
type KvStore struct {
	g.Meta  `orm:"table:kv_store, do:true"`
	ItemKey any    //
	Value   []byte //
}
