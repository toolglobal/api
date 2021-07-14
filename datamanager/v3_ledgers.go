package datamanager

import (
	"github.com/toolglobal/api/database"
)

func (m *DataManager) AddV3Ledger(data *database.V3Ledger) (uint64, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "height", Value: data.Height},
		database.Feild{Name: "blockHash", Value: data.BlockHash},
		database.Feild{Name: "blockSize", Value: data.BlockSize},
		database.Feild{Name: "validator", Value: data.Validator},
		database.Feild{Name: "txCount", Value: data.TxCount},
		database.Feild{Name: "gasLimit", Value: data.GasLimit},
		database.Feild{Name: "gasUsed", Value: data.GasUsed},
		database.Feild{Name: "gasPrice", Value: data.GasPrice},
		database.Feild{Name: "createdAt", Value: data.CreatedAt.Unix()},
	}

	sqlRes, err := m.wdb.Insert(database.TableV3Ledgers, fields)
	if err != nil {
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (m *DataManager) QueryV3Ledger(height int64) (*database.V3Ledger, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "height", Value: height},
	}

	var result []database.V3Ledger
	err := m.rdb.SelectRows(database.TableV3Ledgers, where, nil, nil, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (m *DataManager) QueryV3AllLedger(begin, end uint64, cursor, limit uint64, order string) ([]database.V3Ledger, error) {
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
	// 区块轮询时适用
	//if order == "ASC" || order == "asc" {
	//	where = append(where, database.Where{Name: "id", Value: cursor * limit, Op: ">"})
	//	cursor = 0
	//}

	orderT, err := database.MakeOrder(order, "id")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("id", cursor, limit)

	var result []database.V3Ledger
	err = m.rdb.SelectRows(database.TableV3Ledgers, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
