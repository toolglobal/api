package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
)

type LastBlockInfo struct {
	Height           int64
	StateRoot        ethcmn.Hash
	XStateRoot       ethcmn.Hash
	AppHash          ethcmn.Hash
	PrevHash         ethcmn.Hash
	ValidatorUpdates []ABCIValidator
}

type ABCIValidator struct {
	PubKey [32]byte
	Power  int64
}

func (block *LastBlockInfo) FromBytes(bz []byte) {
	if err := amino.UnmarshalBinaryBare(bz, block); err != nil {
		panic(err)
	}
}

func (block *LastBlockInfo) ToBytes() []byte {
	buf, err := amino.MarshalBinaryBare(block)
	if err != nil {
		panic(err)
	}
	return buf
}
