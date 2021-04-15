package bean

import (
	"errors"
	"github.com/wolot/api/mondo/types"
)

// ----合约相关
type SignedMultisigEvmTx struct {
	Mode     int    `json:"mode"`     // 交易模式，默认为0；0-同步模式 1-全异步 2-半异步；如果tx执行时间较长、网络不稳定、出块慢，建议使用半异步模式。
	Deadline uint64 `json:"deadline"` // 有效截止时间
	GasLimit uint64 `json:"gasLimit"` // gas限额
	GasPrice string `json:"gasPrice"` // gas价格，最低为1
	From     string `json:"from"`     // 多签账户地址
	Nonce    uint64 `json:"nonce"`    // 多签账户nonce
	To       string `json:"to"`       // 交易接受者地址或合约地址
	Value    string `json:"value"`    // 交易金额
	Load     string `json:"load"`     // 合约负载，普通原声币转账时为空
	Memo     string `json:"memo"`     // 备注

	Signature MultiSignature `json:"signature"` // 交易签名
}

type MultiSignature struct {
	PubKey     PubKeyMultisigThreshold `json:"pubKey"`     // 多签公钥
	Signatures []string                `json:"signatures"` // 用户签名列表
}

type PubKeyMultisigThreshold struct {
	K       int      `json:"threshold"` // 多签阈值
	PubKeys []string `json:"pubkeys"`   // 多签用户公钥列表
}

func (tx *SignedMultisigEvmTx) Check() error {
	if !types.ValidAddress(tx.From) {
		return errors.New("invalid from address")
	}
	if tx.Deadline == 0 {
		return errors.New("ignore deadline")
	}
	if tx.GasLimit == 0 {
		return errors.New("ignore gasLimit")
	}
	if tx.GasPrice == "0" {
		return errors.New("ignore gasPrice")
	}
	if len(tx.Memo) > 255 {
		return errors.New("memo length too long")
	}

	if len(tx.To) > 0 && !types.ValidAddress(tx.To) {
		return errors.New("invalid to address")
	}

	if tx.Signature.PubKey.K < 1 {
		return errors.New("bad threshold")
	}

	if len(tx.Signature.PubKey.PubKeys) < tx.Signature.PubKey.K {
		return errors.New("bad public keys")
	}

	if len(tx.Signature.Signatures) < tx.Signature.PubKey.K || len(tx.Signature.Signatures) > len(tx.Signature.PubKey.PubKeys) {
		return errors.New("bad signatures count")
	}

	return nil
}
