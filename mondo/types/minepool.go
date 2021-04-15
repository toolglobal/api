package types

import (
	"github.com/tendermint/go-amino"
	"time"
)

type MinePool struct {
	GenesisTime     time.Time // 创世时间
	LastBlockHeight int64     // 上次出矿区块高度
	LastBlockTime   time.Time // 上次出矿区块时间
}

func (pool *MinePool) FromBytes(bz []byte) {
	if err := amino.UnmarshalBinaryBare(bz, pool); err != nil {
		panic(err)
	}
}

func (pool *MinePool) ToBytes() []byte {
	buf, err := amino.MarshalBinaryBare(pool)
	if err != nil {
		panic(err)
	}
	return buf
}
