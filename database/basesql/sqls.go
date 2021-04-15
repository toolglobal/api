package basesql

// GetInitSQLs get database initialize sqls
//	opt sqls to create operation tables
//	opi sqls to create operation table-indexs
//	qt  sqls to create query tables
//	qi  sqls to create query table-indexs
func (bs *Basesql) GetInitSQLs() (qt, qi []string) {
	qt = []string{
		createV3LedgerSQL,
		createV3TransactionSQL,
		createV3PaymentSQL,
	}
	qi = append(createV3QIndex)

	return
}
