// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Tasks is the golang structure for table tasks.
type Tasks struct {
	TaskKey              string  `json:"taskKey"              orm:"task_key"               description:""` //
	TaskId               string  `json:"taskId"               orm:"task_id"                description:""` //
	TaskName             string  `json:"taskName"             orm:"task_name"              description:""` //
	AccountName          string  `json:"accountName"          orm:"account_name"           description:""` //
	Account              string  `json:"account"              orm:"account"                description:""` //
	CardNumber           string  `json:"cardNumber"           orm:"card_number"            description:""` //
	Bank                 string  `json:"bank"                 orm:"bank"                   description:""` //
	DateRangeStart       string  `json:"dateRangeStart"       orm:"date_range_start"       description:""` //
	DateRangeEnd         string  `json:"dateRangeEnd"         orm:"date_range_end"         description:""` //
	DepositType          string  `json:"depositType"          orm:"deposit_type"           description:""` //
	Currency             string  `json:"currency"             orm:"currency"               description:""` //
	TransactionDateStart string  `json:"transactionDateStart" orm:"transaction_date_start" description:""` //
	TransactionDateEnd   string  `json:"transactionDateEnd"   orm:"transaction_date_end"   description:""` //
	CashExchange         string  `json:"cashExchange"         orm:"cash_exchange"          description:""` //
	OpeningBalance       float32 `json:"openingBalance"       orm:"opening_balance"        description:""` //
	AmountType           string  `json:"amountType"           orm:"amount_type"            description:""` //
	OrderTime            string  `json:"orderTime"            orm:"order_time"             description:""` //
	TransactionCount     int     `json:"transactionCount"     orm:"transaction_count"      description:""` //
	Status               int     `json:"status"               orm:"status"                 description:""` //
	QrCodeUrl            string  `json:"qrCodeUrl"            orm:"qr_code_url"            description:""` //
	StartSerialNumber    string  `json:"startSerialNumber"    orm:"start_serial_number"    description:""` //
	CreatedAt            string  `json:"createdAt"            orm:"created_at"             description:""` //
}
