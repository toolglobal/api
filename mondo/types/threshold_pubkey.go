package types

import (
	"bytes"
	"sort"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type PubKeyMultisigThreshold struct {
	K       uint        `json:"threshold"`
	PubKeys []PublicKey `json:"pubkeys"`
}

func NewPubKeyMultisigThreshold(k int, pubkeys []PublicKey) PubKeyMultisigThreshold {
	if k <= 0 {
		panic("threshold k of n multisignature: k <= 0")
	}
	if len(pubkeys) < k {
		panic("threshold k of n multisignature: len(pubkeys) < k")
	}

	// 按公钥排序
	sort.Slice(pubkeys, func(i, j int) bool {
		return bytes.Compare(pubkeys[i].Bytes(), pubkeys[j].Bytes()) < 0
	})

	//for _, pubkey := range pubkeys {
	//	if pubkey == nil {
	//		panic("nil pubkey")
	//	}
	//}
	return PubKeyMultisigThreshold{uint(k), pubkeys}
}

func (pk PubKeyMultisigThreshold) Address() ethcmn.Address {
	return ethcmn.BytesToAddress(pk.Hash().Bytes()[:20])
}

func (pk PubKeyMultisigThreshold) Bytes() []byte {
	message, _ := rlp.EncodeToBytes(
		[]interface{}{
			pk.K,
			pk.PubKeys,
		})
	return message
}

func (pk PubKeyMultisigThreshold) Hash() (sum ethcmn.Hash) {
	hw := sha3.NewLegacyKeccak256()
	hw.Write(pk.Bytes())
	hw.Sum(sum[:0])
	return
}

func (pk PubKeyMultisigThreshold) VerifyBytes(msg []byte, marshalledSig []byte) bool {
	var sig Multisignature
	err := sig.UnMarshal(marshalledSig)
	//err := cdc.UnmarshalBinaryBare(marshalledSig, &sig)
	if err != nil {
		return false
	}
	size := sig.BitArray.Size()
	// ensure bit array is the correct size
	if len(pk.PubKeys) != size {
		return false
	}
	// ensure size of signature list
	if len(sig.Sigs) < int(pk.K) || len(sig.Sigs) > size {
		return false
	}
	// ensure at least k signatures are set
	if sig.BitArray.NumTrueBitsBefore(size) < int(pk.K) {
		return false
	}
	// index in the list of signatures which we are concerned with.
	sigIndex := 0
	for i := 0; i < size; i++ {
		if sig.BitArray.GetIndex(i) {
			if !pk.PubKeys[i].VerifyBytes(msg, sig.Sigs[sigIndex]) {
				return false
			}
			sigIndex++
		}
	}
	return true
}
