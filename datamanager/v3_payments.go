package datamanager

import (
	"database/sql"
	"github.com/wolot/api/database"
)

func (m *DataManager) PrepareV3Payment() (*sql.Stmt, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}
	fields := []database.Feild{
		database.Feild{Name: "hash"},
		database.Feild{Name: "height"},
		database.Feild{Name: "evName"},
		database.Feild{Name: "idx"},
		database.Feild{Name: "sender"},
		database.Feild{Name: "receiver"},
		database.Feild{Name: "symbol"},
		database.Feild{Name: "contract"},
		database.Feild{Name: "value"},
		database.Feild{Name: "createdAt"},
	}
	return m.wdb.Prepare(database.TableV3Payments, fields)
}

func (m *DataManager) AddV3PaymentStmt(stmt *sql.Stmt, data *database.V3Payment) (err error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "hash", Value: data.Hash},
		database.Feild{Name: "height", Value: data.Height},
		database.Feild{Name: "evName", Value: data.EvName},
		database.Feild{Name: "idx", Value: data.Idx},
		database.Feild{Name: "sender", Value: data.Sender},
		database.Feild{Name: "receiver", Value: data.Receiver},
		database.Feild{Name: "symbol", Value: data.Symbol},
		database.Feild{Name: "contract", Value: data.Contract},
		database.Feild{Name: "value", Value: data.Value},
		database.Feild{Name: "createdAt", Value: data.CreatedAt.Unix()},
	}
	_, err = m.wdb.Excute(stmt, fields)

	return err
}

func (m *DataManager) AddV3Payment(data *database.V3Payment) (uint64, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "hash", Value: data.Hash},
		database.Feild{Name: "height", Value: data.Height},
		database.Feild{Name: "evName", Value: data.EvName},
		database.Feild{Name: "idx", Value: data.Idx},
		database.Feild{Name: "sender", Value: data.Sender},
		database.Feild{Name: "receiver", Value: data.Receiver},
		database.Feild{Name: "symbol", Value: data.Symbol},
		database.Feild{Name: "contract", Value: data.Contract},
		database.Feild{Name: "value", Value: data.Value},
		database.Feild{Name: "createdAt", Value: data.CreatedAt.Unix()},
	}

	sqlRes, err := m.wdb.Insert(database.TableV3Payments, fields)
	if err != nil {
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (m *DataManager) QueryV3PaymentsByAddress(address, symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	var wheres [][]database.Where

	where1 := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	where1 = append(where1, database.Where{Name: "sender", Value: address})
	if begin != 0 {
		where1 = append(where1, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where1 = append(where1, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}
	if symbol != "" {
		where1 = append(where1, database.Where{Name: "symbol", Value: symbol})
	}
	if contract != "" {
		where1 = append(where1, database.Where{Name: "contract", Value: contract})
	}

	wheres = append(wheres, where1)

	where2 := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	if begin != 0 {
		where2 = append(where2, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where2 = append(where2, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}
	if symbol != "" {
		where2 = append(where2, database.Where{Name: "symbol", Value: symbol})
	}
	where2 = append(where2, database.Where{Name: "receiver", Value: address})

	wheres = append(wheres, where2)

	orderT, err := database.MakeOrder(order, "id")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("id", cursor, limit)

	var result []database.V3Payment
	err = m.rdb.SelectRowsUnion(database.TableV3Payments, wheres, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *DataManager) QueryV3PaymentsByHash(hash, symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	where = append(where, database.Where{Name: "hash", Value: hash, Op: "="})
	if begin != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}
	if symbol != "" {
		where = append(where, database.Where{Name: "symbol", Value: symbol})
	}
	if contract != "" {
		where = append(where, database.Where{Name: "contract", Value: contract})
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

	var result []database.V3Payment
	err = m.rdb.SelectRows(database.TableV3Payments, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *DataManager) QueryV3PaymentsByHeight(height int64, symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	where = append(where, database.Where{Name: "height", Value: height, Op: "="})
	if begin != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: begin, Op: ">="})
	}
	if end != 0 {
		where = append(where, database.Where{Name: "createdAt", Value: end, Op: "<"})
	}
	if symbol != "" {
		where = append(where, database.Where{Name: "symbol", Value: symbol})
	}
	if contract != "" {
		where = append(where, database.Where{Name: "contract", Value: contract})
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

	var result []database.V3Payment
	err = m.rdb.SelectRows(database.TableV3Payments, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *DataManager) QueryV3AllPayments(symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
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

	if symbol != "" {
		where = append(where, database.Where{Name: "symbol", Value: symbol})
	}
	if contract != "" {
		where = append(where, database.Where{Name: "contract", Value: contract})
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

	var result []database.V3Payment
	err = m.rdb.SelectRows(database.TableV3Payments, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
