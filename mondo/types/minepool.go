package types

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
)

type DPosPool struct {
	GenesisTime     time.Time // 创世时间
	LastBlockHeight int64     // 上次出矿区块高度
	LastBlockTime   time.Time // 上次出矿区块时间
}

func (pool *DPosPool) FromBytes(bz []byte) {
	if err := amino.UnmarshalBinaryBare(bz, pool); err != nil {
		panic(err)
	}
}

func (pool *DPosPool) ToBytes() []byte {
	buf, err := amino.MarshalBinaryBare(pool)
	if err != nil {
		panic(err)
	}
	return buf
}

var (
	dPosPoolKey = []byte("DPosPool")
)

func LoadDPosPool(db ethdb.Database) *DPosPool {
	var pool DPosPool
	buf, _ := db.Get(dPosPoolKey)
	if len(buf) != 0 {
		if err := amino.UnmarshalBinaryBare(buf, &pool); err != nil {
			panic(fmt.Sprintf("UnmarshalBinaryBare %v", err))
		}
		return &pool
	}

	return nil
}

func SaveDPosPool(db ethdb.Database, pool *DPosPool) {
	buf, err := amino.MarshalBinaryBare(pool)
	if err != nil {
		panic(fmt.Sprintf("MarshalBinaryBare %v", err))
	}

	if err := db.Put(dPosPoolKey, buf); err != nil {
		panic(fmt.Sprintf("chaindb.Put %v", err))
	}
}
