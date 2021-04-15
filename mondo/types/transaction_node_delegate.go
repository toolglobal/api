package types

import (
	"math/big"
	"sync/atomic"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type TxNodeDelegate struct {
	CreatedAt uint64         // 交易发起时间
	Sender    ed25519.PubKey // 交易发起者公钥
	Nonce     uint64         // 交易发起者nonce
	OpType    uint8          // 操作类型 1-抵押 2-赎回 3-领取收益 4-提现
	OpValue   *big.Int       // Op对应的值 OpType=1时为抵押金额；OpType=4时为体现金额
	Receiver  []byte         // 提现地址,OpType=4时有效
	Signature []byte         // 交易签名
	hash      atomic.Value
}

const (
	NODE_OPTYPE_NULL = uint8(iota)
	NODE_OPTYPE_MORTGAGE
	NODE_OPTYPE_REEDEM
	NODE_OPTYPE_COLLECT
	NODE_OPTYPE_WITHDRAW
)

func NewNodeDelegateTx() *TxNodeDelegate {
	return &TxNodeDelegate{
		Sender:  make([]byte, 32),
		OpValue: new(big.Int),
	}
}

func (tx *TxNodeDelegate) Sign(privkey []byte) []byte {
	priv := ed25519.PrivKey(privkey)
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.OpType,
		tx.OpValue,
		tx.Receiver,
	})
	sign, _ := priv.Sign(message)
	return sign
}

func (tx *TxNodeDelegate) Verify() bool {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.OpType,
		tx.OpValue,
		tx.Receiver,
	})
	return tx.Sender.VerifySignature(message, tx.Signature)
}

func (tx *TxNodeDelegate) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *TxNodeDelegate) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *TxNodeDelegate) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *TxNodeDelegate) Cost() *big.Int {
	return big.NewInt(tx.GasWanted())
}

func (tx *TxNodeDelegate) GasWanted() int64 {
	return TxGas * 100
}
