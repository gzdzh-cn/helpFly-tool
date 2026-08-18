// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PausePoints is the golang structure of table pause_points for DAO operations like Where/Data.
type PausePoints struct {
	g.Meta           `orm:"table:pause_points, do:true"`
	TaskId           any //
	LastSerialNumber any //
	CurrentBalance   any //
	CurrentProgress  any //
	Percent          any //
	PausedAt         any //
}
