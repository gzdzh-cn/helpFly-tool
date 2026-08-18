package db

import "time"

// Task is the persistence model for a task and its ordered options.
type Task struct {
	Key                  string
	TaskID               string
	TaskName             string
	AccountName          string
	Account              string
	CardNumber           string
	Bank                 string
	DateRange            []string
	DepositType          string
	Currency             string
	TransactionDateRange []string
	CashExchange         string
	Channels             []string
	Summaries            []string
	OpeningBalance       *float64
	AmountType           string
	OrderTime            string
	TransactionCount     *int
	Status               int
	QRCodeURL            string
	StartSerialNumber    string
	CreatedAt            time.Time
}

type Consumption struct {
	Key                   string
	TaskID                string
	TradeDate             string
	Account               string
	StorageType           string
	SerialNumber          string
	Currency              string
	CashOrRemit           string
	Summary               string
	Region                string
	IncomeOrExpenseAmount float64
	Balance               float64
	Channel               string
}

type Export struct {
	Key       string
	TaskID    string
	FilePath  string
	CreatedAt time.Time
}

type PausePoint struct {
	TaskID           string
	LastSerialNumber int
	CurrentBalance   float64
	CurrentProgress  int
	Percent          float64
	PausedAt         string
}

type SystemParameters struct {
	Banks         []string
	DepositTypes  []string
	CashExchanges []string
	Currencies    []string
	Channels      []string
	Summaries     []string
	Regions       []string
	ExportPath    string
	AddWatermark  bool
	WatermarkPath string
}
