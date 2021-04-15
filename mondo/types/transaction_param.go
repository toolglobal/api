package types

import (
	"sync/atomic"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	KEY_DPOSBeginHeight  = "DPOSBeginHeight"   // 从此高度开启DPOS机制 必须>1
	KEY_DPOSEachHeight   = "DPOSEachHeight"    // 每多少高度清算一次 10240
	KEY_DPOSMaxNodeNum   = "DPOSMaxNodeNum"    // 超级节点数量
	KEY_NodeWorkMortgage = "NodeWorkMortgage"  // 节点至少抵押该数字才会参与DPOS
	KEY_NodeMinMortgage  = "NodeMinMortgage"   // 节点最小单笔抵押金额
	KEY_NodeMinCollect   = "NodeMinCollect"    // 节点最小单笔收集收益金额
	KEY_UserMinMortgage  = "UserMinMortgage"   // 用户最小抵押金额
	KEY_UserMinCollect   = "UserMinCollect"    // 用户最小单笔收集收益金额
	KEY_UpgradeHeight    = "KEY_UpgradeHeight" // 升级高度
)

// unsafe
type TxParams struct {
	CreatedAt uint64         // 时间戳
	Sender    ed25519.PubKey // 交易发起者公钥
	Nonce     uint64         // 交易发起者nonce
	Key       []byte         // 参数名称
	Value     []byte         // 参数值
	Signature []byte         // signature
	hash      atomic.Value
}

func NewTxParams() *TxParams {
	return &TxParams{
		Sender: make([]byte, 32),
	}
}

func (tx *TxParams) Sign(privkey []byte) []byte {
	priv := ed25519.PrivKey(privkey)
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.Key,
		tx.Value,
	})
	sign, _ := priv.Sign(message)
	return sign
}

func (tx *TxParams) Verify() bool {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.Key,
		tx.Value,
	})
	return tx.Sender.VerifySignature(message, tx.Signature)
}

func (tx *TxParams) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *TxParams) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *TxParams) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}
