package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func openTestDB(t *testing.T) context.Context {
	t.Helper()
	// GoFrame caches g.DB() for the process lifetime. Keep one isolated test
	// database for this package instead of replacing the cached default handle.
	if DB() != nil {
		return context.Background()
	}
	t.Setenv("GOFLY_DB_PATH", filepath.Join(t.TempDir(), "app.db"))
	ctx := context.Background()
	if err := Open(ctx); err != nil {
		t.Fatalf("open test database: %v", err)
	}
	return ctx
}

func testTask() Task {
	balance := 1000.5
	count := 3
	return Task{
		Key: "TASK000001", TaskID: "TASK000001", TaskName: "测试任务", AccountName: "张三",
		Account: "6222", CardNumber: "6222", Bank: "测试银行", DateRange: []string{"2026-01-01", "2026-01-31"},
		DepositType: "活期", Currency: "人民币", TransactionDateRange: []string{"2026-01-01", "2026-01-31"},
		CashExchange: "钞", Channels: []string{"渠道一", "渠道二"}, Summaries: []string{"消费"},
		OpeningBalance: &balance, AmountType: "规则一", OrderTime: "2026-01-01 10:00:00",
		TransactionCount: &count, Status: 0, QRCodeURL: "https://example.com/qr", StartSerialNumber: "1",
		CreatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
	}
}

func TestRepositoryRoundTripAndSearch(t *testing.T) {
	ctx := openTestDB(t)
	task := testTask()
	if err := SaveTask(ctx, task); err != nil {
		t.Fatal(err)
	}
	otherTask := testTask()
	otherTask.Key = "OTHERKEY"
	otherTask.TaskID = "OTHER"
	if err := SaveTask(ctx, otherTask); err != nil {
		t.Fatal(err)
	}
	conflictingTask := task
	conflictingTask.Key = "CONFLICTKEY"
	if err := SaveTask(ctx, conflictingTask); err == nil {
		t.Fatal("expected task transaction conflict")
	}
	got, err := GetTask(ctx, task.Key)
	if err != nil || len(got.Channels) != 2 {
		t.Fatalf("task transaction was not rolled back: %#v, %v", got, err)
	}
	got, err = GetTask(ctx, task.Key)
	if err != nil {
		t.Fatal(err)
	}
	if got.TaskName != task.TaskName || len(got.Channels) != 2 || got.Channels[1] != "渠道二" {
		t.Fatalf("task round trip mismatch: %#v", got)
	}

	task.Status = 2
	task.Channels = []string{"新渠道"}
	if err := SaveTask(ctx, task); err != nil {
		t.Fatal(err)
	}
	got, err = GetTask(ctx, task.Key)
	if err != nil || got.Status != 2 || len(got.Channels) != 1 || got.Channels[0] != "新渠道" {
		t.Fatalf("task update mismatch: %#v, %v", got, err)
	}
	if err := UpdateTaskStatus(ctx, task.TaskID, 3); err != nil {
		t.Fatal(err)
	}
	got, err = GetTask(ctx, task.Key)
	if err != nil || got.Status != 3 {
		t.Fatalf("status update mismatch: %#v, %v", got, err)
	}

	consumptions := []Consumption{
		{Key: "C1", TaskID: task.TaskID, TradeDate: "2026-01-02 10:00:00", Account: "6222", StorageType: "活期", SerialNumber: "1", Currency: "人民币", CashOrRemit: "钞", Summary: "工资", Region: "上海", IncomeOrExpenseAmount: 100, Balance: 1100, Channel: "新渠道"},
		{Key: "C2", TaskID: task.TaskID, TradeDate: "2026-01-03 10:00:00", Account: "6222", StorageType: "活期", SerialNumber: "2", Currency: "人民币", CashOrRemit: "钞", Summary: "消费", Region: "北京", IncomeOrExpenseAmount: -20, Balance: 1080, Channel: "新渠道"},
		{Key: "C3", TaskID: "OTHER", TradeDate: "2026-01-03 10:00:00", Account: "9999", StorageType: "活期", SerialNumber: "1", Currency: "人民币", CashOrRemit: "钞", Summary: "工资", Region: "上海", IncomeOrExpenseAmount: 100, Balance: 1100, Channel: "其他"},
	}
	for _, record := range consumptions {
		if err := SaveConsumption(ctx, record); err != nil {
			t.Fatal(err)
		}
	}
	gotConsumptions, err := SearchConsumptions(ctx, task.TaskID, "工资")
	if err != nil || len(gotConsumptions) != 1 || gotConsumptions[0].Key != "C1" {
		t.Fatalf("scoped search mismatch: %#v, %v", gotConsumptions, err)
	}
	gotConsumptions, err = SearchConsumptions(ctx, task.TaskID, "-20")
	if err != nil || len(gotConsumptions) != 1 || gotConsumptions[0].Key != "C2" {
		t.Fatalf("numeric search mismatch: %#v, %v", gotConsumptions, err)
	}
	if err := DeleteConsumptions(ctx, []string{"C1", "C2"}); err != nil {
		t.Fatal(err)
	}

	if err := SavePausePoint(ctx, PausePoint{TaskID: task.TaskID, CurrentBalance: 100, CurrentProgress: 1, Percent: 33, PausedAt: "2026-01-01 12:00:00"}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetPausePoint(ctx, task.TaskID); err != nil {
		t.Fatal(err)
	}
	if err := SaveExport(ctx, Export{Key: "E1", TaskID: task.TaskID, FilePath: "/tmp/test.pdf", CreatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetExport(ctx, "E1"); err != nil {
		t.Fatal(err)
	}
	if err := SetKV(ctx, "test", []byte("value")); err != nil {
		t.Fatal(err)
	}
	value, err := GetKV(ctx, "test")
	if err != nil || string(value) != "value" {
		t.Fatalf("kv mismatch: %q, %v", value, err)
	}
}

func TestParametersReplaceOptionsAndTaskCascade(t *testing.T) {
	ctx := openTestDB(t)
	params := SystemParameters{Banks: []string{"银行一"}, DepositTypes: []string{"活期"}, CashExchanges: []string{"钞"}, Currencies: []string{"人民币"}, Channels: []string{"渠道"}, Summaries: []string{"摘要"}, Regions: []string{"北京"}, ExportPath: "/tmp/export", AddWatermark: true, WatermarkPath: "/tmp/a.png"}
	if err := SaveParameters(ctx, params); err != nil {
		t.Fatal(err)
	}
	params.Banks = []string{"银行二"}
	if err := SaveParameters(ctx, params); err != nil {
		t.Fatal(err)
	}
	got, err := GetParameters(ctx)
	if err != nil || len(got.Banks) != 1 || got.Banks[0] != "银行二" {
		t.Fatalf("parameter replacement mismatch: %#v, %v", got, err)
	}

	task := testTask()
	if err := SaveTask(ctx, task); err != nil {
		t.Fatal(err)
	}
	if err := SaveConsumption(ctx, Consumption{Key: "CASCADE", TaskID: task.TaskID}); err != nil {
		t.Fatal(err)
	}
	if err := SavePausePoint(ctx, PausePoint{TaskID: task.TaskID}); err != nil {
		t.Fatal(err)
	}
	if err := SaveExport(ctx, Export{Key: "CASCADE_EXPORT", TaskID: task.TaskID}); err != nil {
		t.Fatal(err)
	}
	if err := DeleteTasks(ctx, []string{task.Key}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetTask(ctx, task.Key); gerror.Code(err) != gcode.CodeNotFound {
		t.Fatalf("expected task not found, got %v", err)
	}
	if rows, err := ListConsumptions(ctx, task.TaskID); err != nil || len(rows) != 0 {
		t.Fatalf("consumption cascade mismatch: %#v, %v", rows, err)
	}
	if _, err := GetPausePoint(ctx, task.TaskID); gerror.Code(err) != gcode.CodeNotFound {
		t.Fatalf("pause point cascade mismatch: %v", err)
	}
	if _, err := GetExport(ctx, "CASCADE_EXPORT"); gerror.Code(err) != gcode.CodeNotFound {
		t.Fatalf("export cascade mismatch: %v", err)
	}
}

func TestOpenExistingSQLiteDatabase(t *testing.T) {
	if os.Getenv("HELPFLY_EXISTING_DB_TEST") == "" {
		t.Skip("set HELPFLY_EXISTING_DB_TEST=1 to inspect the development database")
	}
	ctx := context.Background()
	if err := Open(ctx); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = Close(ctx) }()
	tasks, err := ListTasks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) == 0 {
		t.Fatal("expected existing development SQLite data")
	}
}
