package bean

import (
	"errors"
	"github.com/wolot/api/mondo/types"
)

type DelegateTx struct {
	CreatedAt uint64 `json:"createdAt"` // 时间戳
	Sender    string `json:"sender"`    // 交易发起者公钥
	Nonce     uint64 `json:"nonce"`     // 交易发起者nonce
	OpType    uint8  `json:"opType"`    // 操作类型：1-抵押选举 2-赎回 3-领取收益
	OpValue   string `json:"opValue"`   // Op对应的值，OpType=1时为抵押金额 OpType=2时不填或填0 OpType=3时填写领取金额
	Receiver  string `json:"receiver"`  // OpType=1时为选举节点的地址
	Signature string `json:"signature"` // 交易签名
}

func (tx *DelegateTx) Check() error {
	if !types.ValidPublicKey(tx.Sender) && !types.ValidAddress(tx.Sender) {
		return errors.New("invalid public key")
	}

	if tx.Receiver != "" && !types.ValidAddress(tx.Receiver) {
		return errors.New("invalid address")
	}
	return nil
}
