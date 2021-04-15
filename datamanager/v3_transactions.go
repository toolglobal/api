package datamanager

import (
	"database/sql"
	"github.com/wolot/api/database"
)

func (m *DataManager) PrepareV3Transaction() (*sql.Stmt, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}
	fields := []database.Feild{
		database.Feild{Name: "hash"},
		database.Feild{Name: "height"},
		database.Feild{Name: "typei"},
		database.Feild{Name: "types"},
		database.Feild{Name: "sender"},
		database.Feild{Name: "nonce"},
		database.Feild{Name: "receiver"},
		database.Feild{Name: "value"},
		database.Feild{Name: "gasLimit"},
		database.Feild{Name: "gasUsed"},
		database.Feild{Name: "gasPrice"},
		database.Feild{Name: "memo"},
		database.Feild{Name: "payload"},
		database.Feild{Name: "events"},
		database.Feild{Name: "codei"},
		database.Feild{Name: "codes"},
		database.Feild{Name: "createdAt"},
	}

	return m.wdb.Prepare(database.TableV3Transactions, fields)
}

func (m *DataManager) AddV3TransactionStmt(stmt *sql.Stmt, data *database.V3Transaction) (err error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "hash", Value: data.Hash},
		database.Feild{Name: "height", Value: data.Height},
		database.Feild{Name: "typei", Value: data.Typei},
		database.Feild{Name: "types", Value: data.Types},
		database.Feild{Name: "sender", Value: data.Sender},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "receiver", Value: data.Receiver},
		database.Feild{Name: "value", Value: data.Value},
		database.Feild{Name: "gasLimit", Value: data.GasLimit},
		database.Feild{Name: "gasUsed", Value: data.GasUsed},
		database.Feild{Name: "gasPrice", Value: data.GasPrice},
		database.Feild{Name: "memo", Value: data.Memo},
		database.Feild{Name: "payload", Value: data.Payload},
		database.Feild{Name: "events", Value: data.Events},
		database.Feild{Name: "codei", Value: data.Codei},
		database.Feild{Name: "codes", Value: data.Codes},
		database.Feild{Name: "createdAt", Value: data.CreatedAt.Unix()},
	}
	_, err = m.wdb.Excute(stmt, fields)
	return err
}

func (m *DataManager) AddV3Transaction(data *database.V3Transaction) (uint64, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "hash", Value: data.Hash},
		database.Feild{Name: "height", Value: data.Height},
		database.Feild{Name: "typei", Value: data.Typei},
		database.Feild{Name: "types", Value: data.Types},
		database.Feild{Name: "sender", Value: data.Sender},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "receiver", Value: data.Receiver},
		database.Feild{Name: "value", Value: data.Value},
		database.Feild{Name: "gasLimit", Value: data.GasLimit},
		database.Feild{Name: "gasUsed", Value: data.GasUsed},
		database.Feild{Name: "gasPrice", Value: data.GasPrice},
		database.Feild{Name: "memo", Value: data.Memo},
		database.Feild{Name: "payload", Value: data.Payload},
		database.Feild{Name: "events", Value: data.Events},
		database.Feild{Name: "codei", Value: data.Codei},
		database.Feild{Name: "codes", Value: data.Codes},
		database.Feild{Name: "createdAt", Value: data.CreatedAt.Unix()},
	}

	sqlRes, err := m.wdb.Insert(database.TableV3Transactions, fields)
	if err != nil {
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (m *DataManager) QueryV3SingleTx(hash string) (*database.V3Transaction, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "hash", Value: hash},
	}

	var result []database.V3Transaction
	err := m.rdb.SelectRows(database.TableV3Transactions, where, nil, nil, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (m *DataManager) QueryV3AccountTxs(address string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Transaction, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	if address != "" {
		where = append(where, database.Where{Name: "sender", Value: address})
	}
	if begin != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}

	orderT, err := database.MakeOrder(order, "id")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("id", cursor, limit)

	var result []database.V3Transaction
	err = m.rdb.SelectRows(database.TableV3Transactions, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *DataManager) QueryV3BlockTxs(height int64, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Transaction, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	if height != 0 {
		where = append(where, database.Where{Name: "height", Value: height})
	}
	if begin != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}

	orderT, err := database.MakeOrder(order, "id")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("id", cursor, limit)

	var result []database.V3Transaction
	err = m.rdb.SelectRows(database.TableV3Transactions, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *DataManager) QueryV3AllTxs(begin, end uint64, cursor, limit uint64, order string) ([]database.V3Transaction, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	if begin != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}
	//if order == "ASC" || order == "asc" {
	//	where = append(where, database.Where{Name: "id", Value: cursor * limit, Op: ">"})
	//	cursor = 0
	//}

	orderT, err := database.MakeOrder(order, "id")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("id", cursor, limit)

	var result []database.V3Transaction
	err = m.rdb.SelectRows(database.TableV3Transactions, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
