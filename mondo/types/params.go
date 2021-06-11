package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
)

type Params struct {
	DPOSBeginHeight  int64 // 从此高度开启DPOS机制 必须>1
	DPOSEachHeight   int64 // 每多少高度清算一次 10240
	DPOSMaxNodeNum   int   // 超级节点数量
	NodeWorkMortgage int64 // 节点至少抵押该数字才会参与DPOS
	NodeMinMortgage  int64 // 节点最小单笔抵押金额
	NodeMinCollect   int64 // 节点最小单笔收集收益金额
	UserMinMortgage  int64 // 用户最小抵押金额
	UserMinCollect   int64 // 用户最小单笔收集收益金额
	UpgradeHeight    int64 // 升级高度，如果不为0，在此高度的EndBlock会Panic等待升级
}

var DefaultParams = Params{
	DPOSBeginHeight:  150000,
	DPOSEachHeight:   20480,
	DPOSMaxNodeNum:   13,
	NodeWorkMortgage: 50000,
	NodeMinMortgage:  10000 * 1e8,
	NodeMinCollect:   10 * 1e8,
	UserMinMortgage:  100 * 1e8,
	UserMinCollect:   10 * 1e8,
	UpgradeHeight:    0,
}

func (p *Params) FromBytes(bz []byte) {
	if err := amino.UnmarshalBinaryBare(bz, p); err != nil {
		panic(err)
	}
}

func (p *Params) ToBytes() []byte {
	buf, err := amino.MarshalBinaryBare(p)
	if err != nil {
		panic(err)
	}
	return buf
}

var (
	mondoConfig = []byte("mondoConfig")
)

func LoadMondoParams(db ethdb.Database) *Params {
	var cfg Params
	buf, _ := db.Get(mondoConfig)
	if len(buf) != 0 {
		if err := amino.UnmarshalBinaryBare(buf, &cfg); err != nil {
			panic(fmt.Sprintf("UnmarshalBinaryBare %v", err))
		}
		return &cfg
	}

	return nil
}

func SaveMondoParams(db ethdb.Database, cfg *Params) {
	buf, err := amino.MarshalBinaryBare(cfg)
	if err != nil {
		panic(fmt.Sprintf("MarshalBinaryBare %v", err))
	}

	if err := db.Put(mondoConfig, buf); err != nil {
		panic(fmt.Sprintf("chaindb.Put %v", err))
	}
}
