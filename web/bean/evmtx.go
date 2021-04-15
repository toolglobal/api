package bean

import (
	"github.com/pkg/errors"
	"github.com/wolot/api/mondo/types"
)

const (
	MODE_DEFAULT = 0
	MODE_ASYNC   = 1
	MODE_SYNC    = 2
)

// ----合约相关
type SignedEvmTx struct {
	Mode      int    `json:"mode"`      // 交易模式，默认为0；0-同步模式 1-全异步 2-半异步；如果tx执行时间较长、网络不稳定、出块慢，建议使用半异步模式。
	CreatedAt uint64 `json:"createdAt"` // 时间戳，可选
	GasLimit  uint64 `json:"gasLimit"`  // gas限额
	GasPrice  string `json:"gasPrice"`  // gas价格，最低为1
	Nonce     uint64 `json:"nonce"`     // 交易发起者nonce
	Sender    string `json:"sender"`    // 交易发起者公钥
	Body      struct {
		To    string `json:"to"`    // 交易接受者地址或合约地址
		Value string `json:"value"` // 交易金额
		Load  string `json:"load"`  // 合约负载，普通原声币转账时为空
		Memo  string `json:"memo"`  // 备注
	} `json:"body"`
	Signature string `json:"signature"` // 交易签名
}

func (tx *SignedEvmTx) Check() error {
	if !types.ValidPublicKey(tx.Sender) && !types.ValidAddress(tx.Sender) {
		return errors.New("invalid sender public key or address")
	}
	if tx.GasLimit == 0 {
		return errors.New("ignore gasLimit")
	}
	if tx.GasPrice == "0" {
		return errors.New("ignore gasPrice")
	}
	if len(tx.Body.Memo) > 255 {
		return errors.New("memo length too long")
	}

	if len(tx.Body.To) > 0 && !types.ValidAddress(tx.Body.To) && !types.ValidPublicKey(tx.Body.To) {
		return errors.New("receiver is not OLO address or public key")
	}
	return nil
}

type EvmCallResult = struct {
	Code    uint32 `json:"code"`    // 错误码
	Msg     string `json:"msg"`     // msg
	Ret     string `json:"ret"`     // 返回数据的hex编码
	GasUsed uint64 `json:"gasUsed"` // 消耗的gas
}

type SignEvmTx struct {
	CreatedAt uint64 `json:"createdAt"` // 时间戳，可选
	GasLimit  uint64 `json:"gasLimit"`  // gas限额
	GasPrice  string `json:"gasPrice"`  // gas价格，最低为1
	Nonce     uint64 `json:"nonce"`     // 交易发起者nonce
	Sender    string `json:"sender"`    // 交易发起者公钥
	Body      struct {
		To    string `json:"to"`    // 交易接受者地址或合约地址
		Value string `json:"value"` // 交易金额
		Load  string `json:"load"`  // 合约负载，普通原声币转账时为空
		Memo  string `json:"memo"`  // 备注
	} `json:"body"`
	PrivateKey string // 私钥
}

func (tx *SignEvmTx) Check() error {
	if !types.ValidPublicKey(tx.Sender) && !types.ValidAddress(tx.Sender) {
		return errors.New("invalid sender public key or address")
	}
	if tx.GasLimit == 0 {
		return errors.New("ignore gasLimit")
	}
	if tx.GasPrice == "0" {
		return errors.New("ignore gasPrice")
	}
	if len(tx.Body.Memo) > 255 {
		return errors.New("memo length too long")
	}

	if len(tx.Body.To) > 0 && !types.ValidAddress(tx.Body.To) && !types.ValidPublicKey(tx.Body.To) {
		return errors.New("receiver is not OLO address or public key")
	}
	return nil
}
