// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SystemParameters is the golang structure of table system_parameters for DAO operations like Where/Data.
type SystemParameters struct {
	g.Meta        `orm:"table:system_parameters, do:true"`
	Id            any //
	ExportPath    any //
	AddWatermark  any //
	WatermarkPath any //
}
