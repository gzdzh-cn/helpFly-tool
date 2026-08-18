// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// SystemParameters is the golang structure for table system_parameters.
type SystemParameters struct {
	Id            int    `json:"id"            orm:"id"             description:""` //
	ExportPath    string `json:"exportPath"    orm:"export_path"    description:""` //
	AddWatermark  int    `json:"addWatermark"  orm:"add_watermark"  description:""` //
	WatermarkPath string `json:"watermarkPath" orm:"watermark_path" description:""` //
}
