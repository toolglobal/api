package types

import (
	"fmt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
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

var (
	lastBlockKey = []byte("lastblock")
)

func LoadLastBlock(db ethdb.Database) (lastBlock LastBlockInfo) {
	buf, _ := db.Get(lastBlockKey)
	if len(buf) != 0 {
		if err := amino.UnmarshalBinaryBare(buf, &lastBlock); err != nil {
			panic(fmt.Sprintf("UnmarshalBinaryBare %v", err))
		}
	}

	return lastBlock
}

func SaveLastBlock(db ethdb.Database, appHash ethcmn.Hash, header *AppHeader, valUpdates []ABCIValidator) {
	lastBlock := LastBlockInfo{
		Height:     header.Height.Int64(),
		StateRoot:  header.StateRoot,
		XStateRoot: header.XStateRoot,
		AppHash:    appHash,
		PrevHash:   header.PrevHash,
	}
	lastBlock.ValidatorUpdates = make([]ABCIValidator, len(valUpdates))
	copy(lastBlock.ValidatorUpdates, valUpdates)

	buf, err := amino.MarshalBinaryBare(&lastBlock)
	if err != nil {
		panic(fmt.Sprintf("MarshalBinaryBare %v", err))
	}

	if err := db.Put(lastBlockKey, buf); err != nil {
		panic(fmt.Sprintf("chaindb.Put %v", err))
	}
}
