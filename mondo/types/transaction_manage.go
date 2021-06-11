package types

import (
	"sync/atomic"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/wolot/api/utils"
)

type TxManage struct {
	CreatedAt uint64         // 时间戳
	Sender    ed25519.PubKey // 交易发起者公钥
	Nonce     uint64         // 交易发起者nonce
	Receiver  ed25519.PubKey // 交易接受者公钥
	OpType    uint8          // 操作类型 1-设置RunFor 2-设置POWER 3-设置Genesis权限
	OpValue   uint64         // Op对应的值 OpType=1时（OpValue=0，取消runFor；OpValue=1，设置runFor）；OpType=2时是TMCore票权
	Signature []byte         // signature
	hash      atomic.Value
}

var (
	// 竞选超级节点
	OpType_SetRunFor = uint8(1)
	// 更改投票权
	OpType_SetPower = uint8(2)
	// 转移genesis权限
	OpType_SetGenesis = uint8(3)
)

func NewMgrTx() *TxManage {
	return &TxManage{
		Sender:   make([]byte, 32),
		Receiver: make([]byte, 32),
	}
}

func (tx *TxManage) Sign(privkey []byte) []byte {
	priv := ed25519.PrivKey(privkey)
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.Receiver,
		tx.OpType,
		tx.OpValue,
	})
	sign, _ := priv.Sign(message)
	return sign
}

func (tx *TxManage) Verify() bool {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.Receiver,
		tx.OpType,
		tx.OpValue,
	})
	return tx.Sender.VerifySignature(message, tx.Signature)
}

func (tx *TxManage) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *TxManage) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *TxManage) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := utils.RLPHash(tx)
	tx.hash.Store(v)
	return v
}
