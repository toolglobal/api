package types

import (
	"bytes"
	"errors"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrUnknownTxType = errors.New("unknown tx type")
)

const TxTagLength = 2

type TxTag [TxTagLength]byte

func NewTxTag(tag []byte) TxTag {
	t := TxTag{}
	for i := range tag {
		if i > 1 {
			break
		}
		t[i] = tag[i]
	}
	return t
}

func (t TxTag) Bytes() []byte { return t[:] }

var (
	//TxTagAppInit        = TxTag{0, 0}     // 账户迁移，采用与batch相同的交易结构，但是不收取手续费 v4废弃
	//TxTagTinInit        = TxTag{0, 1}     // v3初始化TIN用户 v4废弃
	//TxTagAppOLO         = TxTag{1, 1}     // 原生交易 v1 - v3废弃
	TxTagAppEvm = TxTag{1, 2} // 合约交易
	//TxTagAppFee         = TxTag{1, 3}     // 收取手续费 - v3废弃
	TxTagAppBatch       = TxTag{1, 4}     // 批量交易
	TxTagEthereumTx     = TxTag{1, 5}     // 以太坊兼容交易
	TxTagNodeDelegate   = TxTag{2, 1}     // 节点抵押赎回提现
	TxTagUserDelegate   = TxTag{2, 2}     // 用户抵押赎回提现
	TxTagAppEvmMultisig = TxTag{3, 2}     // EVM多签交易
	TxTagAppParams      = TxTag{255, 0}   // 修改APP参数
	TxTagAppMgr         = TxTag{255, 255} // 节点管理交易
)

type HashTx interface {
	Hash() ethcmn.Hash
}

// DecodeTx 当是不支持的解析类型， interface == nil
func DecodeTx(raw []byte) (TxTag, HashTx, error) {
	var inputData []byte
	if len(raw) > 2 {
		inputData = raw[2:]
	}
	switch {
	case bytes.HasPrefix(raw, TxTagAppEvm[:]):
		var tx TxEvm
		err := tx.FromBytes(inputData)
		return TxTagAppEvm, &tx, err
	case bytes.HasPrefix(raw, TxTagAppBatch[:]):
		var tx TxBatch
		err := tx.FromBytes(inputData)
		return TxTagAppBatch, &tx, err
	case bytes.HasPrefix(raw, TxTagEthereumTx[:]):
		var tx ethtypes.Transaction
		err := rlp.DecodeBytes(inputData, &tx)
		return TxTagEthereumTx, &tx, err
	case bytes.HasPrefix(raw, TxTagAppMgr[:]):
		var tx TxManage
		err := tx.FromBytes(inputData)
		return TxTagAppMgr, &tx, err
	case bytes.HasPrefix(raw, TxTagNodeDelegate[:]):
		var tx TxNodeDelegate
		err := tx.FromBytes(inputData)
		return TxTagNodeDelegate, &tx, err
	case bytes.HasPrefix(raw, TxTagUserDelegate[:]):
		var tx TxUserDelegate
		err := tx.FromBytes(inputData)
		return TxTagUserDelegate, &tx, err
	case bytes.HasPrefix(raw, TxTagAppEvmMultisig[:]):
		var tx MultisigEvmTx
		err := tx.FromBytes(inputData)
		return TxTagAppEvmMultisig, &tx, err
	case bytes.HasPrefix(raw, TxTagAppParams[:]):
		var tx TxParams
		err := tx.FromBytes(inputData)
		return TxTagAppParams, &tx, err
	}
	return [2]byte{}, nil, nil
}
