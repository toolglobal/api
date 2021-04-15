package basesql

func (bs *Basesql) GetV3InitSQLs() (qt, qi []string) {
	qt = []string{
		createV3LedgerSQL,
		createV3TransactionSQL,
		createV3PaymentSQL,
	}
	qi = createV3QIndex

	return
}

var (
	createV3QIndex = []string{
		"CREATE INDEX IF NOT EXISTS idx_height ON v3_ledgers (height)",
		"CREATE INDEX IF NOT EXISTS idx_blockHash ON v3_ledgers (blockHash)",
		"CREATE INDEX IF NOT EXISTS idx_createdAt ON v3_ledgers (createdAt)",

		"CREATE INDEX IF NOT EXISTS idx_hash ON v3_transactions (hash)",
		"CREATE INDEX IF NOT EXISTS idx_tx_height ON v3_transactions (height)",
		"CREATE INDEX IF NOT EXISTS idx_typei ON v3_transactions (typei)",
		"CREATE INDEX IF NOT EXISTS idx_sender ON v3_transactions (sender)",
		"CREATE INDEX IF NOT EXISTS idx_receiver ON v3_transactions (receiver)",
		"CREATE INDEX IF NOT EXISTS idx_tx_createdAt ON v3_transactions (createdAt)",

		"CREATE INDEX IF NOT EXISTS idx_pm_hash ON v3_payments (hash)",
		"CREATE INDEX IF NOT EXISTS idx_pm_height ON v3_payments (height)",
		"CREATE INDEX IF NOT EXISTS idx_pm_sender ON v3_payments (sender)",
		"CREATE INDEX IF NOT EXISTS idx_pm_receiver ON v3_payments (receiver)",
		"CREATE INDEX IF NOT EXISTS idx_symbol ON v3_payments (symbol)",
		"CREATE INDEX IF NOT EXISTS idx_contract ON v3_payments (contract)",
		"CREATE INDEX IF NOT EXISTS idx_pm_createdAt ON v3_payments (createdAt)",
	}
)

const (
	createV3LedgerSQL = `CREATE TABLE IF NOT EXISTS v3_ledgers
	( 
		id         INTEGER  PRIMARY KEY AUTOINCREMENT,
		height     INTEGER  NOT NULL,
		blockHash  TEXT     NOT NULL,
		blockSize  INTEGER  NOT NULL,
		validator  TEXT     NOT NULL,
		txCount    INTEGER  NOT NULL,
		gasLimit   INTEGER  NOT NULL,
		gasUsed    INTEGER  NOT NULL,
		gasPrice   TEXT     NOT NULL,
		createdAt  DATETIME NOT NULL 
	);`
	createV3TransactionSQL = `CREATE TABLE IF NOT EXISTS v3_transactions
	( 
		id        INTEGER  PRIMARY KEY AUTOINCREMENT,
		hash      TEXT     NOT NULL,
		height    INTEGER  NOT NULL,
		typei     INTEGER  NOT NULL,
		types     TEXT     NOT NULL,
		sender    TEXT     NOT NULL,
		nonce     INTEGER  NOT NULL,
		receiver  TEXT     NOT NULL,
		value     TEXT     NOT NULL,
		gasLimit  NUMERIC  NOT NULL,
		gasUsed   INTEGER  NOT NULL,
		gasPrice  TEXT     NOT NULL,
		memo      TEXT,
		payload   TEXT,
		events    TEXT,
		codei     INTEGER  NOT NULL,
		codes     TEXT,
		createdAt DATETIME NOT NULL 
	);`
	createV3PaymentSQL = `CREATE TABLE IF NOT EXISTS v3_payments
	( 
		id        INTEGER  PRIMARY KEY AUTOINCREMENT,
		hash      TEXT     NOT NULL,
		height    INTEGER  NOT NULL,
		evName    TEXT     NOT NULL,
		idx       INTEGER  NOT NULL,
		sender    TEXT     NOT NULL,
		receiver  TEXT     NOT NULL,
		symbol    TEXT     NOT NULL,
		contract  TEXT     NOT NULL,
		value     TEXT     NOT NULL,
		createdAt DATETIME NOT NULL 
	);`
)
