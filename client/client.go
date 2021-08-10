package client

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/toolglobal/api/datamanager"
	"github.com/toolglobal/api/libs/log"
	"go.uber.org/zap"
	"strings"
	"time"
)

const abijson = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"dst","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"src","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Withdrawal","type":"event"}]`

type Client struct {
	currentHeight int64
	ctx           context.Context
	fetch         Fetcher
	dataMgr       *datamanager.DataManager
	version       int
	tokenMgr      *TokenMgr
	abi           abi.ABI
}

func NewClient(ctx context.Context, tgsBaseURL, chainId string, version int, rpcRemote string, mgr *datamanager.DataManager, startHeight int64) (*Client, error) {
	cli := &Client{
		ctx:      ctx,
		fetch:    Fetcher(NewFetch(rpcRemote)),
		dataMgr:  mgr,
		version:  version,
		tokenMgr: NewTokenMgr(tgsBaseURL, chainId),
	}

	{
		abi, err := abi.JSON(strings.NewReader(abijson))
		if err != nil {
			panic(err)
		}
		cli.abi = abi
	}

	cli.tokenMgr.Start()

	// 获取库里最新的height
	var (
		height int64
		err    error
	)
	height, err = cli.GetCurrentHeightV3()
	if err != nil {
		return nil, err
	}

	// 只有当设置的开始值大于当前数据库的高度，设置才有效
	if startHeight > height {
		height = startHeight
	}

	cli.currentHeight = height

	return cli, nil
}

// Start 保存获取到的数据，手动开启
func (cli *Client) Start() {
	for {
		select {
		case <-cli.ctx.Done():
			return
		default:
			lastBlockHeight, err := cli.LastBlockHeight()
			if err != nil {
				time.Sleep(time.Second)
				log.Logger.Error("LastBlockHeight", zap.Error(err))
				continue
			}

			if cli.currentHeight > lastBlockHeight {
				time.Sleep(time.Second)
				continue
			}

			data, err := cli.GetV3BlockData(cli.currentHeight)
			if err != nil {
				log.Logger.Error("GetV3BlockData", zap.Error(err), zap.Int64("height", cli.currentHeight))
				continue
			}
			if err := cli.SaveV3Data(data); err != nil {
				log.Logger.Error("SaveV3Data", zap.Error(err))
				continue
			}

			log.Logger.Info("fetch ok", zap.Int64("height", cli.currentHeight), zap.Int("version", cli.version))
			cli.currentHeight++
		}
	}
}

func (cli *Client) LastBlockHeight() (height int64, err error) {
	return cli.fetch.LastBlockHeight()
}

func (cli *Client) GetCurrentHeightV3() (int64, error) {
	var height int64

	result, err := cli.dataMgr.QueryV3AllLedger(0, 0, 0, 1, "DESC")
	if err != nil {
		return height, err
	}

	if len(result) == 0 {
		return 1, nil
	}

	height = result[0].Height + 1

	return height, nil
}

func (cli Client) SaveV3Data(data *V3BlockData) error {
	err := cli.dataMgr.QTxBegin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err := cli.dataMgr.QTxRollback(); err != nil {
				log.Logger.Error("client insert rollback", zap.Error(err))
			}
		}
	}()
	// 区块
	_, err = cli.dataMgr.AddV3Ledger(data.ledger)
	if err != nil {
		return err
	}
	txStmt, err := cli.dataMgr.PrepareV3Transaction()
	if err != nil {
		return err
	}
	defer txStmt.Close()
	paymentStmt, err := cli.dataMgr.PrepareV3Payment()
	if err != nil {
		return err
	}
	defer paymentStmt.Close()
	// 交易
	for _, tx := range data.txs {
		err = cli.dataMgr.AddV3TransactionStmt(txStmt, &tx)
		if err != nil {
			return err
		}
	}

	for _, payment := range data.payments {
		err = cli.dataMgr.AddV3PaymentStmt(paymentStmt, &payment)
		if err != nil {
			return err
		}
	}

	if err := cli.dataMgr.QTxCommit(); err != nil {
		return err
	}

	return nil
}
