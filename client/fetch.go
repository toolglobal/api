package client

import (
	"context"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Fetcher
type Fetcher interface {
	LastBlockHeight() (int64, error)
	FetchBlockInfo(height int64) (*Block, error)
	FetchBlockResultInfo(height int64) ([]*abcitypes.ResponseDeliverTx, error)
}

// DefaultFetcher Fetcher impl
type DefaultFetcher struct {
	abciRpcClient *http.HTTP
}

type Block struct {
	BlockID tmtypes.BlockID `json:"block_id"`
	Block   *tmtypes.Block  `json:"block"`
}

func NewFetch(rpcRemote string) *DefaultFetcher {
	cli, err := http.New(rpcRemote, "/websocket")
	if err != nil {
		panic(err)
	}
	f := DefaultFetcher{
		abciRpcClient: cli,
	}
	return &f
}

// FetchBlockInfo 获取区块信息
func (f *DefaultFetcher) FetchBlockInfo(height int64) (*Block, error) {
	resp, err := f.abciRpcClient.Block(context.Background(), &height)
	if err != nil {
		return nil, err
	}

	block := Block{
		BlockID: resp.BlockID,
		Block:   resp.Block,
	}

	return &block, nil
}

func (f *DefaultFetcher) FetchBlockResultInfo(height int64) ([]*abcitypes.ResponseDeliverTx, error) {
	resp, err := f.abciRpcClient.BlockResults(context.Background(), &height)
	if err != nil {
		return nil, err
	}
	return resp.TxsResults, nil
}

func (f *DefaultFetcher) LastBlockHeight() (int64, error) {
	result, err := f.abciRpcClient.ABCIInfo(context.Background())
	if err != nil {
		return 0, err
	}
	return result.Response.LastBlockHeight, nil
}
