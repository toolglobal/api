package types

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"sync/atomic"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type MultisigEvmTx struct {
	Deadline  uint64         // 有效期 截止时间 秒
	GasLimit  uint64         // gas限额
	GasPrice  *big.Int       // gas价格
	From      ethcmn.Address // 多签账户地址
	Nonce     uint64         // 多签账户nonce
	To        ethcmn.Address // 交易接收方地址
	Value     *big.Int       // 交易金额
	Load      []byte         // 合约交易负荷
	Memo      []byte         // 备注信息
	Signature Signature      // 签名信息
	sigHash   atomic.Value
	hash      atomic.Value
}

func NewMultisigEvmTx(k int, pks []PublicKey) *MultisigEvmTx {
	return &MultisigEvmTx{
		Deadline: 0,
		GasLimit: 0,
		GasPrice: big.NewInt(0),
		From:     ethcmn.Address{},
		Nonce:    0,
		To:       ethcmn.Address{},
		Value:    big.NewInt(0),
		Load:     nil,
		Memo:     nil,
		Signature: Signature{
			PubKey:   NewPubKeyMultisigThreshold(k, pks),
			MultiSig: NewMultisig(len(pks)),
		},
		hash: atomic.Value{},
	}
}

type Signature struct {
	PubKey   PubKeyMultisigThreshold
	MultiSig *Multisignature
}

func (tx MultisigEvmTx) Signer() []ethcmn.Address {
	signers := make([]ethcmn.Address, 0, 1)
	for _, v := range tx.Signature.MultiSig.Sigs {
		pkey, err := crypto.SigToPub(tx.SigHash().Bytes(), v)
		if err != nil {
			return nil
		}
		address := crypto.PubkeyToAddress(*pkey)
		signers = append(signers, address)
	}
	return signers
}

func (tx *MultisigEvmTx) Sign(privkey string) error {
	if strings.HasPrefix(privkey, "0x") {
		privkey = strings.TrimPrefix(privkey, "0x")
	}

	priv, err := crypto.ToECDSA(ethcmn.Hex2Bytes(privkey))
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(tx.SigHash().Bytes(), priv)

	var pkey PublicKey
	copy(pkey[:], crypto.CompressPubkey(&priv.PublicKey))

	return tx.Signature.MultiSig.AddSignatureFromPubKey(sig, pkey, tx.Signature.PubKey.PubKeys)
}

func (tx *MultisigEvmTx) AddSign(signature string) error {
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}
	if len(sig) != 65 {
		return errors.New("bad signature")
	}
	pubkey, err := crypto.SigToPub(tx.SigHash().Bytes(), sig)

	var pkey PublicKey
	copy(pkey[:], crypto.CompressPubkey(pubkey))

	return tx.Signature.MultiSig.AddSignatureFromPubKey(sig, pkey, tx.Signature.PubKey.PubKeys)
}

func (tx *MultisigEvmTx) Verify() bool {
	if tx.Signature.PubKey.Address() != tx.From {
		return false
	}

	ok := tx.Signature.PubKey.VerifyBytes(tx.SigHash().Bytes(), tx.Signature.MultiSig.Marshal())
	if !ok {
		return false
	}

	return true
}

func (tx *MultisigEvmTx) SigHash() ethcmn.Hash {
	if sigHash := tx.sigHash.Load(); sigHash != nil {
		return sigHash.(ethcmn.Hash)
	}
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.Deadline,
		tx.GasLimit,
		tx.GasPrice,
		tx.From,
		tx.Nonce,
		tx.To,
		tx.Value,
		tx.Load,
		tx.Memo,
	})
	v := RLPHash(message)
	tx.sigHash.Store(v)
	return v
}

func (tx *MultisigEvmTx) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := RLPHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *MultisigEvmTx) ToBytes() []byte {
	bs, _ := rlp.EncodeToBytes(tx)
	return bs
}

func (tx *MultisigEvmTx) FromBytes(bs []byte) error {
	return rlp.DecodeBytes(bs, tx)
}

func (tx *MultisigEvmTx) Cost() *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(tx.GasLimit)), tx.GasPrice)
}
