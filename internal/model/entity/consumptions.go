// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Consumptions is the golang structure for table consumptions.
type Consumptions struct {
	RecordKey             string  `json:"recordKey"             orm:"record_key"               description:""` //
	TaskId                string  `json:"taskId"                orm:"task_id"                  description:""` //
	TradeDate             string  `json:"tradeDate"             orm:"trade_date"               description:""` //
	Account               string  `json:"account"               orm:"account"                  description:""` //
	StorageType           string  `json:"storageType"           orm:"storage_type"             description:""` //
	SerialNumber          string  `json:"serialNumber"          orm:"serial_number"            description:""` //
	Currency              string  `json:"currency"              orm:"currency"                 description:""` //
	CashOrRemit           string  `json:"cashOrRemit"           orm:"cash_or_remit"            description:""` //
	Summary               string  `json:"summary"               orm:"summary"                  description:""` //
	Region                string  `json:"region"                orm:"region"                   description:""` //
	IncomeOrExpenseAmount float32 `json:"incomeOrExpenseAmount" orm:"income_or_expense_amount" description:""` //
	Balance               float32 `json:"balance"               orm:"balance"                  description:""` //
	Channel               string  `json:"channel"               orm:"channel"                  description:""` //
}
