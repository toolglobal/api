package dbo

import (
	"github.com/wolot/api/database"
)

func (app *DBO) QueryV3Ledgers(begin, end uint64, cursor, limit uint64, order string) ([]database.V3Ledger, error) {
	return app.dataM.QueryV3AllLedger(begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3LedgerByHeight(height int64) ([]database.V3Ledger, error) {
	ledger, err := app.dataM.QueryV3Ledger(height)
	if err != nil {
		return nil, err
	}
	if ledger == nil {
		return nil, nil
	}
	return []database.V3Ledger{*ledger}, nil
}

func (app *DBO) QueryV3Txs(begin, end uint64, cursor, limit uint64, order string) ([]database.V3Transaction, error) {
	return app.dataM.QueryV3AllTxs(begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3SingleTx(txhash string) ([]database.V3Transaction, error) {
	tx, err := app.dataM.QueryV3SingleTx(txhash)
	if err != nil {
		return nil, err
	}

	if tx == nil {
		return nil, nil
	}
	return []database.V3Transaction{*tx}, nil
}

func (app *DBO) QueryV3AccountTxs(address string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Transaction, error) {
	return app.dataM.QueryV3AccountTxs(address, begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3BlockTxs(height int64, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Transaction, error) {
	return app.dataM.QueryV3BlockTxs(height, begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3Payments(symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	return app.dataM.QueryV3AllPayments(symbol, contract, begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3TxPayments(txhash, symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	return app.dataM.QueryV3PaymentsByHash(txhash, symbol, contract, begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3AccountPayments(address, symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	return app.dataM.QueryV3PaymentsByAddress(address, symbol, contract, begin, end, cursor, limit, order)
}

func (app *DBO) QueryV3BlockPayments(height int64, symbol, contract string, begin, end uint64, cursor, limit uint64, order string) ([]database.V3Payment, error) {
	return app.dataM.QueryV3PaymentsByHeight(height, symbol, contract, begin, end, cursor, limit, order)
}
