package types

import (
	"encoding/binary"
	"fmt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
)

type AppBlock struct {
	Height    int64
	StateRoot ethcmn.Hash
}

func NewAppBlock(height int64, root ethcmn.Hash) *AppBlock {
	return &AppBlock{
		Height:    height,
		StateRoot: root,
	}
}

func (b *AppBlock) FromBytes(bz []byte) {
	if err := amino.UnmarshalBinaryBare(bz, b); err != nil {
		panic(err)
	}
}

func (b *AppBlock) ToBytes() []byte {
	buf, err := amino.MarshalBinaryBare(b)
	if err != nil {
		panic(err)
	}
	return buf
}

func heightKey(height int64) []byte {
	keyBz := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBz, uint64(height))
	return keyBz
}

func (b *AppBlock) FromDB(db ethdb.Database, height int64) (*AppBlock, error) {
	buf, _ := db.Get(heightKey(height))
	if len(buf) != 0 {
		if err := amino.UnmarshalBinaryBare(buf, &b); err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func (b *AppBlock) ToDB(db ethdb.Database) error {
	key := heightKey(b.Height)
	buf, err := amino.MarshalBinaryBare(b)
	if err != nil {
		panic(fmt.Sprintf("MarshalBinaryBare %v", err))
	}
	if err := db.Put(key, buf); err != nil {
		panic(fmt.Sprintf("chaindb.Put %v", err))
	}
	return nil
}
