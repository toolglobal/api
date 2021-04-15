package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
)

// 注意amino的兼容类型
type TxResult struct {
	Ret      []byte      `json:"ret"`      // 返回数据
	Hash     common.Hash `json:"hash"`     // tx hash
	Reversed []byte      `json:"reversed"` // 保留
}

func NewTxResult(ret []byte, hash common.Hash) *TxResult {
	return &TxResult{
		Ret:  ret,
		Hash: hash,
	}
}

func (tx TxResult) Copy() *TxResult {
	return &TxResult{
		Ret:  tx.Ret,
		Hash: tx.Hash,
	}
}

func (tx *TxResult) FromBytes(bs []byte) error {
	return amino.UnmarshalBinaryBare(bs, tx)
}

func (tx TxResult) ToBytes() []byte {
	bs, err := amino.MarshalBinaryBare(tx)
	if err != nil {
		panic("rlp.EncodeToBytes tx:" + err.Error())
	}
	return bs
}
