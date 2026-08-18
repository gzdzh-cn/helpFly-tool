// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PausePoints is the golang structure for table pause_points.
type PausePoints struct {
	TaskId           string  `json:"taskId"           orm:"task_id"            description:""` //
	LastSerialNumber int     `json:"lastSerialNumber" orm:"last_serial_number" description:""` //
	CurrentBalance   float32 `json:"currentBalance"   orm:"current_balance"    description:""` //
	CurrentProgress  int     `json:"currentProgress"  orm:"current_progress"   description:""` //
	Percent          float32 `json:"percent"          orm:"percent"            description:""` //
	PausedAt         string  `json:"pausedAt"         orm:"paused_at"          description:""` //
}
