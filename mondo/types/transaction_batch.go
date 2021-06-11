package types

import (
	"math/big"
	"strings"
	"sync/atomic"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wolot/api/utils"
)

// 1对多批量交易
type TxBatch struct {
	CreatedAt uint64    // 交易发起时间
	GasLimit  uint64    // gas限额
	GasPrice  *big.Int  // gas价格
	Nonce     uint64    // 交易发起者nonce
	Sender    PublicKey // 交易发起者公钥
	Ops       []TxOp    // 交易列表
	Memo      []byte    // 备注
	Signature []byte    // 交易签名
	hash      atomic.Value
}

type TxOp struct {
	To    PublicKey // 转账接收方公钥或地址，地址时填后20字节
	Value *big.Int  // 交易金额
}

func NewTxBatch() *TxBatch {
	return &TxBatch{}
}

func (tx *TxBatch) Sign(privkey string) ([]byte, error) {
	if strings.HasPrefix(privkey, "0x") {
		privkey = string([]byte(privkey)[2:])
	}

	priv, err := crypto.ToECDSA(ethcmn.Hex2Bytes(privkey))
	if err != nil {
		return nil, err
	}
	return crypto.Sign(tx.SigHash().Bytes(), priv)
}

func (tx *TxBatch) Verify() bool {
	if len(tx.Signature) != 65 {
		return false
	}

	pubkey, err := crypto.SigToPub(tx.SigHash().Bytes(), tx.Signature)
	if err != nil {
		return false
	}

	if crypto.PubkeyToAddress(*pubkey) != tx.Sender.ToAddress().Address {
		return false
	}
	return true
}

func (tx *TxBatch) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *TxBatch) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *TxBatch) SigHash() ethcmn.Hash {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.GasPrice,
		tx.GasLimit,
		tx.Nonce,
		tx.Sender,
		tx.Ops,
	})

	return utils.RLPHash(message)
}

func (tx *TxBatch) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := utils.RLPHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *TxBatch) Cost() *big.Int {
	return new(big.Int).Mul(big.NewInt(tx.GasWanted()), tx.GasPrice)
}

func (tx *TxBatch) GasWanted() int64 {
	return int64(len(tx.Ops)) * TxGas
}
