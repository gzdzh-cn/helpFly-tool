package db

// schemaStatements is intentionally limited to SQLite DDL. Runtime CRUD is
// implemented through the generated GoFrame DAO layer.
var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS tasks (
		task_key TEXT PRIMARY KEY,
		task_id TEXT NOT NULL UNIQUE,
		task_name TEXT NOT NULL DEFAULT '',
		account_name TEXT NOT NULL DEFAULT '',
		account TEXT NOT NULL DEFAULT '',
		card_number TEXT NOT NULL DEFAULT '',
		bank TEXT NOT NULL DEFAULT '',
		date_range_start TEXT NOT NULL DEFAULT '',
		date_range_end TEXT NOT NULL DEFAULT '',
		deposit_type TEXT NOT NULL DEFAULT '',
		currency TEXT NOT NULL DEFAULT '',
		transaction_date_start TEXT NOT NULL DEFAULT '',
		transaction_date_end TEXT NOT NULL DEFAULT '',
		cash_exchange TEXT NOT NULL DEFAULT '',
		opening_balance REAL,
		amount_type TEXT NOT NULL DEFAULT '',
		order_time TEXT NOT NULL DEFAULT '',
		transaction_count INTEGER,
		status INTEGER NOT NULL DEFAULT 0,
		qr_code_url TEXT NOT NULL DEFAULT '',
		start_serial_number TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL
	)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC)`,
	`CREATE TABLE IF NOT EXISTS task_channels (
		task_id TEXT NOT NULL,
		position INTEGER NOT NULL,
		value TEXT NOT NULL,
		PRIMARY KEY(task_id, position),
		FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS task_summaries (
		task_id TEXT NOT NULL,
		position INTEGER NOT NULL,
		value TEXT NOT NULL,
		PRIMARY KEY(task_id, position),
		FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS consumptions (
		record_key TEXT PRIMARY KEY,
		task_id TEXT NOT NULL,
		trade_date TEXT NOT NULL DEFAULT '',
		account TEXT NOT NULL DEFAULT '',
		storage_type TEXT NOT NULL DEFAULT '',
		serial_number TEXT NOT NULL DEFAULT '',
		currency TEXT NOT NULL DEFAULT '',
		cash_or_remit TEXT NOT NULL DEFAULT '',
		summary TEXT NOT NULL DEFAULT '',
		region TEXT NOT NULL DEFAULT '',
		income_or_expense_amount REAL NOT NULL DEFAULT 0,
		balance REAL NOT NULL DEFAULT 0,
		channel TEXT NOT NULL DEFAULT '',
		FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
	)`,
	`CREATE INDEX IF NOT EXISTS idx_consumptions_task_id ON consumptions(task_id)`,
	`CREATE INDEX IF NOT EXISTS idx_consumptions_trade_date ON consumptions(trade_date)`,
	`CREATE INDEX IF NOT EXISTS idx_consumptions_serial_number ON consumptions(serial_number)`,
	`CREATE TABLE IF NOT EXISTS exports (
		record_key TEXT PRIMARY KEY,
		task_id TEXT NOT NULL DEFAULT '',
		file_path TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL,
		FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
	)`,
	`CREATE INDEX IF NOT EXISTS idx_exports_created_at ON exports(created_at DESC)`,
	`CREATE TABLE IF NOT EXISTS pause_points (
		task_id TEXT PRIMARY KEY,
		last_serial_number INTEGER NOT NULL DEFAULT 0,
		current_balance REAL NOT NULL DEFAULT 0,
		current_progress INTEGER NOT NULL DEFAULT 0,
		percent REAL NOT NULL DEFAULT 0,
		paused_at TEXT NOT NULL DEFAULT '',
		FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS system_parameters (
		id INTEGER PRIMARY KEY CHECK(id = 1),
		export_path TEXT NOT NULL DEFAULT '',
		add_watermark INTEGER NOT NULL DEFAULT 0,
		watermark_path TEXT NOT NULL DEFAULT ''
	)`,
	`CREATE TABLE IF NOT EXISTS system_parameter_options (
		group_name TEXT NOT NULL,
		position INTEGER NOT NULL,
		value TEXT NOT NULL,
		PRIMARY KEY(group_name, position)
	)`,
	`CREATE TABLE IF NOT EXISTS kv_store (
		item_key TEXT PRIMARY KEY,
		value BLOB NOT NULL
	)`,
}
