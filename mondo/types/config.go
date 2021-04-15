package types

import "github.com/tendermint/go-amino"

type Config struct {
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

var DefaultConfig = Config{
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

func (p *Config) FromBytes(bz []byte) {
	if err := amino.UnmarshalBinaryBare(bz, p); err != nil {
		panic(err)
	}
}

func (p *Config) ToBytes() []byte {
	buf, err := amino.MarshalBinaryBare(p)
	if err != nil {
		panic(err)
	}
	return buf
}
