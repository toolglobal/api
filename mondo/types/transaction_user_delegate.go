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

type TxUserDelegate struct {
	CreatedAt uint64    // 交易发起时间
	Sender    PublicKey // 交易发起者公钥
	Nonce     uint64    // 交易发起者nonce
	OpType    uint8     // 操作类型 1-抵押选举 2-赎回 3-领取收益
	OpValue   *big.Int  // Op对应的值，OpType=1时为抵押金额
	Receiver  []byte    // 接受者地址 OpType=1时为节点地址
	Signature []byte    // 交易签名
	hash      atomic.Value
}

const (
	USER_OPTYPE_NULL = uint8(iota)
	USER_OPTYPE_MORTGAGE
	USER_OPTYPE_REEDEM
	USER_OPTYPE_COLLECT
)

func NewUserDelegateTx() *TxUserDelegate {
	return &TxUserDelegate{}
}

func (tx *TxUserDelegate) Sign(privkey string) ([]byte, error) {
	if strings.HasPrefix(privkey, "0x") {
		privkey = string([]byte(privkey)[2:])
	}

	priv, err := crypto.ToECDSA(ethcmn.Hex2Bytes(privkey))
	if err != nil {
		return nil, err
	}

	return crypto.Sign(tx.SigHash().Bytes(), priv)
}

func (tx *TxUserDelegate) Verify() bool {
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

func (tx *TxUserDelegate) SigHash() ethcmn.Hash {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.CreatedAt,
		tx.Sender,
		tx.Nonce,
		tx.OpType,
		tx.OpValue,
		tx.Receiver,
	})

	return utils.RLPHash(message)
}

func (tx *TxUserDelegate) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *TxUserDelegate) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *TxUserDelegate) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := utils.RLPHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *TxUserDelegate) Cost() *big.Int {
	return big.NewInt(tx.GasWanted())
}

func (tx *TxUserDelegate) GasWanted() int64 {
	return TxGas * 100
}
