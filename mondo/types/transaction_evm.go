package types

import (
	"math/big"
	"strings"
	"sync/atomic"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// evm交易时sender和receiver都传public key，如果是合约帐户，receiver填写最后20个字节，前12字节填0
type TxEvm struct {
	CreatedAt uint64      // 交易发起时间
	GasLimit  uint64      // gas限额
	GasPrice  *big.Int    // gas价格
	Nonce     uint64      // 交易发起者nonce
	Sender    PublicKey   // 交易发起者公钥
	Body      TxEvmCommon // 交易结构
	Signature []byte      // 交易签名
	hash      atomic.Value
}

type TxEvmCommon struct {
	To    PublicKey // 交易接收方公钥或地址，地址时填后20字节，创建合约是为全0
	Value *big.Int  // 交易金额
	Load  []byte    // 合约交易负荷
	Memo  []byte    // 备注信息
}

func NewTxEvm() *TxEvm {
	return &TxEvm{}
}

func (tx *TxEvm) Sign(privkey string) ([]byte, error) {
	if strings.HasPrefix(privkey, "0x") {
		privkey = string([]byte(privkey)[2:])
	}

	priv, err := crypto.ToECDSA(ethcmn.Hex2Bytes(privkey))
	if err != nil {
		return nil, err
	}
	return crypto.Sign(tx.SigHash().Bytes(), priv)
}

func (tx *TxEvm) Verify() bool {
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

func (tx *TxEvm) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *TxEvm) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *TxEvm) SigHash() ethcmn.Hash {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.GasPrice,
		tx.GasLimit,
		tx.Nonce,
		tx.Sender,
		tx.Body,
	})

	return RLPHash(message)
}

func (tx *TxEvm) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := RLPHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *TxEvm) Cost() *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(tx.GasLimit)), tx.GasPrice)
}
