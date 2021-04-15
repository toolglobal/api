package datamanager

import (
	"sync"

	"github.com/wolot/api/database"
)

// DBCreator to create db instance
type DBCreator func(dbname string) database.Database

// DataManager data access between app and database
type DataManager struct {
	wdb       database.Database
	rdb       database.Database
	qNeedLock bool
	qLock     sync.Mutex
}

// NewDataManager create data manager
func NewDataManager(dbname string, dbc DBCreator) (*DataManager, error) {
	wdb := dbc(dbname)
	qt, qi := wdb.GetInitSQLs()
	err := wdb.PrepareTables(qt, qi)

	if err != nil {
		return nil, err
	}
	dm := &DataManager{
		wdb:       wdb,
		rdb:       dbc(dbname),
		qNeedLock: true,
	}

	return dm, nil
}

// Close close all dbs
func (m *DataManager) Close() {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}
	if m.wdb != nil {
		m.wdb.Close()
		m.wdb = nil
	}
	if m.rdb != nil {
		m.rdb.Close()
		m.rdb = nil
	}
}

// QTxBegin start database transaction of wdb
func (m *DataManager) QTxBegin() error {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	return m.wdb.Begin()
}

// QTxCommit commit database transaction of wdb
func (m *DataManager) QTxCommit() error {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	return m.wdb.Commit()
}

// QTxRollback rollback database transaction of wdb
func (m *DataManager) QTxRollback() error {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	return m.wdb.Rollback()
}
