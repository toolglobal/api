package bean

import (
	"github.com/pkg/errors"
	"github.com/toolglobal/api/mondo/types"
)

// 批量交易
type SignedBatchTx struct {
	Mode      int         `json:"mode"`       // 模式:0-default/commit 1-async 2-sync
	CreatedAt uint64      `json:"createdAt"`  // 时间戳，可选字段，秒/毫秒均可
	GasLimit  uint64      `json:"gasLimit"`   // gas限额
	GasPrice  string      `json:"gasPrice"`   // gas价格，至少为1
	Nonce     uint64      `json:"nonce"`      // 用户nonce，每次交易前从链上获取，每次交易nonce+1
	Sender    string      `json:"sender"`     // 交易发起者公钥
	Ops       []Operation `json:"operations"` // 交易列表，数量不可大于10000笔
	Memo      string      `json:"memo"`       // 备注，必须<256字节
	Signature string      `json:"signature"`  // 交易签名的hex字符串
}

type Operation struct {
	To    string `json:"to"`    // 交易接受方地址，可以是普通用户地址、合约地址、节点账户地址
	Value string `json:"value"` // 交易金额
}

func (tx *SignedBatchTx) Check() error {
	if !types.ValidPublicKey(tx.Sender) && !types.ValidAddress(tx.Sender) {
		return errors.New("invalid sender public key")
	}
	if tx.GasLimit == 0 {
		return errors.New("ignore gasLimit")
	}
	if tx.GasPrice == "0" {
		return errors.New("ignore gasPrice")
	}
	if len(tx.Ops) == 0 || len(tx.Ops) > 10000 {
		return errors.New("no operation or too many operations")
	}

	for _, v := range tx.Ops {
		if !types.ValidAddress(v.To) && !types.ValidPublicKey(v.To) {
			return errors.New("bad receiver")
		}
	}

	if len(tx.Memo) > 255 {
		return errors.New("memo length too long")
	}
	return nil
}
